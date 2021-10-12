package tx_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	"github.com/renproject/multichain"
	"github.com/renproject/pack/packutil"
	"github.com/renproject/surge/surgeutil"
	"github.com/renproject/tx"
	"github.com/renproject/tx/txutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/tx"
)

var _ = Describe("Selector", func() {

	t := reflect.TypeOf(Selector(""))
	numTrials := 100

	allAssets := txutil.SupportedAssets()
	allHosts := txutil.SupportedHostChains()

	selectorsByAsset := txutil.SelectorsByAsset()
	selectorsByDestination := txutil.SelectorsByDestination()
	selectorsBySource := txutil.SelectorsBySource()

	lockSelectors := txutil.LockSelectors()
	nonLockSelectors := txutil.NonLockSelectors()

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

	Context("when fetching the destination", func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		Context("if the selector is valid", func() {
			It("should return the correct destination chain", func() {
				loop := func() bool {
					selector := txutil.RandomGoodTxSelector(r)

					// Its either a lock and mint tx.
					for _, host := range allHosts {
						if contains(selectorsByDestination[host], selector) {
							Expect(selector.Destination()).To(Equal(host))
							return true
						}
					}
					// Or a burn and release tx.
					for _, asset := range allAssets {
						assetChain := asset.OriginChain()
						if contains(selectorsByDestination[assetChain], selector) {
							Expect(selector.Destination()).To(Equal(assetChain))
							return true
						}
					}

					Fail(fmt.Sprintf("destination for selector=%v not found", selector))
					return false
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})

		Context("if the selector is invalid", func() {
			It("should return an empty string", func() {
				loop := func() bool {
					selector := txutil.RandomBadTxSelector(r)
					Expect(selector.Destination()).To(Equal(multichain.Chain("")))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})
	})

	Context("when fetching the source", func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		Context("if the selector is valid", func() {
			It("should return the correct source chain", func() {
				loop := func() bool {
					selector := txutil.RandomGoodTxSelector(r)

					// Its either a lock and mint tx.
					for _, asset := range allAssets {
						assetChain := asset.OriginChain()
						if contains(selectorsBySource[assetChain], selector) {
							Expect(selector.Source()).To(Equal(assetChain))
							return true
						}
					}
					// Or a burn and release tx.
					for _, host := range allHosts {
						if contains(selectorsBySource[host], selector) {
							Expect(selector.Source()).To(Equal(host))
							return true
						}
					}

					Fail(fmt.Sprintf("source for selector=%v not found", selector))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})

		Context("if the selector is invalid", func() {
			It("should return an empty string", func() {
				loop := func() bool {
					selector := txutil.RandomBadTxSelector(r)
					Expect(selector.Source()).To(Equal(multichain.Chain("")))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})
	})

	Context("when fetching the asset", func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		Context("if the selector is valid", func() {
			It("should return the correct value", func() {
				loop := func() bool {
					selector := txutil.RandomGoodTxSelector(r)

					for _, asset := range allAssets {
						if contains(selectorsByAsset[asset], selector) {
							Expect(selector.Asset()).To(Equal(asset))
							return true
						}
					}

					Fail(fmt.Sprintf("asset for selector=%v not found", selector))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})

		Context("if the selector is invalid", func() {
			It("should return an empty string", func() {
				loop := func() bool {
					selector := txutil.RandomBadTxSelector(r)
					Expect(selector.Asset()).To(Equal(multichain.Asset("")))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})
	})

	Context("when fetching the lock chain", func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		Context("if the selector is valid", func() {
			It("should return the correct value", func() {
				loop := func() bool {
					selector := txutil.RandomGoodTxSelector(r)

					for _, asset := range allAssets {
						assetChain := asset.OriginChain()
						if contains(lockSelectors[assetChain], selector) {
							Expect(selector.IsLock()).To(BeTrue())
							Expect(selector.Source()).To(Equal(assetChain))
							return true
						}
						if contains(nonLockSelectors[assetChain], selector) {
							Expect(selector.IsLock()).To(BeFalse())
							return true
						}
					}

					Fail(fmt.Sprintf("lockChain for selector=%v not found", selector))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())
			})
		})

		Context("if the selector is invalid", func() {
			It("should return an empty string", func() {
				loop := func() bool {
					selector := txutil.RandomBadTxSelector(r)
					ok := selector.IsLock()
					lockChain := selector.Source()
					Expect(ok).To(BeFalse())
					Expect(lockChain).To(Equal(multichain.Chain("")))
					return true
				}
				Expect(quick.Check(loop, nil)).To(Succeed())

				selector := tx.Selector("BTC/randomFn")
				ok := selector.IsLock()
				lockChain := selector.Source()
				Expect(ok).To(BeFalse())
				Expect(lockChain).To(Equal(multichain.Chain("")))
			})
		})
	})

	Context("when stringifying", func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		It("should return the expected string", func() {
			loop := func() bool {
				selector := txutil.RandomGoodTxSelector(r)
				Expect(selector.String()).To(Equal(string(selector)))
				return true
			}
			Expect(quick.Check(loop, nil)).To(Succeed())
		})
	})
})

func contains(s []tx.Selector, e tx.Selector) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
