package main

import (
	"blockchain/block"
	"blockchain/wallet"
)

var cache = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port: port}
}

func (s *BlockchainServer) Port() uint16 {
	return s.port
}

func (s *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), s.Port())
		cache["blockchain"] = bc
	}
	return bc
}
func (s *BlockchainServer) Run() {

}
