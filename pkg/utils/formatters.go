package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// FormatUUID converts a string to UUID, panics if invalid
func FormatUUID(id string) uuid.UUID {
	return uuid.MustParse(id)
}

// ToPgUUID converts a UUID to pgtype.UUID
func ToPgUUID(id uuid.UUID) pgtype.UUID {
	var bytes [16]byte
	copy(bytes[:], id[:])
	return pgtype.UUID{
		Valid: true,
		Bytes: bytes,
	}
}

// ToUUID converts a pgtype.UUID to UUID, returns zero UUID if invalid
func ToUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.UUID{}
	}
	return uuid.UUID(pgUUID.Bytes)
}

// ToPgText converts a string to pgtype.Text
func ToPgText(text string) pgtype.Text {
	return pgtype.Text{String: text, Valid: true}
}

// ToText converts a pgtype.Text to string, returns empty string if invalid
func ToText(pgText pgtype.Text) string {
	if !pgText.Valid {
		return ""
	}
	return pgText.String
}

func ToPgTimestamp(t string) pgtype.Timestamp {
	if t == "" {
		return pgtype.Timestamp{Valid: false}
	}
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: parsed, Valid: true}
}

func GetWalletIDFromContext(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.GetString("wallet_id"))
}

func ToNullablePgUUID(id *uuid.UUID, format string) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return ToPgUUID(*id)
}

func CurrentPgTimestamp() pgtype.Timestamp {
	return ToPgTimestamp(time.Now().Format(time.RFC3339))
}
