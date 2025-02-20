package tss

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	rd "mpc/internal/db/redis"
	"mpc/pkg/logger"
	pb "mpc/proto"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddress  = "localhost:50051"
	rpcTimeout   = 5 * time.Minute
	keygenPrefix = "keygen:"
	signPrefix   = "sign:"
)

// Configuration for TSS operations
type Config struct {
	Parties   []uint32
	Threshold uint32
}

// TSS handles threshold signature operations
type TSS struct {
	redisClient *rd.Client
	rpcClient   pb.MPCServiceClient
	config      Config
}

// NewTSS creates a new TSS instance with connection pooling
func NewTSS(redisClient *rd.Client) (*TSS, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &TSS{
		redisClient: redisClient,
		rpcClient:   pb.NewMPCServiceClient(conn),
		config: Config{
			Parties:   []uint32{1, 2, 3},
			Threshold: 2,
		},
	}, nil
}

// CreateWallet initiates key generation for a new wallet
func (t *TSS) CreateWallet(ctx context.Context, sessionID string) (shareData []byte, publicKey string, err error) {
	// Set up context with timeout
	ctx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()

	// Notify key generation action
	_, err = t.rpcClient.NotifyAction(ctx, &pb.ActionRequest{
		SessionId: sessionID,
		Parties:   t.config.Parties,
		Threshold: t.config.Threshold,
		Action:    pb.Action_INIT_KEYGEN,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to notify keygen action: %w", err)
	}

	// Subscribe to keygen results channel
	keygenChannel := fmt.Sprintf("%s%s", keygenPrefix, sessionID)
	logger.Debug("Subscribing to keygen channel: " + keygenChannel)
	pubsub := t.redisClient.Subscribe(ctx, keygenChannel)
	defer pubsub.Close()

	// Wait for results with timeout
	shareData, publicKey, err = processKeygenResult(ctx, pubsub.Channel())
	if err != nil {
		return nil, "", fmt.Errorf("key generation failed: %w", err)
	}

	return shareData, publicKey, nil
}

// Sign creates a threshold signature for the given message
func (t *TSS) Sign(ctx context.Context, sessionID string, shareData []byte, message []byte) ([]byte, error) {
	// Set up context with timeout
	ctx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()

	// Notify signing action
	_, err := t.rpcClient.NotifyAction(ctx, &pb.ActionRequest{
		SessionId: sessionID,
		Parties:   t.config.Parties,
		Threshold: t.config.Threshold,
		MsgHash:   message,
		ShareData: shareData,
		Action:    pb.Action_INIT_SIGN,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to notify signing action: %w", err)
	}

	// Subscribe to signing results channel
	signChannel := fmt.Sprintf("%s%s", signPrefix, sessionID)
	pubsub := t.redisClient.Subscribe(ctx, signChannel)
	defer pubsub.Close()

	// Wait for signature with timeout
	signature, err := processSignResult(ctx, pubsub.Channel())
	if err != nil {
		return nil, fmt.Errorf("signature generation failed: %w", err)
	}

	return signature, nil
}

// processKeygenResult handles the keygen result from Redis PubSub
func processKeygenResult(ctx context.Context, ch <-chan *redis.Message) ([]byte, string, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil, "", fmt.Errorf("channel closed")
			}

			var result map[string]string
			if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
				log.Printf("Warning: Failed to parse JSON: %v", err)
				continue
			}

			shareData, ok := result["share_data"]
			if !ok {
				continue // Incomplete message, wait for next one
			}

			publicKey, ok := result["pub_key"]
			if !ok {
				continue // Incomplete message, wait for next one
			}

			// Decode base64 share
			decodedShare, err := base64.StdEncoding.DecodeString(shareData)
			if err != nil {
				log.Printf("Warning: Failed to decode base64 share: %v", err)
				continue
			}

			return decodedShare, publicKey, nil
		}
	}
}

// processSignResult handles the signing result from Redis PubSub
func processSignResult(ctx context.Context, ch <-chan *redis.Message) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil, fmt.Errorf("channel closed")
			}

			var result map[string]string
			if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
				log.Printf("Warning: Failed to parse JSON: %v", err)
				continue
			}

			signatureBase64, ok := result["signature"]
			if !ok {
				continue // Incomplete message, wait for next one
			}

			// Decode base64 signature
			signature, err := base64.StdEncoding.DecodeString(signatureBase64)
			if err != nil {
				log.Printf("Warning: Failed to decode base64 signature: %v", err)
				continue
			}

			return signature, nil
		}
	}
}
