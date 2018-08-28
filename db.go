package ipmonitor

import (
	"github.com/jinzhu/gorm"
	// for gorm to use sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Conn stores opened DB connection
var Conn *DBConnection

// DBConnection stores DB connection
type DBConnection struct {
	DB *gorm.DB
}

// InitDB initializes DB
func InitDB() {
	Conn.DB.AutoMigrate(&Host{})
}

// OpenDB open DB connection and store it to DBConnection struct
func OpenDB(dbname string) error {
	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	Conn = &DBConnection{DB: db}
	return nil
}
