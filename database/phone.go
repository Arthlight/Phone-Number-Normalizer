package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	db *sql.DB
}

// Phone represents the phone_numbers table in the DB
type Phone struct{
	ID 		int
	Number 	string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}


func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}

	return db.Close()
}
func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return createDB(db, name)
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		fmt.Println(err)
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}

	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)`
	_, err := db.Exec(statement)
	return err
}

func Seed(d *DB) error {
	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range phoneNumbers {
		_, err := insertPhoneNumber(d.db, number)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertPhoneNumber(db *sql.DB, phoneNumber string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phoneNumber).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (d *DB) FindPhoneNumber(number string) (*Phone, error) {
	var p Phone
	err := d.db.QueryRow("SELECT * FROM phone_numbers WHERE value=$1", number).Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (d *DB) UpdateRow(p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := d.db.Exec(statement, p.ID, p.Number)

	return err
}

func (d *DB) DeleteRow(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := d.db.Exec(statement, id)

	return err
}

func (d *DB) AllPhoneNumbers() ([]Phone, error){
	rows, err := d.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}