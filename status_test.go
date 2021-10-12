package tx_test

import (
	"fmt"
	"reflect"

	"github.com/renproject/pack/packutil"
	"github.com/renproject/surge/surgeutil"
	"github.com/renproject/tx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction status", func() {

	t := reflect.TypeOf(tx.WithStatus{})
	numTrials := 50

	Context("when fuzzing", func() {
		It("should not panic", func() {
			for trial := 0; trial < numTrials; trial++ {
				Expect(func() { surgeutil.Fuzz(t) }).ToNot(Panic())
				Expect(func() { packutil.JSONFuzz(t) }).ToNot(Panic())
			}
		})
	})

	Context("when marshaling and then unmarshaling", func() {
		It("should return itself", func() {
			for trial := 0; trial < numTrials; trial++ {
				Expect(surgeutil.MarshalUnmarshalCheck(t)).To(Succeed())
				Expect(JSONMarshalUnmarshalCheck(t)).To(Succeed())
			}
		})
	})

	Context("when marshaling", func() {
		Context("when the buffer is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(surgeutil.MarshalBufTooSmall(t)).To(Succeed())
				}
			})
		})

		Context("when the remaining memory quota is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(surgeutil.MarshalRemTooSmall(t)).To(Succeed())
				}
			})
		})
	})

	Context("when unmarshaling", func() {
		Context("when the buffer is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(surgeutil.UnmarshalBufTooSmall(t)).To(Succeed())
				}
			})
		})

		Context("when the remaining memory quota is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(surgeutil.UnmarshalRemTooSmall(t)).To(Succeed())
				}
			})
		})
	})

	table := []struct {
		status      tx.Status
		stringified string
		value       uint8
	}{
		{tx.StatusNil, "nil", 0},
		{tx.StatusConfirming, "confirming", 1},
		{tx.StatusPending, "pending", 2},
		{tx.StatusExecuting, "executing", 3},
		{tx.StatusDone, "done", 4},
		{tx.Status(255), "", 255}, // Unknown status
	}

	for _, entry := range table {
		entry := entry

		Context(fmt.Sprintf("when stringifying status=%v", entry.status), func() {
			It(fmt.Sprintf("should return \"%v\"", entry.stringified), func() {
				Expect(entry.status.String()).To(Equal(entry.stringified))
			})
		})

		Context(fmt.Sprintf("when checking value of status=%v", entry.status), func() {
			It(fmt.Sprintf("should equal %v", entry.value), func() {
				Expect(uint8(entry.status)).To(Equal(entry.value))
			})
		})
	}
})
