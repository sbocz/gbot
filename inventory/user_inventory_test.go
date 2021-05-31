package inventory_test

import (
	"gbot/inventory"

	"github.com/diamondburned/arikawa/v2/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserInventory", func() {
	var (
		testUserId = discord.UserID(1234)

		sut *inventory.UserInventory
	)
	Context("Creating a new UserInventory", func() {
		It("should not be nil", func() {
			sut = inventory.NewUserInventory(testUserId)
			Expect(sut).ToNot(BeNil())
		})
	})

	Context("AddItem", func() {
		It("should set the value to 1 for new items", func() {
			sut = inventory.NewUserInventory(testUserId)
			sut.AddItem("foo")
			Expect(sut.ItemMap["foo"]).To(Equal(1))
		})
		It("should increment the entry for existing items", func() {
			sut = inventory.NewUserInventory(testUserId)
			sut.AddItem("bar")
			sut.AddItem("bar")
			Expect(sut.ItemMap["bar"]).To(Equal(2))
		})
	})

	Context("RemoveItem", func() {
		When("the item is not in the inventory", func() {
			It("should error", func() {
				sut = inventory.NewUserInventory(testUserId)
				err := sut.RemoveItem("foo")
				Expect(err).To(HaveOccurred())
			})
		})
		When("there are 3 of the item in the inventory", func() {
			BeforeEach(func() {
				sut = inventory.NewUserInventory(testUserId)
				sut.AddItem("foo")
				sut.AddItem("foo")
				sut.AddItem("foo")
			})
			It("should succeed 3 times", func() {
				err := sut.RemoveItem("foo")
				Expect(err).NotTo(HaveOccurred())
				err = sut.RemoveItem("foo")
				Expect(err).NotTo(HaveOccurred())
				err = sut.RemoveItem("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(sut.ItemMap["foo"]).To(BeZero())
			})
			It("should error on the 4th attempt", func() {
				sut.RemoveItem("foo")
				sut.RemoveItem("foo")
				sut.RemoveItem("foo")
				err := sut.RemoveItem("foo")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
