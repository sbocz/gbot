package commerce_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommerce(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commerce Suite")
}
