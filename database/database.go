package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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

	SSLMode       *string
	SSLRootCert   *string
	SSLCert       *string
	SSLKey        *string
	SSLServerName *string

	TimeZone *string
	Params   map[string]*string

	MaxOpenConns    *int
	MaxIdleConns    *int
	ConnMaxLifetime *int
}

func (c Configuration) DSN() string {
	sslMode := "require"
	if c.SSLMode != nil {
		sslMode = *c.SSLMode
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Password,
		sslMode,
	)

	if c.SSLRootCert != nil {
		dsn += fmt.Sprintf(" sslrootcert=%s", *c.SSLRootCert)
	}
	if c.SSLCert != nil {
		dsn += fmt.Sprintf(" sslcert=%s", *c.SSLCert)
	}
	if c.SSLKey != nil {
		dsn += fmt.Sprintf(" sslkey=%s", *c.SSLKey)
	}
	if c.SSLServerName != nil {
		dsn += fmt.Sprintf(" sslservername=%s", *c.SSLServerName)
	}
	if c.TimeZone != nil {
		dsn += fmt.Sprintf(" TimeZone=%s", *c.TimeZone)
	}

	for k, v := range c.Params {
		if v != nil {
			dsn += fmt.Sprintf(" %s=%s", k, *v)
		}
	}

	return dsn
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
		dbConString := config.DSN()
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
