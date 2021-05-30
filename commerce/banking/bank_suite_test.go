package banking_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBank(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Banking Suite")
}
