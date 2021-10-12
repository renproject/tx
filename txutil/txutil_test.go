package txutil_test

import (
	"math/rand"
	"reflect"
	"testing/quick"

	"github.com/renproject/tx"
	"github.com/renproject/tx/txutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction Utility", func() {
	Context("when generating random txs with different seeds", func() {
		It("should not generate the same tx", func() {
			f := func(seed int64) bool {
				x := txutil.RandomTx(rand.New(rand.NewSource(seed)))
				y := txutil.RandomTx(rand.New(rand.NewSource(seed + 1)))
				return !reflect.DeepEqual(x, y)
			}
			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when generating random txs with the same seed", func() {
		It("should generate the same tx", func() {
			f := func(seed int64) bool {
				r := rand.New(rand.NewSource(seed))
				x := txutil.RandomTx(r)
				r = rand.New(rand.NewSource(seed))
				y := txutil.RandomTx(r)
				return reflect.DeepEqual(x, y)
			}
			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when generating random good txs", func() {
		It("should not generate the same tx", func() {
			f := func(seed int64) bool {
				x := txutil.RandomGoodTx(rand.New(rand.NewSource(seed)))
				y := txutil.RandomGoodTx(rand.New(rand.NewSource(seed + 1)))

				// Ensure hashes are as expected.
				xHash, err := tx.NewTxHash(x.Version, x.Selector, x.Input)
				Expect(err).ToNot(HaveOccurred())
				Expect(x.Hash).To(Equal(xHash))

				yHash, err := tx.NewTxHash(y.Version, y.Selector, y.Input)
				Expect(err).ToNot(HaveOccurred())
				Expect(y.Hash).To(Equal(yHash))

				return !reflect.DeepEqual(x, y)
			}
			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when generating random bad txs", func() {
		It("should not generate the same tx", func() {
			f := func(seed int64) bool {
				x := txutil.RandomBadTx(rand.New(rand.NewSource(seed)))
				y := txutil.RandomBadTx(rand.New(rand.NewSource(seed + 1)))

				// Ensure hashes are invalid.
				xHash, err := tx.NewTxHash(x.Version, x.Selector, x.Input)
				Expect(err).ToNot(HaveOccurred())
				Expect(x.Hash).ToNot(Equal(xHash))

				yHash, err := tx.NewTxHash(y.Version, y.Selector, y.Input)
				Expect(err).ToNot(HaveOccurred())
				Expect(y.Hash).ToNot(Equal(yHash))

				return !reflect.DeepEqual(x, y)
			}
			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
