package counter

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

func init() {
	client, _ = ethclient.Dial("https://api.zan.top/node/v1/eth/sepolia/<api key>")
}

func loadContractHexBytes() []byte {
	// 读取文件内容
	content, err := os.ReadFile("./counter_sol_counter.bin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
	data := make([]byte, len(content))
	n, err := hex.Decode(data, content)
	if err != nil {
		log.Fatal(err)
	}
	data = data[:n]
	return data
}

func callCounter() {
	privateKey, err := crypto.HexToECDSA("<private key>")
	if err != nil {
		log.Fatal(err)
	}
	//receipt := Deploy(privateKey)
	//address := receipt.ContractAddress
	//0x9c3A954A30a052C89BD74EAF61c816a52ABa1BF2
	address := common.HexToAddress("0x9c3A954A30a052C89BD74EAF61c816a52ABa1BF2")
	counter, err := NewCounter(address, client)
	if err != nil {
		log.Fatal(err)
	}
	// 初始化交易opt实例
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}
	count, err := counter.Count(opt)
	if err != nil {
		log.Fatal(err)
	}
	_, err = waitForReceipt(count.Hash())
	if err != nil {
		log.Fatal(err)
	}

	num, err := counter.Num(nil)
	fmt.Println("num", num)
}

func Deploy(privateKey *ecdsa.PrivateKey) *types.Receipt {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	txData := &types.LegacyTx{
		Nonce:    nonce,
		Value:    big.NewInt(0),
		Gas:      3000000,
		GasPrice: gasPrice,
		Data:     loadContractHexBytes(),
	}
	tx := types.NewTx(txData)
	// 签名交易
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	receipt, err := waitForReceipt(signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := receipt.ContractAddress.Hex()
	fmt.Println("contract address:", contractAddress)
	return receipt
}

func waitForReceipt(hash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), hash)
		if err == nil {
			return receipt, nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		time.Sleep(1 * time.Second)
	}
}
