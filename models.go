package ipmonitor

import (
	"log"

	"github.com/jinzhu/gorm"
	// for gorm to use sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// HostModel represents host table in DB
type HostModel struct {
	gorm.Model
	Address     string
	Hostname    string
	Description string
}

// InitDB initializes DB
func InitDB() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Println("[ERROR] /hosts: failed to open DB:", err)
		return
	}
	defer db.Close()

	db.AutoMigrate(&HostModel{})
}
