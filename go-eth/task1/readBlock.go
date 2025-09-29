package task1

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func Link() {
	client, err := ethclient.Dial("https://api.zan.top/node/v1/eth/sepolia/<api key>")
	if err != nil {
		log.Fatal(err)
	}
	// 查询最新区块
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块号", block.Number())
	fmt.Println("区块哈希", block.Hash().Hex())
	fmt.Println("区块时间戳", block.Time())
	fmt.Println("交易数", len(block.Transactions()))

	// 转账
	transfer(client)
}

func transfer(client *ethclient.Client) {
	// 加载私钥
	privateKey, err := crypto.HexToECDSA("<private key>")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("地址", fromAddress.Hex())

	value := big.NewInt(100000000000) // in wei (1 eth)
	gasLimit := uint64(21000)         // in units
	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 获取建议的gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("<to address>")

	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     make([]byte, 0),
	}
	tx := types.NewTx(txData)
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

	// 查询交易收据
	receipt(client, signedTx)
}

func receipt(client *ethclient.Client, txHash *types.Transaction) {
	for {
		fmt.Println("signedTx", txHash.Hash().Hex())
		rec, err := client.TransactionReceipt(context.Background(), txHash.Hash())
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				time.Sleep(1 * time.Second)
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Println(rec.Logs)
			fmt.Println(rec.Status)
			break
		}
	}

}

func genAddress() {
	// 1. 生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("生成私钥失败:", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("私钥", hexutil.Encode(privateKeyBytes)[2:])
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("公钥", hexutil.Encode(publicKeyBytes[4:]))
	fmt.Println("地址", crypto.PubkeyToAddress(*publicKeyECDSA).Hex())
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("full:", hexutil.Encode(hash.Sum(nil)[:]))
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去12位，保留后20位
}
