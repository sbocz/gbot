package shop_test

import (
	"fmt"
	"gbot/commerce/banking"
	"gbot/commerce/shop"
	"gbot/database"
	"gbot/inventory"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shop", func() {
	var (
		testDb           = database.NewDb("testing.db")
		testBank, _      = banking.NewBank(testDb)
		testInventory, _ = inventory.NewInventory(testDb)
		testRichUserId   = discord.UserID(1234)
		testPoorUserId   = discord.UserID(12345)
		testItemId       = "sweetapple"
		testItem         = testInventory.FetchItem(testItemId)
		startingBalance  = 1000

		sut *shop.Shop
	)

	Context("Creating a new Shop", func() {
		It("should not be nil", func() {
			sut = shop.NewShop(testBank, testInventory)
			Expect(sut).ToNot(BeNil())
		})
	})

	Context("BuyItem", func() {
		When("it succeeds", func() {
			BeforeEach(func() {
				testBank.SaveAccount(banking.NewAccount(testRichUserId, startingBalance, time.Now()))
				testInventory.WriteUserInventory(inventory.NewUserInventory(testRichUserId))
				sut = shop.NewShop(testBank, testInventory)
			})
			It("should withdraw the correct funds", func() {
				err := sut.BuyItem(testRichUserId, testItemId, 5)
				cost := testItem.Value() * 5
				balance, _ := testBank.Balance(testRichUserId)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(balance).To(Equal(startingBalance - cost))
			})
			It("should add the item to the user inventory", func() {
				err := sut.BuyItem(testRichUserId, testItemId, 5)
				userInventory, _ := testInventory.FetchUserInventory(testRichUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(userInventory.ItemMap[testItemId]).To(Equal(5))
			})
			It("should reduce the quantity in the shop", func() {
				startingCount := sut.Entries[testItemId].Count
				err := sut.BuyItem(testRichUserId, testItemId, 5)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(sut.Entries[testItemId].Count).To(Equal(startingCount - 5))
			})
		})
		When("the shop does not have the item", func() {
			BeforeEach(func() {
				testBank.SaveAccount(banking.NewAccount(testRichUserId, startingBalance, time.Now()))
				testInventory.WriteUserInventory(inventory.NewUserInventory(testRichUserId))
				sut = shop.NewShop(testBank, testInventory)
			})
			It("should error", func() {
				err := sut.BuyItem(testRichUserId, "fakeItem", 1)
				Expect(err).Should(HaveOccurred())
			})
			It("should not add the item to the user inventory", func() {
				err := sut.BuyItem(testRichUserId, "fakeItem", 1)
				userInventory, _ := testInventory.FetchUserInventory(testRichUserId)
				Expect(err).Should(HaveOccurred())
				Expect(userInventory.ItemMap["fakeItem"]).To(Equal(0))
			})
			It("should not change the quantity in the shop", func() {
				err := sut.BuyItem(testRichUserId, "fakeItem", 1)
				Expect(err).Should(HaveOccurred())
				Expect(sut.Entries["fakeItem"]).To(BeNil())
			})
		})
		When("the shop does not have enough quantity of the item", func() {
			BeforeEach(func() {
				testBank.SaveAccount(banking.NewAccount(testRichUserId, startingBalance, time.Now()))
				testInventory.WriteUserInventory(inventory.NewUserInventory(testRichUserId))
				sut = shop.NewShop(testBank, testInventory)
			})
			It("should error", func() {
				err := sut.BuyItem(testRichUserId, testItemId, 20)
				Expect(err).Should(HaveOccurred())
			})
			It("should not add the item to the user inventory", func() {
				err := sut.BuyItem(testRichUserId, testItemId, 20)
				userInventory, _ := testInventory.FetchUserInventory(testRichUserId)
				Expect(err).Should(HaveOccurred())
				Expect(userInventory.ItemMap[testItemId]).To(Equal(0))
			})
			It("should not change the quantity in the shop", func() {
				startingCount := sut.Entries[testItemId].Count
				err := sut.BuyItem(testRichUserId, testItemId, 20)
				Expect(err).Should(HaveOccurred())
				Expect(sut.Entries[testItemId].Count).To(Equal(startingCount))
			})
		})
		When("the user does not have enough funds", func() {
			BeforeEach(func() {
				testBank.SaveAccount(banking.NewAccount(testPoorUserId, 0, time.Now()))
				testInventory.WriteUserInventory(inventory.NewUserInventory(testPoorUserId))
				sut = shop.NewShop(testBank, testInventory)
			})
			It("should error", func() {
				err := sut.BuyItem(testPoorUserId, testItemId, 1)
				Expect(err).Should(HaveOccurred())
			})
			It("should not add the item to the user inventory", func() {
				err := sut.BuyItem(testPoorUserId, testItemId, 1)
				userInventory, _ := testInventory.FetchUserInventory(testPoorUserId)
				Expect(err).Should(HaveOccurred())
				Expect(userInventory.ItemMap[testItemId]).To(Equal(0))
			})
			It("should not withdraw funds", func() {
				startingBalance, _ := testBank.Balance(testPoorUserId)
				err := sut.BuyItem(testPoorUserId, testItemId, 1)
				Expect(err).Should(HaveOccurred())
				Expect(testBank.Balance(testPoorUserId)).To(Equal(startingBalance))
			})
			It("should not change the quantity in the shop", func() {
				startingCount := sut.Entries[testItemId].Count
				err := sut.BuyItem(testPoorUserId, testItemId, 1)
				Expect(err).Should(HaveOccurred())
				Expect(sut.Entries[testItemId].Count).To(Equal(startingCount))
			})
		})
	})

	Context("PrintInventory", func() {
		It("should not panic", func() {
			sut = shop.NewShop(testBank, testInventory)
			fmt.Println("------STARTING INVENTORY------")
			fmt.Println(sut.PrintInventory())
			fmt.Println("-------ENDING INVENTORY-------")
		})
	})
})
