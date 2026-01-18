package database

import (
	"fmt"
	"log"
	"sun-stockanalysis-api/internal/configurations"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	*gorm.DB
}

var (
	postgreDatebaseInstance 	*postgresDatabase
	once						sync.Once				
)

func (p *postgresDatabase) ConnectionGetting() *gorm.DB {
	return p.DB
}

func NewPostgresDatabase(conf *configurations.Database) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
			conf.Host, conf.Port, conf.User, conf.Password, conf.DBname, conf.SSLmode, conf.Schema,
		)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		log.Printf("Connected to database %s", conf.DBname)
		postgreDatebaseInstance = &postgresDatabase{conn}
	})

	return postgreDatebaseInstance
}
