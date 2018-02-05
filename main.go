package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

//Randomize salt, which is specifically necessary for unique device-id creation
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	db := Connect()
	logger := CreateLogger("log")

	appID := flag.String("app_id", "0", "App id to populate users to")
	flag.Parse()

	for {
		time.Sleep(1 * time.Second)
		CreateNewUser(logger, db, *appID)
	}
}

// CreateLogger creates a logger that writes to the given filename
func CreateLogger(filename string) *log.Logger {
	file, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
