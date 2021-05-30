package banking_test

import (
	"gbot/commerce/banking"
	"gbot/database"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bank", func() {
	var (
		testUserId            = discord.UserID(1234)
		testNonExistantUserId = discord.UserID(12345)
		testDb                = database.NewDb("testing.db")

		sut *banking.Bank
	)

	Context("Creating a new Bank", func() {
		It("should not error when DB provided", func() {
			sut, err := banking.NewBank(testDb)
			Expect(sut).ToNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should error if DB is nil", func() {
			_, err := banking.NewBank(nil)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("Balance", func() {
		BeforeEach(func() {
			sut, _ = banking.NewBank(testDb)
		})
		When("value is positive", func() {
			BeforeEach(func() {
				sut.SaveAccount(banking.NewAccount(testUserId, 5, time.Now()))
			})
			It("should store the correct value", func() {
				balance, err := sut.Balance(testUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(balance).To(Equal(5))
			})
		})
		When("value is 0", func() {
			BeforeEach(func() {
				sut.SaveAccount(banking.NewAccount(testUserId, 0, time.Now()))
			})
			It("should return 0", func() {
				balance, err := sut.Balance(testUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(balance).To(BeZero())
			})
		})
		When("account does not exist", func() {
			It("should return 0", func() {
				balance, err := sut.Balance(testNonExistantUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(balance).To(BeZero())
			})
		})
	})

	Context("Deposit", func() {
		BeforeEach(func() {
			sut, _ = banking.NewBank(testDb)
		})
		When("value is positive", func() {
			AfterEach(func() {
				sut.Withdraw(testUserId, 20)
			})
			It("should store the correct value", func() {
				err := sut.Deposit(testUserId, 20)
				balance, _ := sut.Balance(testUserId)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(balance).To(Equal(20))
			})
		})
		When("value is negative", func() {
			It("should error", func() {
				err := sut.Deposit(testUserId, -20)
				Expect(err).Should(HaveOccurred())
			})
			It("should retain the old balance", func() {
				oldBalance, _ := sut.Balance(testUserId)
				sut.Deposit(testUserId, -20)
				newBalance, _ := sut.Balance(testUserId)

				Expect(newBalance).To(Equal(oldBalance))
			})
		})
	})

	Context("Withdraw", func() {
		BeforeEach(func() {
			sut, _ = banking.NewBank(testDb)
		})
		When("value is positive", func() {
			When("user has sufficient funds", func() {
				BeforeEach(func() {
					sut.Deposit(testUserId, 100)
				})
				AfterEach(func() {
					balance, _ := sut.Balance(testUserId)
					sut.Withdraw(testUserId, balance)
				})
				It("should return the correct value", func() {
					err := sut.Withdraw(testUserId, 50)
					balance, _ := sut.Balance(testUserId)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(balance).To(Equal(100 - 50))
				})
			})
			When("user does not have sufficient funds", func() {
				BeforeEach(func() {
					balance, _ := sut.Balance(testUserId)
					sut.Withdraw(testUserId, balance)
				})
				It("should error", func() {
					err := sut.Withdraw(testUserId, 1000)
					Expect(err).Should(HaveOccurred())
				})
				It("should retain the old balance", func() {
					oldBalance, _ := sut.Balance(testUserId)
					sut.Withdraw(testUserId, 1000)
					newBalance, _ := sut.Balance(testUserId)

					Expect(newBalance).To(Equal(oldBalance))
				})
			})
		})
		When("value is negative", func() {
			It("should error", func() {
				err := sut.Withdraw(testUserId, -1)
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("PayInterest", func() {
		BeforeEach(func() {
			sut, _ = banking.NewBank(testDb)
		})

		When("account is new", func() {
			BeforeEach(func() {
				account := banking.NewAccount(testUserId, 0, time.Now().UTC())
				sut.SaveAccount(account)
			})
			It("should update account last interest date and not change the balance", func() {
				err := sut.PayInterest(testUserId)
				account, _ := sut.FetchAccount(testUserId)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(account.Balance).To(Equal(0))
				Expect(account.CreationDate).Should(BeTemporally("<", account.LastInterest))
				Expect(account.LastInterest).Should(BeTemporally("~", time.Now().UTC(), time.Second))
			})
		})
		When("account is old", func() {
			BeforeEach(func() {
				account := banking.NewAccount(testUserId, 0, time.Now().UTC().Add(-time.Hour*24))
				sut.SaveAccount(account)
			})
			It("should update account last interest date and increase the balance", func() {
				err := sut.PayInterest(testUserId)
				account, _ := sut.FetchAccount(testUserId)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(account.Balance == 0).NotTo(BeTrue())
				Expect(account.CreationDate).Should(BeTemporally("<", account.LastInterest))
				Expect(account.LastInterest).Should(BeTemporally("~", time.Now().UTC(), time.Second))
			})
		})
	})
})
