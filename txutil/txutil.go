package txutil

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/pack"
	"github.com/renproject/tx"
)

//
// RANDOM
//

func RandomTx(r *rand.Rand) tx.Tx {
	switch r.Int() % 2 {
	case 0:
		return RandomGoodTx(r)
	default:
		return RandomBadTx(r)
	}
}

func RandomTxHash(r *rand.Rand) id.Hash {
	switch r.Int() % 2 {
	case 0:
		return RandomGoodTxHash(r)
	default:
		return RandomBadTxHash(r)
	}
}

func RandomTxStatus(r *rand.Rand) tx.Status {
	switch r.Int() % 4 {
	case 0:
		return tx.StatusConfirming
	case 1:
		return tx.StatusPending
	case 2:
		return tx.StatusExecuting
	default:
		return tx.StatusDone
	}
}

//
// GOOD
//

func RandomGoodTx(r *rand.Rand) tx.Tx {
	// Generate a random transaction selector.
	randomSelector := RandomGoodTxSelector(r)

	// Generate random transaction inputs.
	input := RandomTxInput(r)

	// Construct the transaction.
	transaction, err := tx.NewTx(randomSelector, input)
	if err != nil {
		panic(err)
	}
	transaction.Version = tx.Version1 // The inputs we generate are for version 1 transactions.
	transaction.Output = pack.NewTyped()

	// Compute the transaction hash correctly, based on the other randomly
	// generated fields.
	hash, err := tx.NewTxHash(transaction.Version, transaction.Selector, transaction.Input)
	if err != nil {
		panic(err)
	}
	transaction.Hash = hash

	return transaction
}

func RandomGoodTxs(r *rand.Rand, n int) []tx.Tx {
	txs := make([]tx.Tx, n)
	for i := range txs {
		txs[i] = RandomGoodTx(r)
	}
	return txs
}

func RandomGoodTxWithStatus(r *rand.Rand) tx.WithStatus {
	transaction := RandomGoodTx(r)
	return tx.WithStatus{
		Tx:     transaction,
		Status: RandomTxStatus(r),
	}
}

func RandomGoodTxsWithStatus(r *rand.Rand, n int) []tx.WithStatus {
	txs := RandomGoodTxs(r, n)
	txsWithStatus := make([]tx.WithStatus, n)
	for i := range txs {
		txsWithStatus[i] = tx.WithStatus{
			Tx:     txs[i],
			Status: RandomTxStatus(r),
		}
	}
	return txsWithStatus
}

func RandomGoodTxHash(r *rand.Rand) id.Hash {
	transaction := RandomGoodTx(r)
	return transaction.Hash
}

func RandomGoodTxVersion(r *rand.Rand) tx.Version {
	switch r.Int() % 4 {
	case 0:
		return ""
	case 1:
		return tx.Version0
	default:
		return tx.Version1
	}
}

func RandomGoodTxSelector(r *rand.Rand) tx.Selector {
	selectors := AllSelectors()
	return selectors[r.Intn(len(selectors))]
}

func RandomTxInput(r *rand.Rand) pack.Typed {
	input := pack.NewStruct(
		"txid", pack.Bytes{}.Generate(r, r.Int()%100).Interface().(pack.Bytes),
		"txindex", pack.U32(0).Generate(r, 1).Interface().(pack.U32),
		"amount", pack.U256{}.Generate(r, 1).Interface().(pack.U256),
		"payload", pack.Bytes{}.Generate(r, r.Int()%100).Interface().(pack.Bytes),
		"phash", pack.Bytes32{}.Generate(r, 1).Interface().(pack.Bytes32),
		"to", pack.String("").Generate(r, r.Int()%100).Interface().(pack.String),
		"nonce", pack.Bytes32{}.Generate(r, 1).Interface().(pack.Bytes32),
		"nhash", pack.Bytes32{}.Generate(r, 1).Interface().(pack.Bytes32),
		"gpubkey", pack.Bytes{}.Generate(r, r.Int()%100).Interface().(pack.Bytes),
		"ghash", pack.Bytes32{}.Generate(r, 1).Interface().(pack.Bytes32),
	)
	return pack.Typed(input)
}

