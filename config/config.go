package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// 環境変数
var (
	SessionKey = makeSessionKey()
	SQLEnv     = makeSQLEnv()
)

func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func makeSessionKey() string {
	envLoad()
	return os.Getenv("SESSION_KEY")
}

func makeSQLEnv() string {
	envLoad()
	return os.Getenv("SQL_ENV")
}
