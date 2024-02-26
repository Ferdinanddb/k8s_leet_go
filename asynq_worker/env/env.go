package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func VerifyEnv() {

	_, isSetDB_HOST := os.LookupEnv("DB_HOST")
	_, isSetDB_USER := os.LookupEnv("DB_USER")
	_, isSetPGPASSWORD := os.LookupEnv("PGPASSWORD")
	_, isSetDB_NAME := os.LookupEnv("DB_NAME")
	_, isSetDB_PORT := os.LookupEnv("DB_PORT")

	_, isSetREDIS_HOST := os.LookupEnv("REDIS_HOST")
	_, isSetREDIS_PORT := os.LookupEnv("REDIS_PORT")
	_, isSetREDIS_PASSWORD := os.LookupEnv("REDIS_PASSWORD")

	if 	!isSetDB_HOST || !isSetDB_USER || !isSetPGPASSWORD || !isSetDB_NAME || !isSetDB_PORT || 
		!isSetREDIS_HOST || !isSetREDIS_PORT || !isSetREDIS_PASSWORD {
		fmt.Printf(
			"Env variables that are set :\nDB_HOST -> %v\nDB_USER -> %v\nPGPASSWORD-> %v\nDB_NAME -> %v\nDB_PORT -> %v\n",
			isSetDB_HOST,
			isSetDB_USER,
			isSetPGPASSWORD,
			isSetDB_NAME,
			isSetDB_PORT,
		)

		fmt.Printf(
			"Env variables that are set :\nREDIS_HOST -> %v\nREDIS_PORT -> %v\nREDIS_PASSWORD-> %v\n\n",
			isSetREDIS_HOST,
			isSetREDIS_PORT,
			isSetREDIS_PASSWORD,
		)

		fmt.Println("Attempting to log an .env file located at the root of the project...")

		err := godotenv.Load(".do_not_push/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		fmt.Println("All the env variables needed are set, processing further...")
	}

}