//
// BAD
//

func RandomBadTx(r *rand.Rand) tx.Tx {
	transaction := tx.Tx{}
	hash, _ := pack.Bytes32{}.Generate(r, 1).Interface().(pack.Bytes32)
	transaction.Hash = id.Hash(hash)
	return transaction
}

func RandomBadTxs(r *rand.Rand, n int) []tx.Tx {
	txs := make([]tx.Tx, n)
	for i := range txs {
		txs[i] = RandomBadTx(r)
	}
	return txs
}

func RandomBadTxHash(r *rand.Rand) id.Hash {
	return id.Hash{}
}

func RandomBadTxSelector(r *rand.Rand) tx.Selector {
	switch r.Int() % 2 {
	case 0:
		return tx.Selector("")
	default:
		str := pack.String("").Generate(r, r.Int()%100).Interface().(pack.String)
		return tx.Selector(strings.Replace(str.String(), "/", "", -1))
	}
}

//
// MISC
//

func TxsToTxsWithStatus(txs []tx.Tx, txStatus tx.Status) []tx.WithStatus {
	txsWithStatus := make([]tx.WithStatus, len(txs))
	for i := range txsWithStatus {
		txsWithStatus[i] = tx.WithStatus{
			Tx:     txs[i],
			Status: txStatus,
		}
	}
	return txsWithStatus
}

func TxsWithStatusToTxs(txsWithStatus []tx.WithStatus) []tx.Tx {
	txs := make([]tx.Tx, len(txsWithStatus))
	for i := range txs {
		txs[i] = txsWithStatus[i].Tx
	}
	return txs
}

func SupportedAssets() []multichain.Asset {
	return []multichain.Asset{
		multichain.BCH, multichain.BTC, multichain.DGB, multichain.DOGE,
		multichain.FIL, multichain.LUNA, multichain.ZEC,
	}
}

func SupportedHostChains() []multichain.Chain {
	return []multichain.Chain{
		multichain.Arbitrum, multichain.Avalanche, multichain.BinanceSmartChain,
		multichain.Ethereum, multichain.Fantom, multichain.Goerli,
		multichain.Moonbeam, multichain.Polygon, multichain.Solana,
	}
}

func AllSelectors() []tx.Selector {
	assets := SupportedAssets()
	hosts := SupportedHostChains()
	selectors := make([]tx.Selector, 0, 2*len(assets)*len(hosts))
	for _, asset := range assets {
		for _, host := range hosts {
			host := host
			selectors = append(selectors, tx.Selector(fmt.Sprintf("%v/to%v", asset, host)))
			selectors = append(selectors, tx.Selector(fmt.Sprintf("%v/from%v", asset, host)))

			for _, otherHost := range hosts {
				otherHost := otherHost
				if otherHost == host {
					continue
				}
				selectors = append(selectors, tx.Selector(fmt.Sprintf("%v/to%vFrom%v", asset, host, otherHost)))
			}
		}
	}
	return selectors
}

func SelectorsByAsset() map[multichain.Asset][]tx.Selector {
	assets := SupportedAssets()
	hosts := SupportedHostChains()
	selectors := make(map[multichain.Asset][]tx.Selector)

	// Populate all selectors by specific asset.
	for _, asset := range assets {
		selectors[asset] = make([]tx.Selector, 0, 2*len(hosts))
		for _, host := range hosts {
			host := host
			selectors[asset] = append(selectors[asset], tx.Selector(fmt.Sprintf("%v/to%v", asset, host)))
			selectors[asset] = append(selectors[asset], tx.Selector(fmt.Sprintf("%v/from%v", asset, host)))
			for _, otherHost := range hosts {
				otherHost := otherHost
				if otherHost == host {
					continue
				}
				selectors[asset] = append(selectors[asset], tx.Selector(fmt.Sprintf("%v/to%vFrom%v", asset, host, otherHost)))
			}
		}
	}
	return selectors
}

