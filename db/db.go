package db

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "root"
    dbname   = "postgres"
)
var dbHandler *sql.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("open db err %s", err)
	}
	err = db.Ping()
    if err != nil {
        log.Printf("ping db err %s", err)
    } else {
		fmt.Println("Successfully connected!")
	}
	dbHandler = db
	
	_, err = dbHandler.Exec("CREATE TABLE IF NOT EXISTS friends (id integer, email character varying(100), friends integer[] ) WITH(OIDS=FALSE);")
	if err != nil {
		log.Printf("create table err %s", err)
	}
}
