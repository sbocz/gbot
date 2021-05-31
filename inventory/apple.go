package inventory

import (
	"fmt"
	"math/rand"
)

type apple struct {
	flavor string
}

func NewApple(flavorIn string) Item {
	return &apple{
		flavor: flavorIn,
	}
}

func (a apple) Identifier() string {
	return a.flavor + "apple"
}

func (a apple) Value() int {
	return 1
}

func (a apple) Description() string {
	return fmt.Sprintf("A %s apple. Mmm...", a.flavor)
}

func (a apple) Rarity() RarityLevel {
	return Common
}

func (a apple) Use() string {
	freshness := rand.Float32()
	if freshness < 0.1 {
		return fmt.Sprintf("...absolutely disgusting. You spit it out.")
	} else if freshness < 0.4 {
		return fmt.Sprintf("A bit too %s. Not the worst, but not great.", a.flavor)
	} else if freshness < 0.7 {
		return fmt.Sprintf("Ahhh.. Refreshing.")
	} else if freshness < 0.9 {
		return fmt.Sprintf("Exactly what you hope for in a %s apple. That hit the spot", a.flavor)
	} else {
		return `Mmmm, plplplplpl. Cover all 9000 taste buds, plplplplpl. Drive up that top note 
		feel, plplplplpl. Yep, that's a 10.`
	}
}
