package db

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func init() {
	user := *flag.String("user", "root", "db user")
	password := *flag.String("password", "mysql8", "db password")
	host := *flag.String("host", "10.211.55.18", "db host")
	port := *flag.String("port", "3306", "db port")
	database := *flag.String("db", "test", "database")
	dns := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4"
	db, err := sqlx.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
