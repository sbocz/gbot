package inventory

import (
	"encoding/json"
	"fmt"
	"gbot/database"

	"github.com/diamondburned/arikawa/v2/discord"
)

type Inventory struct {
	inventoryBucket *database.Bucket
	itemSet         map[string]Item
}

func NewInventory(db *database.DB) (*Inventory, error) {
	items := []Item{}
	items = append(items, NewApple("sweet"), NewApple("sour"))

	initialItemSet := make(map[string]Item)
	for _, item := range items {
		initialItemSet[item.Identifier()] = item
	}

	bucket, err := database.NewBucket(db, database.InventoryData)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize Inventory: %v", err)
	}
	return &Inventory{
		inventoryBucket: bucket,
		itemSet:         initialItemSet,
	}, nil
}

func (i *Inventory) FetchItem(itemId string) Item {
	return i.itemSet[itemId]
}

func (i *Inventory) ItemSet() map[string]Item {
	targetMap := make(map[string]Item)

	for key, value := range i.itemSet {
		targetMap[key] = value
	}
	return targetMap
}

func (i *Inventory) FetchUserInventory(userId discord.UserID) (*UserInventory, error) {
	rawValue, err := i.inventoryBucket.Get(fmt.Sprint(userId))
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve inventory for %s: %s", userId, err)
	}

	// User does not exist
	if rawValue == nil {
		return NewUserInventory(userId), nil
	}

	var userInventory *UserInventory
	err = json.Unmarshal(rawValue, &userInventory)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshall response %s: %s", rawValue, err)
	}
	return userInventory, nil
}

func (i *Inventory) WriteUserInventory(userInventory *UserInventory) error {
	return i.inventoryBucket.Put(fmt.Sprint(userInventory.UserId), userInventory)
}
