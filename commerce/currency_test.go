package commerce_test

import (
	"gbot/commerce"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Currency", func() {
	var (
		intValue       = 123
		currencyValue  = commerce.Currency(intValue)
		expectedString = "123𝔻"
	)
	Context("String", func() {
		When("it is 3 digits", func() {
			It("should add a '𝔻'", func() {
				Expect(currencyValue.String()).To(Equal(expectedString))
			})
		})
		When("it is 4 digits", func() {
			BeforeEach(func() {
				intValue = 1234
				currencyValue = commerce.Currency(intValue)
				expectedString = "1,234𝔻"
			})
			It("should return the expected string", func() {
				Expect(currencyValue.String()).To(Equal(expectedString))
			})
		})
		When("it is 5 digits", func() {
			BeforeEach(func() {
				intValue = 12345
				currencyValue = commerce.Currency(intValue)
				expectedString = "12,345𝔻"
			})
			It("should return the expected string", func() {
				Expect(currencyValue.String()).To(Equal(expectedString))
			})
		})
		When("it is 6 digits", func() {
			BeforeEach(func() {
				intValue = 123456
				currencyValue = commerce.Currency(intValue)
				expectedString = "123,456𝔻"
			})
			It("should return the expected string", func() {
				Expect(currencyValue.String()).To(Equal(expectedString))
			})
		})
		When("it is 7 digits", func() {
			BeforeEach(func() {
				intValue = 1234567
				currencyValue = commerce.Currency(intValue)
				expectedString = "1,234,567𝔻"
			})
			It("should return the expected string", func() {
				Expect(currencyValue.String()).To(Equal(expectedString))
			})
		})
		When("it is a negative value", func() {
			BeforeEach(func() {
				intValue = -1234567
				currencyValue = commerce.Currency(intValue)
				expectedString = "-1,234,567𝔻"
			})
			It("should return the expected string", func() {
				Expect(currencyValue.String()).To(Equal(expectedString))
			})
		})
	})
})
