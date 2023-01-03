package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

// DB global var
var DB *gorm.DB

type Configuration struct {
	Driver     string
	Name       string
	User       string
	Password   string
	Host       string
	Port       string
	LogLevel   string
	DBLogLevel logger.LogLevel
}

func NewConfiguration(config *Configuration) *Configuration {
	logLevelMap := map[string]logger.LogLevel{
		"Silent": logger.Silent,
		"Error":  logger.Error,
		"Warn":   logger.Warn,
		"Info":   logger.Info,
	}
	config.DBLogLevel = logLevelMap[config.LogLevel]
	return config
}

type TableNameReplacer struct{}

func (r TableNameReplacer) Replace(name string) string {
	if len(name) < 4 {
		return name
	}
	const postfix = "Tbl"
	const postfixCharacters = len(postfix)
	lastCharacters := name[len(name)-postfixCharacters:]
	if lastCharacters == postfix {
		return name[0 : len(name)-postfixCharacters]
	}
	return name
}

// InitDB opens a database and saves the reference to `Database` struct.
func InitDB(config *Configuration) *gorm.DB {
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
		DB, err = gorm.Open(postgres.Open(dbConString), &gorm.Config{
			Logger: logger.Default.LogMode(config.DBLogLevel),
			NamingStrategy: schema.NamingStrategy{
				NameReplacer: TableNameReplacer{}, // use name replacer to change struct/field name before convert it to db name
			},
		})

		if err != nil {
			log.Printf("Error when connecting to postgres db, %s", err)
		}
	case "mysql":
		dbConString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Locals",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Name,
		)
		DB, err = gorm.Open(mysql.Open(dbConString), &gorm.Config{
			Logger: logger.Default.LogMode(config.DBLogLevel),
			NamingStrategy: schema.NamingStrategy{
				NameReplacer: TableNameReplacer{}, // use name replacer to change struct/field name before convert it to db name
			},
		})

		if err != nil {
			log.Printf("Error when connecting to postgres db, %s", err)
		}
	default:
		panic("Invalid DB Driver")
	}

	return DB
}

// GetDB helps you to get a connection
func GetDB() *gorm.DB {
	return DB
}
