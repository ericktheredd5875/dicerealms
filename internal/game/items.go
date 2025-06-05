package game

import (
	"log"

	"gorm.io/gorm"

	"github.com/ericktheredd5875/dicerealms/internal/db"
)

var itemRegistry = make(map[string]*db.ItemModel)

func LoadAllItems(gdb *gorm.DB) {
	items := []db.ItemModel{}
	if err := gdb.Find(&items).Error; err != nil {
		log.Printf("Failed to load items: %v", err)
		return
	}

	for _, i := range items {
		item := i
		itemRegistry[item.Name] = &item
	}

	log.Printf("Loaded %d items", len(itemRegistry))
}

func GetItemByName(name string) *db.ItemModel {
	item, ok := itemRegistry[name]
	if !ok {
		item = &db.ItemModel{}
		if err := db.DB.Where("name = ?", name).First(&item).Error; err != nil {
			log.Printf("Item not found: %v", err)
			return nil
		}
		itemRegistry[name] = item
	}

	return item
}

func GetRandomItem() *db.ItemModel {
	item := &db.ItemModel{}
	if err := db.DB.Order("RANDOM()").First(&item).Error; err != nil {
		log.Printf("Failed to get random item: %v", err)
		return nil
	}
	itemRegistry[item.Name] = item
	return item
}

func GetRandomItems(n int) []*db.ItemModel {
	var items []db.ItemModel
	if err := db.DB.Order("RANDOM()").Limit(n).Find(&items).Error; err != nil {
		log.Printf("Failed to get random items: %v", err)
		return nil
	}

	var result []*db.ItemModel
	for _, i := range items {
		item := i
		itemRegistry[item.Name] = &item
		result = append(result, &item)
	}

	return result
}

func GetAllItemsInRoom(roomID uint) []*db.ItemModel {
	var items []db.ItemModel
	if err := db.DB.Where("room_found_id = ?", roomID).Find(&items).Error; err != nil {
		log.Printf("Failed to get items in room %d: %v", roomID, err)
		return nil
	}

	var result []*db.ItemModel
	for _, i := range items {
		item := i
		itemRegistry[item.Name] = &item
		result = append(result, &item)
	}

	return result
}
