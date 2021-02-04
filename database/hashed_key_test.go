package database_test

import (
	"gbot/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	Context("Generating a random key", func() {
		It("should not panic", func() {
			database.NewRandomHashedKey()
		})
		It("should be the correct size", func() {
			key := database.NewRandomHashedKey()
			Expect(len(key)).To(Equal(database.HashedKeySize))
		})
	})

	Context("Generating a key from a value", func() {
		It("should be the correct size", func() {
			key := database.NewHashedKey("foo")
			Expect(len(key)).To(Equal(database.HashedKeySize))
		})

		It("should be the same result for the same value", func() {
			Context("for string values", func() {
				key1 := database.NewHashedKey("foo")
				key2 := database.NewHashedKey("foo")

				Expect(key1).NotTo(BeNil())
				for index, value := range key1 {
					Expect(value).To(Equal(key2[index]))
				}
			})
			Context("for int values", func() {
				key1 := database.NewHashedKey(3)
				key2 := database.NewHashedKey(3)

				Expect(key1).NotTo(BeNil())
				for index, value := range key1 {
					Expect(value).To(Equal(key2[index]))
				}
			})
			Context("for nil value", func() {
				key1 := database.NewHashedKey(nil)
				key2 := database.NewHashedKey(nil)

				Expect(key1).NotTo(BeNil())
				for index, value := range key1 {
					Expect(value).To(Equal(key2[index]))
				}
			})
			Context("for struct values", func() {
				key1 := database.NewHashedKey(TestStructure{
					StringField:  "bar",
					IntegerField: 10,
				})
				key2 := database.NewHashedKey(TestStructure{
					StringField:  "bar",
					IntegerField: 10,
				})

				Expect(key1).NotTo(BeNil())
				for index, value := range key1 {
					Expect(value).To(Equal(key2[index]))
				}
			})
			Context("for struct pointers", func() {
				key1 := database.NewHashedKey(&TestStructure{
					StringField:  "baz",
					IntegerField: 13,
				})
				key2 := database.NewHashedKey(&TestStructure{
					StringField:  "baz",
					IntegerField: 13,
				})

				Expect(key1).NotTo(BeNil())
				for index, value := range key1 {
					Expect(value).To(Equal(key2[index]))
				}
			})
		})
	})
})
