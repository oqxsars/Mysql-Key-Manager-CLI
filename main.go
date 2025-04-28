package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Key       string
	Note      string
	UUID      string
	TimeMade  string
	IPDomain  string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:PASSWORD@tcp(127.0.0.1:3306)/DATABASENAME")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Key Manager")
	fmt.Println("---------------------")
	fmt.Println("Available commands:")
	fmt.Println("addkey <key> <note> <ip/domain>")
	fmt.Println("removekey <key>")
	fmt.Println("exit")
	fmt.Println("---------------------")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		command := args[0]
		switch command {
		case "addkey":
			if len(args) != 4 {
				fmt.Println("Usage: addkey <key> <note> <ip/domain>")
				continue
			}
			addKey(args[1], args[2], args[3])
		case "removekey":
			if len(args) != 2 {
				fmt.Println("Usage: removekey <key>")
				continue
			}
			removeKey(args[1])
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Available commands: addkey, removekey, exit")
		}
	}
}

func addKey(key, note, ipDomain string) {
	newUUID := uuid.New().String()
	timeMade := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO users (`key`, note, uuid, timeofcreation, domain_or_ip) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, key, note, newUUID, timeMade, ipDomain)
	if err != nil {
		log.Println("Error inserting record:", err)
		return
	}

	fmt.Printf("Key added successfully: UUID = %s\n", newUUID)
}

func removeKey(key string) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE `key` = ?", key).Scan(&count)
	if err != nil {
		log.Println("Error checking key:", err)
		return
	}

	if count == 0 {
		fmt.Printf("Key '%s' not found in database.\n", key)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE `key` = ?", key)
	if err != nil {
		log.Println("Error removing record:", err)
		return
	}
	fmt.Printf("Key '%s' removed successfully.\n", key)
}
