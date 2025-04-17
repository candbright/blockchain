package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func (b *Block) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		b.nonce,
		fmt.Sprintf("%x", b.previousHash),
		b.timestamp,
		b.transactions,
	})
}

func (b *Block) UnmarshalJson(data []byte) error {
	var previousHash string
	v := struct {
		Nonce        *int            `json:"nonce"`
		PreviousHash *string         `json:"previous_hash"`
		Timestamp    *int64          `json:"timestamp"`
		Transactions *[]*Transaction `json:"transactions"`
	}{
		Nonce:        &b.nonce,
		PreviousHash: &previousHash,
		Timestamp:    &b.timestamp,
		Transactions: &b.transactions,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.nonce = nonce
	b.previousHash = previousHash
	b.timestamp = time.Now().UnixNano()
	b.transactions = transactions
	return b
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Transactions() []*Transaction {
	return b.transactions
}

func (b *Block) Nonce() int {
	return b.nonce
}
