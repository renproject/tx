package tx_test

import (
	"math/rand"
	"reflect"
	"testing/quick"

	"github.com/renproject/surge"
	"github.com/renproject/surge/surgeutil"
	"github.com/renproject/tx"
	"github.com/renproject/tx/txutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transactions", func() {

	t := reflect.TypeOf(tx.Tx{})
	numTrials := 50

	Context("when fuzzing transactions", func() {
		It("should not panic", func() {
			Expect(func() { surgeutil.Fuzz(t) }).ToNot(Panic())
			Expect(func() { JSONFuzz(t) }).ToNot(Panic())
		})
	})

	Context("when marshaling and then unmarshaling transactions", func() {
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

	Context("when constructing a new tx hash", func() {
		It("should be correct", func() {
			f := func(seed int64) bool {
				r := rand.New(rand.NewSource(seed))
				transaction := txutil.RandomGoodTx(r)
				txHash, err := tx.NewTxHash(
					transaction.Version,
					transaction.Selector,
					transaction.Input,
				)
				Expect(err).ToNot(HaveOccurred())
				buf := make([]byte, surge.SizeHintString(
					string(transaction.Version))+
					surge.SizeHintString(string(transaction.Selector))+
					surge.SizeHint(transaction.Input),
				)
				expectedHash, err := tx.NewTxHashIntoBuffer(
					transaction.Version,
					transaction.Selector,
					transaction.Input,
					buf,
				)
				Expect(txHash).To(Equal(expectedHash))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})

		It("should be reconstructible", func() {
			f := func(seed int64) bool {
				r := rand.New(rand.NewSource(seed))
				transaction := txutil.RandomGoodTx(r)
				expectedHash, err := tx.NewTxHash(
					transaction.Version,
					transaction.Selector,
					transaction.Input,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(transaction.Hash).To(Equal(expectedHash))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})

		It("should be unique", func() {
			f := func(seed int64) bool {
				r := rand.New(rand.NewSource(seed))
				fstHash := txutil.RandomGoodTxHash(r)
				sndHash := txutil.RandomGoodTxHash(r)
				Expect(fstHash).ToNot(Equal(sndHash))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})
