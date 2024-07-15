package blockchain

import (
	"strings"

	"github.com/phantasma-io/phantasma-go/pkg/cryptography"
	crypto "github.com/phantasma-io/phantasma-go/pkg/cryptography"
	"github.com/phantasma-io/phantasma-go/pkg/io"
	hashing "github.com/phantasma-io/phantasma-go/pkg/util/hashing"
)

// Transaction a
type Transaction struct {

	// Code to run in PhantasmaVM for this transaction.
	Script []byte

	NexusName string

	ChainName string

	Expiration uint32

	Payload []byte

	Signatures []crypto.Signature

	Hash crypto.Hash
}

// NewTransaction creates a new transaction object
func NewTransaction(nexusName, chainName string, script []byte, timestamp uint32, payload []byte) Transaction {

	tx := Transaction{
		NexusName:  nexusName,
		ChainName:  chainName,
		Script:     script,
		Expiration: timestamp,
		Payload:    payload,
		Signatures: []crypto.Signature{},
		Hash:       crypto.Hash{},
	}

	tx.updateHash()

	return tx
}

// updateHash sets the hash of the transaction
func (tx *Transaction) updateHash() {
	data := tx.Bytes(false)
	bytes := hashing.Sha256(data)
	hash, err := cryptography.HashFromBytes(bytes)
	if err != nil {
		panic("Updating hash on tx failed!")
	}
	tx.Hash = hash
}

// HasSignatures checks if the transaction was signed already
func (tx *Transaction) HasSignatures() bool {
	return len(tx.Signatures) > 0
}

// Serialize implements ther Serializable interface
func (tx *Transaction) Serialize(writer *io.BinWriter, withSignatures bool) {
	writer.WriteString(tx.NexusName)
	writer.WriteString(tx.ChainName)
	writer.WriteVarBytes(tx.Script)
	writer.WriteU32LE(tx.Expiration)
	writer.WriteVarBytes(tx.Payload)

	if withSignatures {
		writer.WriteVarUint(uint64(len(tx.Signatures)))
		for _, signature := range tx.Signatures {

			if signature == nil {
				writer.WriteB(byte(cryptography.None)) // signaturekind.none
			}
			writer.WriteB(byte(signature.Kind()))
			signature.Serialize(writer)
		}
	}
}

// Deserialize implements ther Serializable interface
func (tx *Transaction) Deserialize(reader *io.BinReader) {
	tx.NexusName = reader.ReadString()
	tx.ChainName = reader.ReadString()
	tx.Script = reader.ReadVarBytes()
	tx.Expiration = reader.ReadU32LE()
	tx.Payload = reader.ReadVarBytes()

	signatureCount := int(reader.ReadVarUint())
	if signatureCount > 0 {
		reader.ReadArray(&tx.Signatures, signatureCount)
	} else {
		tx.Signatures = []crypto.Signature{}
	}
	tx.updateHash()
}

// String a
func (tx *Transaction) String() string {
	return tx.Hash.String()
}

// Bytes a
func (tx *Transaction) Bytes(withSignatures bool) []byte {
	bw := *io.NewBufBinWriter()
	tx.Serialize(bw.BinWriter, withSignatures)
	return bw.Bytes()
}

// Sign the transaction
func (tx *Transaction) Sign(keyPair crypto.KeyPair) {
	if keyPair == nil {
		panic("KeyPair can't be nil!")
	}

	msg := tx.Bytes(false)

	signature := keyPair.Sign(msg)

	tx.Signatures = append(tx.Signatures, signature)
}

// IsSignedBy checks if a transaction is signed by a specific address
func (tx *Transaction) IsSignedBy(addresses []crypto.Address) bool {
	if !tx.HasSignatures() {
		return false
	}

	msg := tx.Bytes(false)

	for _, signature := range tx.Signatures {
		if signature.Verify(msg, addresses) {
			return true
		}
	}

	return false
}

// Mine the transaction with the passed in difficulty
func (tx *Transaction) Mine(difficulty int) {

	//TODO checks

	if difficulty == 0 {
		return
	}

	var nonce uint32 = 0

	for {
		if tx.Hash.GetDifficulty() >= difficulty {
			break
		}

		if nonce == 0 {
			tx.Payload = make([]byte, 4)
		}

		nonce++

		tx.Payload[0] = byte((nonce >> 0) & 0xFF)
		tx.Payload[1] = byte((nonce >> 8) & 0xFF)
		tx.Payload[2] = byte((nonce >> 16) & 0xFF)
		tx.Payload[3] = byte((nonce >> 24) & 0xFF)
		tx.updateHash()
	}
}

func TxStateIsSuccess(state string) bool {
	if strings.ToUpper(state) == "HALT" {
		return true
	} else {
		return false
	}
}

func TxStateIsFault(state string) bool {
	if strings.ToUpper(state) == "FAULT" {
		return true
	} else {
		return false
	}
}

// TODO
//func (tx *Transaction) IsValid(chain Chain) {}
