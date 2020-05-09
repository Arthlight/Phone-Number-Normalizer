package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"strconv"
	"strings"
)

const (
	host = "localhost"
	port = 5432
	user = "arthred"
	dbname = "phone_number_serializer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable", host, port, user)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	err = resetDB(db, dbname)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	err = createPhoneNumbersTable(db)
	if err != nil {
		fmt.Println(err)
	}

	id, err := insertPhoneNumber(db, "1234567890")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("ID=%d", id)
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

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)`
	_, err := db.Exec(statement)
	return err
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
func normalize(phoneNumber string) string {
	var sb strings.Builder

	for _, char := range phoneNumber {
		if _, err := strconv.ParseInt(string(char), 10, 64); err != nil {
			continue
		} else {
			sb.WriteString(string(char))
		}
	}

	return sb.String()
}