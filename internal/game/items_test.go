package game

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ericktheredd5875/dicerealms/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestLoadAllItems(t *testing.T) {
	// Set up the sqlmock database connection
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Inject the mock DB into GORM
	gormDB, err := db.NewGormFromSQLDB(sqlDB)
	assert.NoError(t, err)
	db.DB = gormDB

	// Define the expected query pattern GORM will generate for PostgreSQL
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "item_models" WHERE "item_models"."deleted_at" IS NULL`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rarity", "effect", "category"}).
				AddRow(1, "Sword of Truth", "A legendary sword", "legendary", "2d6 slashing", "weapon").
				AddRow(2, "Potion of Healing", "Heals 2d4+2 HP", "common", "2d4+2 healing", "potion"),
		)

	// Call the function being tested
	LoadAllItems()

	// Retrieve the loaded item and validate
	item := GetItemByName("Sword of Truth")
	assert.NotNil(t, item)
	assert.Equal(t, "Sword of Truth", item.Name)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetItemByName(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := db.NewGormFromSQLDB(sqlDB)
	assert.NoError(t, err)
	db.DB = gormDB

	name := "Amulet of Power"

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "item_models" WHERE name = $1 AND "item_models"."deleted_at" IS NULL ORDER BY "item_models"."id" LIMIT $2`,
	)).
		WithArgs(name, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rarity", "effect", "category"}).
				AddRow(1, name, "A glowing amulet", "rare", "+2 strength", "jewelry"),
		)

	item := GetItemByName(name)
	assert.NotNil(t, item)
	assert.Equal(t, name, item.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRandomItem(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := db.NewGormFromSQLDB(sqlDB)
	assert.NoError(t, err)
	db.DB = gormDB

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "item_models" WHERE "item_models"."deleted_at" IS NULL ORDER BY RANDOM(),"item_models"."id" LIMIT $1`,
	)).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rarity", "effect", "category"}).
				AddRow(1, "Mystery Box", "Who knows what’s inside?", "uncommon", "???", "wondrous"),
		)

	item := GetRandomItem()
	assert.NotNil(t, item)
	assert.Equal(t, "Mystery Box", item.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRandomItems(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := db.NewGormFromSQLDB(sqlDB)
	assert.NoError(t, err)
	db.DB = gormDB

	mock.ExpectQuery(
		`SELECT \* FROM "item_models" WHERE "item_models"\."deleted_at" IS NULL ORDER BY RANDOM\(\) LIMIT \$1`,
	).
		WithArgs(3).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rarity", "effect", "category"}).
				AddRow(1, "Ring of Speed", "Increases your movement", "rare", "+10 feet", "jewelry").
				AddRow(2, "Lantern of Revealing", "Illuminates the invisible", "uncommon", "Detects invisible", "wondrous").
				AddRow(3, "Boots of Jumping", "Jump higher", "common", "Triple jump distance", "armor"),
		)

	items := GetRandomItems(3)
	assert.Len(t, items, 3)
	assert.Equal(t, "Ring of Speed", items[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllItemsInRoom(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := db.NewGormFromSQLDB(sqlDB)
	assert.NoError(t, err)
	db.DB = gormDB

	roomID := uint(42)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "item_models" WHERE room_found_id = $1 AND "item_models"."deleted_at" IS NULL`,
	)).
		WithArgs(roomID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description", "rarity", "effect", "category"}).
				AddRow(1, "Torch", "Lights the way", "mundane", "Provides light", "gear").
				AddRow(2, "Backpack", "Carries your stuff", "common", "+10 inventory", "gear"),
		)

	items := GetAllItemsInRoom(roomID)
	assert.Len(t, items, 2)
	assert.Equal(t, "Torch", items[0].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}
