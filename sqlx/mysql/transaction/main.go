package main

import (
	"sqlx/mysql/db"
)

func main() {

}

func Transaction() (err error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	_, err = tx.Exec("DELETE user WHERE id = 1")
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE user WHERE id = 2")
	if err != nil {
		return err
	}
	return
}
