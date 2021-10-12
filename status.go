package tx

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing/quick"

	"github.com/renproject/surge"
)

// Status of a transaction.
type Status uint8

const (
	// StatusNil is used for transactions that have an unknown status. For
	// example, invalid transactions that have been rejected.
	StatusNil = Status(0)

	// StatusConfirming is used for transactions that are waiting for their
	// underlying blockchain transactions to confirm. For example, before a
	// BTC0Btc2Eth transaction is accepted by a Darknode, its inner UTXO must
	// have the required confirmations on Bitcoin. This status is not directly
	// used within the Darknodes, but is used by the Lightnodes and other
	// components in the network.
	StatusConfirming = Status(1)

	// StatusPending is used for transactions that have been successfully
	// validated and are waiting to be included in a Hyperdrive block.
	StatusPending = Status(2)

	// StatusExecuting is used for transactions that have been included in a
	// Hyperdrive block, where the block has not been executed.
	StatusExecuting = Status(3)

	// StatusDone is used for transactions that have been included in a
	// Hyperdrive block, where the block has been executed. This is irrespective
	// of whether or not the transaction was reverted. It is worth noting that
	// cross-chain transaction cannot be reverted; they will instead be rejected
	// before reaching the "pending" status).
	StatusDone = Status(4)
)

func (status Status) String() string {
	switch status {
	case StatusNil:
		return "nil"
	case StatusConfirming:
		return "confirming"
	case StatusPending:
		return "pending"
	case StatusExecuting:
		return "executing"
	case StatusDone:
		return "done"
	default:
		return ""
	}
}

// MarshalText from the status.
func (status Status) MarshalText() ([]byte, error) {
	return []byte(status.String()), nil
}

// UnmarshalText into the status.
func (status *Status) UnmarshalText(data []byte) error {
	switch string(data) {
	case "nil":
		*status = StatusNil
	case "confirming":
		*status = StatusConfirming
	case "pending":
		*status = StatusPending
	case "executing":
		*status = StatusExecuting
	case "done":
		*status = StatusDone
	default:
		return fmt.Errorf("non-exhaustive pattern: status %v", string(data))
	}
	return nil
}

// Generate allows us to quickly generate random transaction statuses. This is
// mostly used for writing tests.
func (Status) Generate(r *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(Status(r.Intn(5)))
}

// WithStatus is a combination of a transaction and its current status. It is a
// helper struct for moving both values together.
type WithStatus struct {
	Tx     `json:"tx"`
	Status Status `json:"status"`
}

// SizeHint returns the number of bytes required to represent the WithStatus
// type in binary.
func (w WithStatus) SizeHint() int {
	return surge.SizeHint(w.Tx) + surge.SizeHintU8
}

// Marshal the WithStatus type to binary.
func (w WithStatus) Marshal(buf []byte, rem int) ([]byte, int, error) {
	var err error
	if buf, rem, err = surge.Marshal(w.Tx, buf, rem); err != nil {
		return buf, rem, err
	}
	return surge.MarshalU8(uint8(w.Status), buf, rem)
}

// Unmarshal the WithStatus type from binary.
func (w *WithStatus) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	var err error
	if buf, rem, err = surge.Unmarshal(&w.Tx, buf, rem); err != nil {
		return buf, rem, err
	}
	var status uint8
	if buf, rem, err = surge.UnmarshalU8(&status, buf, rem); err != nil {
		return buf, rem, err
	}
	w.Status = Status(status)
	return buf, rem, nil
}

// Generate allows us to quickly generate random transactions with statuses.
// This is mostly used for writing tests.
func (w WithStatus) Generate(r *rand.Rand, size int) reflect.Value {
	tx, _ := quick.Value(reflect.TypeOf(Tx{}), r)
	status, _ := quick.Value(reflect.TypeOf(Status(0)), r)
	return reflect.ValueOf(WithStatus{
		Tx:     tx.Interface().(Tx),
		Status: status.Interface().(Status),
	})
}

// ZipStatus to a slice of transactions, and return the zipped slice.
func ZipStatus(txs []Tx, status Status) []WithStatus {
	txsWithStatus := make([]WithStatus, len(txs))
	for i := range txs {
		txsWithStatus[i] = WithStatus{Tx: txs[i], Status: status}
	}
	return txsWithStatus
}
