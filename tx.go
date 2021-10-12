// Package tx defines data types for transactions, transaction statuses, and
// other related types. This package should not contain business-logic, and
// should not depend on business-logic packages.
//
// Transactions are defined in this package, away from business-logic, so that
// they can be imported into other projects without needing to import unwanted
// components. For example, block explorers need to know about transactions, but
// do not need to know about the transaction pool.
package tx

import (
	"math/rand"
	"reflect"
	"testing/quick"

	"github.com/renproject/id"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
)

// Tx represents a RenVM cross-chain transaction.
type Tx struct {
	// Hash of the transaction that uniquely identifies it. It is the SHA256
	// hash of the "to" contract address and the "in" values.
	Hash id.Hash `json:"hash"`

	// Version of the transaction. This allows for backwards/forwards
	// compatibility when processing transactions.
	Version Version `json:"version"`

	// Selector of the gateway function that is being called by this
	// transaction.
	Selector Selector `json:"selector"`

	// Input values are provided by an external user when the transaction is
	// submitted.
	Input pack.Typed `json:"in"`

	// Output values are generated as part of execution.
	Output pack.Typed `json:"out"`
}

// NewTxHash returns the transaction hash for a transaction with the given
// recipient and inputs. An error is returned when the recipient and inputs is
// too large and cannot be marshaled into bytes without exceeding memory
// allocation restrictions.
func NewTxHash(version Version, selector Selector, input pack.Typed) (id.Hash, error) {
	buf := make([]byte, surge.SizeHintString(string(version))+surge.SizeHintString(string(selector))+surge.SizeHint(input))
	return NewTxHashIntoBuffer(version, selector, input, buf)
}

// NewTxHashIntoBuffer write the transaction hash for a transaction with the
// given recipient and inputs into a bytes buffer. An error is returned when the
// recipient and inputs is too large and cannot be marshaled into bytes without
// exceeding memory allocation restrictions. This function is useful when doing
// a lot of hashing, because it allows for buffer re-use.
func NewTxHashIntoBuffer(version Version, selector Selector, input pack.Typed, data []byte) (id.Hash, error) {
	var err error
	buf := data
	rem := surge.MaxBytes
	if buf, rem, err = version.Marshal(buf, rem); err != nil {
		return id.Hash{}, err
	}
	if buf, rem, err = selector.Marshal(buf, rem); err != nil {
		return id.Hash{}, err
	}
	if buf, rem, err = input.Marshal(buf, rem); err != nil {
		return id.Hash{}, err
	}
	return id.NewHash(data), nil
}

// NewTx returns a transaction with the given recipient and inputs. The hash of
// the transaction is automatically computed and stored in the transaction. An
// error is returned when the recipient and inputs is too large and cannot be
// marshaled into bytes without exceeding memory allocation restrictions.
func NewTx(selector Selector, input pack.Typed) (Tx, error) {
	hash, err := NewTxHash(Version1, selector, input)
	if err != nil {
		return Tx{}, err
	}
	return Tx{Version: Version1, Hash: hash, Selector: selector, Input: input}, nil
}

// Generate allows us to quickly generate random transactions. This is mostly
// used for writing tests.
func (tx Tx) Generate(r *rand.Rand, size int) reflect.Value {
	hash, _ := quick.Value(reflect.TypeOf(id.Hash{}), r)
	version, _ := quick.Value(reflect.TypeOf(Version("")), r)
	selector, _ := quick.Value(reflect.TypeOf(Selector("")), r)
	input, _ := quick.Value(reflect.TypeOf(pack.Typed{}), r)
	output, _ := quick.Value(reflect.TypeOf(pack.Typed{}), r)
	return reflect.ValueOf(Tx{
		Hash:     hash.Interface().(id.Hash),
		Version:  version.Interface().(Version),
		Selector: selector.Interface().(Selector),
		Input:    input.Interface().(pack.Typed),
		Output:   output.Interface().(pack.Typed),
	})
}

// MapToPtrs maps a slice of transactions to a slice of transaction pointers.
func MapToPtrs(txs []Tx) []*Tx {
	txPtrs := make([]*Tx, len(txs))
	for i := range txs {
		txPtrs[i] = &txs[i]
	}
	return txPtrs
}

// MapFromPtrs maps a slice of transaction pointers to a slice of transactions.
func MapFromPtrs(txPtrs []*Tx) []Tx {
	txs := make([]Tx, len(txPtrs))
	for i := range txPtrs {
		txs[i] = *txPtrs[i]
	}
	return txs
}

// MapToHashes maps a slice of transactions to a slice of transaction hashes.
func MapToHashes(txs []Tx) []pack.Bytes32 {
	txHashes := make([]pack.Bytes32, len(txs))
	for i := range txs {
		txHashes[i] = pack.Bytes32(txs[i].Hash)
	}
	return txHashes
}

// MapToIDs maps a slice of transactions to a slice of transaction hashes.
func MapToIDs(txs []Tx) []id.Hash {
	txHashes := make([]id.Hash, len(txs))
	for i := range txs {
		txHashes[i] = id.Hash(txs[i].Hash)
	}
	return txHashes
}

// MapToNoStatus maps a slice of transactions with statuses to a slice of just
// transactions.
func MapToNoStatus(txsWithStatus []WithStatus) []Tx {
	txs := make([]Tx, len(txsWithStatus))
	for i := range txs {
		txs[i] = txsWithStatus[i].Tx
	}
	return txs
}
