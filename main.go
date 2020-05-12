package main

import (
	"Phone-Number-Serializer/database"
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

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable", host, port, user)
	err := database.Reset("postgres", psqlInfo, dbname)
	if err != nil {
		fmt.Println(err)
	}
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	err = database.Migrate("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}

	db, err := database.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	err = database.Seed(db)
	if err != nil {
		fmt.Println(err)
	}

	phones, err := db.AllPhoneNumbers()
	if err != nil {
		fmt.Println(err)
	}

	for _, n := range phones {
		fmt.Printf("Working on...%+v\n", n)
		number := normalize(n.Number)
		if number != n.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhoneNumber(number)
			if err != nil {
				fmt.Println(err)
			}
			if existing != nil {
				err := db.DeleteRow(n.ID)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				n.Number = number
				err := db.UpdateRow(&n)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("No changes required")
		}
	}
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