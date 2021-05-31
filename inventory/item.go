package inventory

type Item interface {
	Identifier() string
	Value() int
	Description() string
	Rarity() RarityLevel
	Use() string
}

type RarityLevel int

const (
	Common RarityLevel = iota
	Uncommon
	Rare
	Epic
	Legendary
)

func (r RarityLevel) String() string {
	return [...]string{"Common", "Uncommon", "Rare", "Epic", "Legendary"}[r]
}
