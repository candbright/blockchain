package block

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

const MINING_TIMER_SEC = 20

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint16
	mux               sync.Mutex

	neibors    []string
	muxNeibors sync.Mutex
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.port = port
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeibors()
}

func (bc *Blockchain) SetNeibors() {

}

func (bc *Blockchain) SyncNeibors() {
	bc.muxNeibors.Lock()
	defer bc.muxNeibors.Unlock()
	bc.SetNeibors()
}

func (bc *Blockchain) StartSyncNeibors() {
	bc.SyncNeibors()
	_ = time.AfterFunc(time.Second*10, bc.StartSyncNeibors)
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	// 创建一个新的区块，传入的参数为随机数和前一个区块的哈希值
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	// 将新的区块添加到区块链中
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)
	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		if bc.CalculateTotalAmount(sender) < value {
			fmt.Println("ERROR: Not enough balance")
		}
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		fmt.Println("ERROR: Verify transaction signature failed")
	}
	return false

}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderBlockchainAddress, t.recipentBlockchainAddress, t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{timestamp: 0, nonce: nonce, previousHash: previousHash, transactions: transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce++
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	if len(bc.transactionPool) == 0 {
		return false
	}
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	return true
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipentBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress && blockchainAddress != MINING_SENDER {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256(m)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	preBlock := chain[0]
	currentBlockIndex := 1
	for currentBlockIndex < len(chain) {
		currentBlock := chain[currentBlockIndex]
		if currentBlock.previousHash != preBlock.Hash() {
			return false
		}
		if !bc.ValidProof(currentBlock.nonce, preBlock.Hash(), currentBlock.transactions, MINING_DIFFICULTY) {
			return false
		}
		preBlock = currentBlock
		currentBlockIndex++
	}
	return true
}
