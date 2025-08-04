package main

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockchainParser struct {
	contractABI string
	contractAddr common.Address
	client       bind.ContractCaller
}

type ParsedData struct {
	BlockNumber uint64
	TxHash      common.Hash
	TxValue     uint64
	Sender      common.Address
	Receiver    common.Address
}

func NewBlockchainParser(contractABI string, contractAddr common.Address, client bind.ContractCaller) *BlockchainParser {
	return &BlockchainParser{
		contractABI: contractABI,
		client:       client,
		contractAddr: contractAddr,
	}
}

func (bp *BlockchainParser) ParseBlock(blockNum uint64) ([]ParsedData, error) {
	var parsedData []ParsedData
	block, err := bp.client.PendingCodeAt(common.HexToAddress("0x0000000000000000000000000000000000000000"), blockNum)
	if err != nil {
		return nil, err
	}
	for _, tx := range block.Transactions() {
		txHash := tx.Hash()
		receipt, err := bp.client.TransactionReceipt(txHash)
		if err != nil {
			return nil, err
		}
		parsedData = append(parsedData, ParsedData{
			BlockNumber: blockNum,
			TxHash:      txHash,
			TxValue:     tx.Value(),
			Sender:      tx.From(),
			Receiver:    tx.To(),
		})
	}
	return parsedData, nil
}

func (bp *BlockchainParser) ParseTx(txHash common.Hash) (*ParsedData, error) {
	tx, _, err := bp.client.TransactionByHash(txHash)
	if err != nil {
		return nil, err
	}
	receipt, err := bp.client.TransactionReceipt(txHash)
	if err != nil {
		return nil, err
	}
	return &ParsedData{
		BlockNumber: receipt.BlockNumber,
		TxHash:      txHash,
		TxValue:     tx.Value(),
		Sender:      tx.From(),
		Receiver:    tx.To(),
	}, nil
}

func main() {
	contractABI := "path/to/contract/abi.json"
	contractAddr := common.HexToAddress("0x000000000000000000000000000000000000000000")
	client, err := bind.NewBoundContract(contractAddr, contractABI, nil, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	parser := NewBlockchainParser(contractABI, contractAddr, client)
	parsedData, err := parser.ParseBlock(100)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonData, err := json.Marshal(parsedData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonData))
}