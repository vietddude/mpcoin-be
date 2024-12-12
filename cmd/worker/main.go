// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"sync"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// // Kiểm tra xem một địa chỉ có trong mảng địa chỉ không
// func isAddressInList(address common.Address, addressList map[common.Address]struct{}) bool {
// 	_, exists := addressList[address]
// 	return exists
// }

// // Xử lý các giao dịch chuyển ETH trong một dải block
// func processBlockRange(client *ethclient.Client, startBlock, endBlock uint64, addressList map[common.Address]struct{}, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	// Lặp qua các block từ startBlock đến endBlock
// 	for i := startBlock; i <= endBlock; i++ {
// 		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
// 		if err != nil {
// 			log.Printf("Error getting block %d: %v", i, err)
// 			continue
// 		}
// 		// Lấy Network ID (chainID)
// 		networkID, err := client.NetworkID(context.Background())
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// Tạo signer với chainID
// 		signer := types.NewEIP155Signer(networkID)

// 		// Lặp qua các giao dịch trong block
// 		for _, tx := range block.Transactions() {
// 			// Kiểm tra địa chỉ `from` trong giao dịch
// 			from, err := types.Sender(signer, tx)
// 			if err != nil {
// 				log.Printf("Error getting sender of tx %s: %v", tx.Hash().Hex(), err)
// 				continue
// 			}

// 			// Nếu `from` hoặc `to` có trong danh sách thì hiển thị
// 			if isAddressInList(from, addressList) || (tx.To() != nil && isAddressInList(*tx.To(), addressList)) {
// 				fmt.Println("Transaction Hash:", tx.Hash().Hex())
// 				fmt.Println("From:", from.Hex())
// 				if to := tx.To(); to != nil {
// 					fmt.Println("To:", to.Hex())
// 				}
// 			}
// 		}
// 	}
// }

// // Quét các block từ startBlock đến block hiện tại
// func watchTransactionsFromBlock(client *ethclient.Client, startBlock uint64, addressList map[common.Address]struct{}) {
// 	// Lấy block hiện tại
// 	currentBlock, err := client.BlockByNumber(context.Background(), nil) // Block hiện tại
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Tính số block mỗi lần quét (số block cần xử lý song song)
// 	blockRangeSize := uint64(100) // Số block trong một lần quét

// 	var wg sync.WaitGroup

// 	// Lặp qua các block ranges để xử lý song song
// 	for start := startBlock; start <= currentBlock.NumberU64(); start += blockRangeSize {
// 		end := start + blockRangeSize - 1
// 		if end > currentBlock.NumberU64() {
// 			end = currentBlock.NumberU64()
// 		}

// 		wg.Add(1)
// 		go processBlockRange(client, start, end, addressList, &wg)
// 	}

// 	// Chờ tất cả goroutines hoàn thành
// 	wg.Wait()
// }

// func main() {
// 	// Kết nối tới client Ethereum
// 	client, err := ethclient.Dial("wss://sepolia.gateway.tenderly.co")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	// Danh sách các địa chỉ ví để theo dõi
// 	addresses := []string{
// 		"0x9b6bE46ed05EE77a22928ba88cA46d9FFf09e3f8",
// 	}

// 	// Chuyển đổi địa chỉ ví sang common.Address và tạo map cho việc tìm kiếm nhanh
// 	addressList := make(map[common.Address]struct{})
// 	for _, address := range addresses {
// 		addressList[common.HexToAddress(address)] = struct{}{}
// 	}

// 	// Gọi hàm theo dõi giao dịch từ block bắt đầu
// 	startBlock := uint64(7203228) // Block bắt đầu từ đây
// 	watchTransactionsFromBlock(client, startBlock, addressList)
// }

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Kiểm tra xem một địa chỉ có trong mảng địa chỉ không
func isAddressInList(address common.Address, addressList map[common.Address]struct{}) bool {
	_, exists := addressList[address]
	return exists
}

func watchTransactions(client *ethclient.Client, addressList map[common.Address]struct{}) {
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Lấy Network ID (chainID)
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Tạo signer với chainID
	signer := types.NewEIP155Signer(networkID)

	// Theo dõi các block mới
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			// Lấy thông tin block từ header mới
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			// Xử lý các giao dịch trong block mới
			for _, tx := range block.Transactions() {
				// Lấy sender của giao dịch
				from, err := types.Sender(signer, tx)
				if err != nil {
					continue
				}

				// Kiểm tra xem từ địa chỉ hoặc tới địa chỉ có trong danh sách không
				// if isAddressInList(from, addressList) {
				// 	fmt.Println("Transaction Hash:", tx.Hash().Hex())
				// 	fmt.Println("From:", from.Hex())
				// }
				// if to := tx.To(); to != nil && isAddressInList(*to, addressList) {
				// 	fmt.Println("To:", to.Hex())
				// }

				fmt.Println("Transaction Hash:", tx.Hash().Hex())
				fmt.Println("From:", from.Hex())
				if to := tx.To(); to != nil {
					fmt.Println("To:", to.Hex())
				}
				fmt.Print("Timestamp: ", block.Time(), "\n-----------------\n")
			}
		}
	}
}

func processTransaction(txHash string, from string, to string, timestamp uint64) {

}

func main() {
	// Kết nối tới client Ethereum
	client, err := ethclient.Dial("wss://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Danh sách các địa chỉ ví để theo dõi
	addresses := []string{
		"0x9b6bE46ed05EE77a22928ba88cA46d9FFf09e3f8",
	}

	// Chuyển đổi địa chỉ ví sang common.Address và tạo map cho việc tìm kiếm nhanh
	addressList := make(map[common.Address]struct{})
	for _, address := range addresses {
		addressList[common.HexToAddress(address)] = struct{}{}
	}

	// Gọi hàm theo dõi giao dịch từ các block mới
	watchTransactions(client, addressList)
}
