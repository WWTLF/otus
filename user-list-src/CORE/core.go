package Core

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

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
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"))
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Errorf("sql.Open(\"postgres\", dbinfo) ", err)
		panic(err)
	}
	core.DB = db
	_, err = db.Exec("select 1")
	if err != nil {
		panic("dbinfo = " + dbinfo)
	}
	log.Info("Done")
}

// DBClose закрытие соединения с БД
func (core *CORE) DBClose() {
	core.DB.Close()
}

func HandelError(err error, usePanic bool) {
	if err != nil {
		log.Error(err)
		if usePanic {
			log.Fatal(err)
		}
	}
}
