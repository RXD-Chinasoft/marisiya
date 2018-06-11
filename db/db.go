package db

import (
	"database/sql"
	"log"
	"fmt"
	"github.com/lib/pq"
	"time"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "root"
    dbname   = "postgres"
)
var dbHandler *sql.DB
var psqlInfo string

func init() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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
	
	_, err = dbHandler.Exec("CREATE TABLE IF NOT EXISTS friends (id serial NOT NULL, email character varying(100), friends integer[], subscribMgr text[] ) WITH(OIDS=FALSE);")
	if err != nil {
		log.Printf("create table err %s", err)
	}
}

func trigger(statment string) {
	// dbHandler.Exec("select pg_notify('hello','world')")
	dbHandler.Exec(statment)
}

func listen(subject string) {
	go func(){
		report := func(et pq.ListenerEventType, err error) {
			if err != nil {
				fmt.Println(err)
			}
		}
	
		listener := pq.NewListener(psqlInfo, 10 * time.Second, time.Minute, report)
		err := listener.Listen(subject)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Println("-------start listen------------")
		for {
			n := <-listener.Notify
			log.Println("get notify : ", n.Extra)
			// switch n.Channel {
			// case "hello":
			// 	log.Println("get notify : ", n.Extra)
			// }
		}
	}()
}
