package banking_test

import (
	"gbot/commerce/banking"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	var (
		testUserId          = discord.UserID(1234)
		testStartingBalance = 123
		approxTime          = time.Now().UTC()
		oneDayAgo           = approxTime.AddDate(0, 0, -1)
		periodMillis        = time.Hour
		periods             = 24 // one day
		interestRate        = 10
		expectedInterest    = periods * interestRate

		sut *banking.Account
	)

	Context("Creating a new Account", func() {
		It("should not error", func() {
			sut = banking.NewAccount(testUserId, testStartingBalance, time.Now().UTC())
		})
		It("should have the correct values", func() {
			creationDate := time.Now().UTC().Add(-time.Hour * 100000)
			sut = banking.NewAccount(testUserId, testStartingBalance, creationDate)
			Expect(sut.UserId).To(Equal(testUserId))
			Expect(sut.Balance).To(Equal(testStartingBalance))
			Expect(sut.InterestValue).To(Equal(10))
			Expect(sut.CreationDate).To(BeTemporally("~", creationDate, time.Second))
			Expect(sut.LastInterest).To(BeTemporally("~", creationDate, time.Second))
		})
	})

	Context("CalculateInterest", func() {
		When("interest was last paid one day ago", func() {
			BeforeEach(func() {
				sut = banking.NewAccount(testUserId, testStartingBalance, time.Now().UTC())
				sut.LastInterest = oneDayAgo
			})
			It("should not error", func() {
				sut.CalculateInterest(periodMillis)
			})
			It("should be the correct value", func() {
				result := sut.CalculateInterest(periodMillis)
				Expect(result).To(Equal(expectedInterest))
			})
		})
		When("interest is in the future", func() {
			BeforeEach(func() {
				sut = banking.NewAccount(testUserId, testStartingBalance, time.Now().UTC())
				sut.LastInterest = approxTime.AddDate(0, 0, 1)
			})
			It("should not error", func() {
				sut.CalculateInterest(periodMillis)
			})
			It("should return 0", func() {
				result := sut.CalculateInterest(periodMillis)
				Expect(result).To(BeZero())
			})
		})
		When("interest is in the past but less than the period", func() {
			BeforeEach(func() {
				sut = banking.NewAccount(testUserId, testStartingBalance, time.Now().UTC())
				sut.LastInterest = approxTime
			})
			It("should not error", func() {
				sut.CalculateInterest(periodMillis)
			})
			It("should return 0", func() {
				result := sut.CalculateInterest(periodMillis)
				Expect(result).To(BeZero())
			})
		})
	})
})
