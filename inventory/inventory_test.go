package inventory_test

import (
	"gbot/database"
	"gbot/inventory"

	"github.com/diamondburned/arikawa/v2/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inventory", func() {
	var (
		testDb         = database.NewDb("testing.db")
		existingItemId = "sweetapple"
		missingItemId  = "foo"
		testUserId     = discord.UserID(1234)
		missingUserId  = discord.UserID(4321)
		sut            *inventory.Inventory
	)

	Context("Creating a new Inventory", func() {
		It("should not error or be nil", func() {
			sut, err := inventory.NewInventory(testDb)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sut).ToNot(BeNil())
		})
	})

	Context("FetchItem", func() {
		BeforeEach(func() {
			sut, _ = inventory.NewInventory(testDb)
		})
		When("item exists", func() {
			It("should return the item", func() {
				item := sut.FetchItem(existingItemId)
				Expect(item).ToNot(BeNil())
				Expect(item.Identifier()).To(Equal(existingItemId))
			})
		})
		When("item does not exist", func() {
			It("should return nil", func() {
				item := sut.FetchItem(missingItemId)
				Expect(item).To(BeNil())
			})
		})
	})

	Context("ItemSet", func() {
		BeforeEach(func() {
			sut, _ = inventory.NewInventory(testDb)
		})
		It("should not be nil", func() {
			result := sut.ItemSet()
			Expect(result).NotTo(BeNil())
		})
		It("should have 2 items", func() {
			result := sut.ItemSet()
			Expect(len(result)).To(Equal(2))
		})
	})

	Context("FetchUserInventory", func() {
		BeforeEach(func() {
			sut, _ = inventory.NewInventory(testDb)
			user := inventory.NewUserInventory(testUserId)
			user.AddItem("banana")
			sut.WriteUserInventory(user)
		})
		When("user exists with items", func() {
			It("should not be nil", func() {
				user, err := sut.FetchUserInventory(testUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ToNot(BeNil())
			})
			It("should not have 0 items", func() {
				user, err := sut.FetchUserInventory(testUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(user.ItemMap)).ToNot(BeZero())
			})
		})
		When("user does not exist", func() {
			It("should be a new defaulted entry", func() {
				user, err := sut.FetchUserInventory(missingUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(user.ItemMap)).To(BeZero())
			})
		})
	})

	Context("WriteUserInventory", func() {
		BeforeEach(func() {
			sut, _ = inventory.NewInventory(testDb)
		})
		It("should not error", func() {
			user := inventory.NewUserInventory(testUserId)
			user.AddItem("banana")
			err := sut.WriteUserInventory(user)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should overwrite old data", func() {
			user := inventory.NewUserInventory(testUserId)
			user.AddItem("banana")
			err := sut.WriteUserInventory(user)
			Expect(err).ShouldNot(HaveOccurred())

			oldUser, _ := sut.FetchUserInventory(testUserId)

			err = sut.WriteUserInventory(inventory.NewUserInventory(testUserId))
			Expect(err).ShouldNot(HaveOccurred())

			newUser, _ := sut.FetchUserInventory(testUserId)

			Expect(oldUser.UserId).To(Equal(newUser.UserId))
			Expect(len(oldUser.ItemMap)).NotTo(Equal(len(newUser.ItemMap)))
		})
	})
})
