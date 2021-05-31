package inventory

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
)

type UserInventory struct {
	UserId  discord.UserID
	ItemMap map[string]int
}

func NewUserInventory(id discord.UserID) *UserInventory {
	return &UserInventory{
		UserId:  discord.UserID(id),
		ItemMap: make(map[string]int),
	}
}

func (u *UserInventory) AddItem(itemId string) {
	u.ItemMap[itemId] += 1
}

func (u *UserInventory) RemoveItem(itemId string) error {
	if u.ItemMap[itemId] < 1 {
		return fmt.Errorf("Cannot remove item %s from user %v. They have %v in their inventory.",
			itemId, u.UserId, u.ItemMap[itemId])
	}
	u.ItemMap[itemId] -= 1
	return nil
}
