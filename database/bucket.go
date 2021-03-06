package database

type BucketType string

const (
	BankAccounts BucketType = "BankAccounts"
	StockData    BucketType = "StockData"
	UserNotes    BucketType = "UserNotes"
	ShoutPhrases BucketType = "ShoutPhrases"
	Tests        BucketType = "Tests"
)

type Bucket struct {
	db    *DB
	bType BucketType
}

func NewBucket(db *DB, bucketType BucketType) *Bucket {
	return &Bucket{db: db, bType: bucketType}
}

func (b *Bucket) Put(key string, value interface{}) error {
	return b.db.Put(b.bType, key, value)
}

func (b *Bucket) Get(key string) ([]byte, error) {
	return b.db.Get(b.bType, key)
}

func (b *Bucket) PutRandom(value interface{}) error {
	return b.db.PutRandom(b.bType, value)
}

func (b *Bucket) GetRandom() ([]byte, error) {
	return b.db.GetRandom(b.bType)
}
