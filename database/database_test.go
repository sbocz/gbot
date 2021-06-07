package database_test

import (
	"encoding/json"
	"gbot/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestStructure struct {
	StringField        string
	IntegerField       int
	FloatField         float64
	StringPointerField *string
}

var _ = Describe("Database", func() {
	var (
		testKey             = "Key"
		unwrittenKey        = "unwrittenKey"
		testString          = "FakeString"
		testInt             = 123
		testFloat           = 1.3
		testBucketType      = database.BucketType("testBucket")
		testEmptyBucketType = database.BucketType("emptyBucket")
		testKeysBucketType  = database.BucketType("keysBucket")
		databaseFile        = "testing.db"

		testStruct = &TestStructure{
			StringField:        testString,
			IntegerField:       testInt,
			FloatField:         testFloat,
			StringPointerField: &testString,
		}

		sut *database.DB
	)
	sut = database.NewDb(databaseFile)

	Context("Initializing a new database", func() {
		It("should return a database", func() {
			testDb := database.NewDb("foo.db")
			Expect(testDb).NotTo(BeNil())
		})
	})

	Context("CreateBucketIfNotExists", func() {
		It("should not error", func() {
			err := sut.CreateBucketIfNotExists(testBucketType)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Put", func() {
		It("should not error on write", func() {
			err := sut.Put(testBucketType, testKey, testStruct)
			Expect(err).To(BeNil())
		})

		It("should not error on read", func() {
			sut.Put(testBucketType, testKey, testStruct)
			_, err := sut.Get(testBucketType, testKey)
			Expect(err).To(BeNil())
		})

		It("should not error being unmarshalled into the type written", func() {
			sut.Put(testBucketType, testKey, testStruct)
			rawValue, _ := sut.Get(testBucketType, testKey)
			var result *TestStructure
			err := json.Unmarshal(rawValue, &result)

			Expect(err).To(BeNil())
		})

		Context("the read value should have the expected data", func() {
			sut.Put(testBucketType, testKey, testStruct)
			rawValue, _ := sut.Get(testBucketType, testKey)
			var result *TestStructure
			json.Unmarshal(rawValue, &result)

			Expect(result.StringField).To(Equal(testString))
			Expect(result.IntegerField).To(Equal(testInt))
			Expect(result.FloatField).To(Equal(testFloat))
			Expect(*result.StringPointerField).To(Equal(testString))
		})
	})

	Context("Get", func() {
		When("the key exists", func() {
			BeforeEach(func() {
				sut.Put(testBucketType, testKey, testStruct)
			})
			It("should not error on read", func() {
				_, err := sut.Get(testBucketType, testKey)
				Expect(err).To(BeNil())
			})
			It("should not error being unmarshalled into the type written", func() {
				rawValue, _ := sut.Get(testBucketType, testKey)
				var result *TestStructure
				err := json.Unmarshal(rawValue, &result)

				Expect(err).To(BeNil())
			})
			Context("the read value should have the expected data", func() {
				rawValue, _ := sut.Get(testBucketType, testKey)
				var result *TestStructure
				json.Unmarshal(rawValue, &result)

				Expect(result.StringField).To(Equal(testString))
				Expect(result.IntegerField).To(Equal(testInt))
				Expect(result.FloatField).To(Equal(testFloat))
				Expect(*result.StringPointerField).To(Equal(testString))
			})
		})
		When("the key does not exist", func() {
			It("should not error on read", func() {
				_, err := sut.Get(testBucketType, unwrittenKey)
				Expect(err).To(BeNil())
			})

			It("should return nil", func() {
				result, _ := sut.Get(testBucketType, unwrittenKey)
				Expect(result).To(BeNil())
			})
		})
	})

	Context("Keys", func() {
		When("there are no keys", func() {
			BeforeEach(func() {
				sut.CreateBucketIfNotExists(testEmptyBucketType)
			})
			It("should return an empty slice", func() {
				keys, err := sut.Keys(testEmptyBucketType)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(keys)).To(BeZero())
			})
		})
		When("there are multiple keys", func() {
			BeforeEach(func() {
				sut.CreateBucketIfNotExists(testKeysBucketType)
				sut.Put(testKeysBucketType, "foo", "bar")
				sut.Put(testKeysBucketType, "baz", "qux")
			})
			It("should return all keys", func() {
				keys, err := sut.Keys(testKeysBucketType)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(keys)).To(Equal(2))
				Expect(keys).To(ContainElement([]byte("foo")))
				Expect(keys).To(ContainElement([]byte("baz")))
			})
		})
	})

	Context("PutRandom", func() {
		It("should not error on write", func() {
			err := sut.PutRandom(testBucketType, testStruct)

			Expect(err).To(BeNil())
		})
		It("should not error being unmarshalled into the type written", func() {
			rawValue, _ := sut.GetRandom(testBucketType)
			var result *TestStructure
			err := json.Unmarshal(rawValue, &result)

			Expect(err).To(BeNil())
		})
	})

	Context("GetRandom", func() {
		BeforeEach(func() {
			sut.PutRandom(testBucketType, testStruct)
		})
		It("should not error on read", func() {
			sut.GetRandom(testBucketType)
			_, err := sut.Get(testBucketType, testKey)

			Expect(err).To(BeNil())
		})

		It("should not error being unmarshalled into the type written", func() {
			rawValue, _ := sut.GetRandom(testBucketType)
			var result *TestStructure
			err := json.Unmarshal(rawValue, &result)

			Expect(err).To(BeNil())
		})

		Context("the read value should have the expected data", func() {
			rawValue, _ := sut.GetRandom(testBucketType)
			var result *TestStructure
			json.Unmarshal(rawValue, &result)

			Expect(result.StringField).To(Equal(testString))
			Expect(result.IntegerField).To(Equal(testInt))
			Expect(result.FloatField).To(Equal(testFloat))
			Expect(*result.StringPointerField).To(Equal(testString))
		})
	})
})
