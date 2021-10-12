package tx_test

import (
	"fmt"

	"github.com/renproject/tx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction version", func() {

	table := []struct {
		version     tx.Version
		stringified string
	}{
		{tx.Version0, "0"},
		{tx.Version1, "1"},
		{tx.Version("2"), ""}, // Unknown versions
		{tx.Version("3"), ""}, // Unknown versions
		{tx.Version(""), ""},  // Unknown versions
	}

	for _, entry := range table {
		entry := entry

		Context(fmt.Sprintf("when stringifying version=%v", entry.version), func() {
			It(fmt.Sprintf("should return \"%v\"", entry.stringified), func() {
				Expect(entry.version.String()).To(Equal(entry.stringified))
			})
		})
	}
})
