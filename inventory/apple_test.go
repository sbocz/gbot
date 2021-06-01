package inventory_test

import (
	"gbot/inventory"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apple", func() {
	var (
		flavor string = "sweet"
	)
	Context("Creating a new Apple", func() {
		It("should not be nil", func() {
			sut := inventory.NewApple(flavor)
			Expect(sut).ToNot(BeNil())
		})
	})

	Context("Identifier", func() {
		It("should include the flavor", func() {
			sut := inventory.NewApple(flavor)
			Expect(sut.Identifier()).To(Equal(flavor + "apple"))
		})
	})

	Context("Description", func() {
		It("should include the flavor", func() {
			sut := inventory.NewApple(flavor)
			Expect(strings.Contains(sut.Description(), flavor)).To(BeTrue())
		})
	})

	Context("Rarity", func() {
		It("should be common", func() {
			sut := inventory.NewApple(flavor)
			Expect(sut.Rarity()).To(Equal(inventory.Common))
		})
	})

	Context("Use", func() {
		It("should have 5 outcomes", func() {
			sut := inventory.NewApple(flavor)
			results := make(map[string]int)
			for i := 0; i < 100; i++ {
				results[sut.Use()] += 1
			}
			Expect(len(results)).To(Equal(5))
		})
	})
})
