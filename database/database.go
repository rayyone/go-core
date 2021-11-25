package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/inflection"

	"github.com/jinzhu/gorm"

	// blank import but used for driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB global var
var DB *gorm.DB

type Configuration struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
}

// InitDB opens a database and saves the reference to `Database` struct.
func InitDB(config Configuration, debug bool) *gorm.DB {
	var err error

	switch config.Driver {
	case "postgres":
		dbConString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			config.Host,
			config.Port,
			config.User,
			config.Name,
			config.Password,
		)
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			if len(defaultTableName) < 6 {
				return defaultTableName
			}
			last5Chars := defaultTableName[len(defaultTableName)-5:]
			if last5Chars == "_tbls" {
				return inflection.Plural(defaultTableName[0 : len(defaultTableName)-5])
			}
			return defaultTableName
		}
		DB, err = gorm.Open("postgres", dbConString)

		if err != nil {
			log.Printf("Error when connecting to postgres db, %s", err)
		}
	case "mysql":
		dbConString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Name,
		)
		DB, err = gorm.Open("mysql", dbConString)
		if err != nil {
			log.Printf("Error when connecting to mysql db, %s", err)
		}
	default:
		panic("Invalid DB Driver")
	}

	DB.LogMode(debug)

	return DB
}

// GetDB helps you to get a connection
func GetDB() *gorm.DB {
	return DB
}
