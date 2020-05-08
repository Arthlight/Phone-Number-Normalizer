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
	dbname = "Phone-Number-Serializer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable", host, port, user)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	err = createDB(db, dbname)
	if err != nil {
		fmt.Println(err)
	}

	db.Close()
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE" + name)
	if err != nil {
		return err
	}
	return nil
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