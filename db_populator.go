package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"math/rand"
	_ "github.com/go-sql-driver/mysql"
)

// Connect to Marian DB database
func Connect() *sql.DB {

	// Catch Panics
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic caught in Connect : %s \n", err)
		}
	}()

	// Needs to be dynamically defined based on Cli Arguments in the future
	schema := "filius_stultus:*nkjrko2cp9!L0NV&qoqSy0rZCXTShc&nsDq#TSOz@tcp(bollywood-test.cdcgawleh4ss.us-east-1.rds.amazonaws.com:3306)/bollywood"
	log.Println(schema)

	// Return DB Object
	db, err := sql.Open("mysql", schema)
	if err != nil || db == nil {
		fmt.Println(err)
	}

	// Set Max Connections to the db to 100
	db.SetMaxOpenConns(100)

	// return *sql.DB object
	fmt.Println("Database connection established", db)
	return db

}

// CreateNewUser adds a new user to the DB with Unique Device ID; the user agent not generated.
// Instead, a random user-agent is retrieved from a different device in the DB
// IP Address is randomly generated
func CreateNewUser(logger *log.Logger, db *sql.DB, appID string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Panicf("Panic Caught in CreateNewUser : %s \n", err)
		}
	}()

	deviceID := getDeviceID()
	userAgent := getUserAgent(logger, db)
	ipAddress := getIPAddress()

	rows, err := db.Query("INSERT INTO user ("+
		"user_device_id,"+
		"user_agent,"+
		"user_ip_address,"+
		"user_installed,"+
		"user_friends_invited,"+
		"user_total_purchases,"+
		"user_total_sessions,"+
		"user_items_in_cart,"+
		"user_quality,"+
		"app_id) VALUES (?,?,?,0,0,0,0,0,900,?)", deviceID, userAgent, ipAddress, appID)

	if err != nil {
		logger.Printf("Error in CreateNewUser: ")
		logger.Panicln(err)
	}
	defer rows.Close()
	var response string

	for rows.Next() {
		if err := rows.Scan(
			response,
		); err != nil {
			logger.Println(err)
		}
	}

	logger.Println("User Created :", response)

}

// Generates a UUID for the new User
func getDeviceID() string {

	return randString(8) + "-" + randString(4) + "-" + randString(4) + "-" + randString(4) + "-" + randString(12)

}

// Grabs a random user-agent from the database
func getUserAgent(logger *log.Logger, db *sql.DB) string {
	var userAgent string

	defer func() {
		if err := recover(); err != nil {
			logger.Panicf("Panic Caught in getUserAgent : %s \n", err)
		}
	}()

	rows, err := db.Query("SELECT user.user_agent FROM user WHERE (user.user_agent like '%iPhone%') OR (user.user_agent like '%iPad%') ORDER BY RAND() LIMIT 1")

	if err != nil {
		logger.Println(err)
	}
	if rows == nil {
		logger.Panicln("rows is nil")
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&userAgent,
		); err != nil {
			logger.Println(err)
		}
	}

	return userAgent

}

// Generate Random IP Address
func getIPAddress() string {
	return strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255))

}

// Defines Random AlphaNumeric String of a given Length
func randString(n int) string {

	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