func SelectorsByDestination() map[multichain.Chain][]tx.Selector {
	assets := SupportedAssets()
	hosts := SupportedHostChains()
	selectors := make(map[multichain.Chain][]tx.Selector)

	// Populate asset/toHost selectors.
	for _, host := range hosts {
		selectors[host] = make([]tx.Selector, 0, len(assets))
		for _, asset := range assets {
			selectors[host] = append(selectors[host], tx.Selector(fmt.Sprintf("%v/to%v", asset, host)))
		}
	}

	// Populate asset/fromHost selectors.
	for _, asset := range assets {
		assetChain := asset.OriginChain()
		selectors[assetChain] = make([]tx.Selector, 0, len(hosts))
		for _, host := range hosts {
			selectors[assetChain] = append(selectors[assetChain], tx.Selector(fmt.Sprintf("%v/from%v", asset, host)))
		}
	}

	// Populate asset/toHostFromOtherHost selectors.
	for _, host := range hosts {
		for _, otherHost := range hosts {
			for _, asset := range assets {
				selectors[host] = append(selectors[host], tx.Selector(fmt.Sprintf("%v/to%vFrom%v", asset, host, otherHost)))
			}
		}
	}

	return selectors
}

func SelectorsBySource() map[multichain.Chain][]tx.Selector {
	assets := SupportedAssets()
	hosts := SupportedHostChains()
	selectors := make(map[multichain.Chain][]tx.Selector)

	// Populate asset/toHost selectors.
	for _, asset := range assets {
		assetChain := asset.OriginChain()
		selectors[assetChain] = make([]tx.Selector, 0, len(hosts))
		for _, host := range hosts {
			selectors[assetChain] = append(selectors[assetChain], tx.Selector(fmt.Sprintf("%v/to%v", asset, host)))
		}
	}

	// Populate asset/fromHost selectors.
	for _, host := range hosts {
		selectors[host] = make([]tx.Selector, len(assets))
		for i, asset := range assets {
			selectors[host][i] = tx.Selector(fmt.Sprintf("%v/from%v", asset, host))
		}
	}

	// Populate asset/toHostFromOtherHost selectors.
	for _, host := range hosts {
		for _, otherHost := range hosts {
			for _, asset := range assets {
				selectors[host] = append(selectors[host], tx.Selector(fmt.Sprintf("%v/to%vFrom%v", asset, otherHost, host)))
			}
		}
	}

	return selectors
}

func LockSelectors() map[multichain.Chain][]tx.Selector {
	assets := SupportedAssets()
	hosts := SupportedHostChains()
	selectors := make(map[multichain.Chain][]tx.Selector)

	// Populate asset/toHost selectors.
	for _, asset := range assets {
		assetChain := asset.OriginChain()
		selectors[assetChain] = make([]tx.Selector, 0, len(hosts))
		for _, host := range hosts {
			selectors[assetChain] = append(selectors[assetChain], tx.Selector(fmt.Sprintf("%v/to%v", asset, host)))
		}
	}

	return selectors
}

func NonLockSelectors() map[multichain.Chain][]tx.Selector {
	assets := SupportedAssets()
	hosts := SupportedHostChains()
	selectors := make(map[multichain.Chain][]tx.Selector)

	// Populate asset/fromHost selectors.
	for _, asset := range assets {
		assetChain := asset.OriginChain()
		for _, host := range hosts {
			selectors[assetChain] = append(selectors[assetChain], tx.Selector(fmt.Sprintf("%v/from%v", asset, host)))
		}
	}

	// Populate asset/toHostFromOtherHost selectors.
	for _, asset := range assets {
		assetChain := asset.OriginChain()
		for _, host := range hosts {
			host := host
			for _, otherHost := range hosts {
				otherHost := otherHost
				if otherHost == host {
					continue
				}
				selectors[assetChain] = append(selectors[assetChain], tx.Selector(fmt.Sprintf("%v/to%vFrom%v", asset, host, otherHost)))
			}
		}
	}

	return selectors
}
