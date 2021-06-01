package shop

import (
	"bytes"
	"fmt"
	"gbot/commerce"
	"gbot/commerce/banking"
	"gbot/inventory"
	"text/tabwriter"

	"github.com/diamondburned/arikawa/v2/discord"
)

type InventoryEntry struct {
	Count int
	Item  inventory.Item
}

type Shop struct {
	Entries map[string]*InventoryEntry

	bank            *banking.Bank
	globalInventory *inventory.Inventory
}

func NewShop(shopBank *banking.Bank, globalInventory *inventory.Inventory) *Shop {
	allItems := globalInventory.ItemSet()

	initialEntries := make(map[string]*InventoryEntry)
	for itemId, item := range allItems {
		var count int
		if item.Rarity() == inventory.Legendary {
			count = 1
		} else if item.Rarity() == inventory.Epic {
			count = 2
		} else if item.Rarity() == inventory.Rare {
			count = 3
		} else if item.Rarity() == inventory.Uncommon {
			count = 5
		} else if item.Rarity() == inventory.Common {
			count = 10
		}

		initialEntries[itemId] = &InventoryEntry{
			Count: count,
			Item:  globalInventory.FetchItem(itemId),
		}
	}

	return &Shop{
		Entries:         initialEntries,
		bank:            shopBank,
		globalInventory: globalInventory,
	}
}

func (s *Shop) BuyItem(userId discord.UserID, itemId string, count int) error {
	if shopEntry, ok := s.Entries[itemId]; ok {
		if count > shopEntry.Count {
			return fmt.Errorf("Cannot buy %v of item %s. Only %v available.", count, itemId, shopEntry.Count)
		}

		price := count * shopEntry.Item.Value()
		err := s.bank.Withdraw(userId, price)
		if err != nil {
			return fmt.Errorf("Could not withdraw funds for purchase: %s", err)
		}

		userInventory, err := s.globalInventory.FetchUserInventory(userId)
		if err != nil {
			return fmt.Errorf("Could fetch inventory for purchase: %s", err)
		}
		userInventory.ItemMap[shopEntry.Item.Identifier()] += count
		err = s.globalInventory.WriteUserInventory(userInventory)
		if err != nil {
			return fmt.Errorf("Could write inventory for purchase: %s", err)
		}
		shopEntry.Count -= count

		return nil
	}
	return fmt.Errorf("Item %s does not exist in the shop.", itemId)
}

func (s *Shop) PrintInventory() string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "Name\tRarity\tDescription\tPrice\tQuantity\n")

	for itemId, entry := range s.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n",
			itemId,
			entry.Item.Rarity(),
			entry.Item.Description(),
			commerce.Currency(entry.Item.Value()),
			entry.Count)
	}
	w.Flush()
	return b.String()
}
