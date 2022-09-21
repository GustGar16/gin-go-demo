package configs

import (
	"log"

	"github.com/joho/godotenv"
)

var myEnv map[string]string

func readEnv() {
	read, err := godotenv.Read()
	myEnv = read
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
}

func EnvMongoURI() string {
	readEnv()
	return myEnv["MONGOURI"]
}

func CurrentDatabase() string {
	readEnv()
	return myEnv["MONGODB"]
}
