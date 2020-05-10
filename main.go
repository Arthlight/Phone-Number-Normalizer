package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

const (
	host = "localhost"
	port = 5432
	user = "arthred"
	dbname = "phone_number_serializer"
)

var (
	phoneNumbers = []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
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

	for _, number := range phoneNumbers {
		_, err = insertPhoneNumber(db, number)
		if err != nil {
			fmt.Println(err)
		}
	}

	phones, _ := allPhoneNumbers(db)
	for _, n := range phones {
		fmt.Printf("Working on...%+v\n", n)
		number := normalize(n.number)
		if number != n.number {
			fmt.Println("Updating or removing...", number)
		} else {
			fmt.Println("No changes required")
		}
	}
}

func getPhoneNumber(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

type phone struct{
	id 		int
	number 	string
}

func allPhoneNumbers(db *sql.DB) ([]phone, error){
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
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