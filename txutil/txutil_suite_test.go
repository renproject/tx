package txutil_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTxUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transaction Utility Suite")
}
