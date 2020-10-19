package Core

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// CORE ядро нашего прилоежния
type CORE struct {
	DB *sql.DB
}

var instance *CORE
var once sync.Once

// GetInstance Получение экземпляра ядра в виде синглтона
func GetInstance() *CORE {
	once.Do(func() {
		instance = &CORE{}
	})
	return instance
}

// DBInit Инициализация соединения с БД
func (core *CORE) DBInit() {
	log.Info("DB Connection...")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRESQL_HOSTNAME"), os.Getenv("POSTGRESQL_PORT_NUMBER"))
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Panic(err)
	}
	core.DB = db
	_, err = db.Exec("select 1")
	log.Info("Done")
}

// DBClose закрытие соединения с БД
func (core *CORE) DBClose() {
	core.DB.Close()
}

func HandelError(err error, usePanic bool) {
	if err != nil {
		if usePanic {
			log.Panic(err)
		} else {
			log.Error(err)
		}
	}
}
