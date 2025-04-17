package block

import "encoding/json"

type Transaction struct {
	senderBlockchainAddress   string
	recipentBlockchainAddress string
	value                     float32
}

func NewTransaction(sender string, recipent string, value float32) *Transaction {
	return &Transaction{sender, recipent, value}
}

func (t *Transaction) MarshalJson() ([]byte, error) {
	return json.Marshal(
		struct {
			SenderBlockchainAddress   string  `json:"sender_blockchain_address"`
			RecipentBlockchainAddress string  `json:"recipent_blockchain_address"`
			Value                     float32 `json:"value"`
		}{
			SenderBlockchainAddress:   t.senderBlockchainAddress,
			RecipentBlockchainAddress: t.recipentBlockchainAddress,
			Value:                     t.value,
		})
}

func (t *Transaction) UnmarshalJson(data []byte) error {
	v := &struct {
		Sender   *string  `json:"sender_blockchain_address"`
		Recipent *string  `json:"recipent_blockchain_address"`
		Value    *float32 `json:"value"`
	}{
		Sender:   &t.senderBlockchainAddress,
		Recipent: &t.recipentBlockchainAddress,
		Value:    &t.value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
