package main

import (
	"strconv"
	"strings"
	"database/sql"
	"github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "arthred"
	password = ""
	dbname = "gophercises_phone"
)

func main() {

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