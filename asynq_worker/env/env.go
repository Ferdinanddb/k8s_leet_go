package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_, isSetREDIS_HOST := os.LookupEnv("REDIS_HOST")
	_, isSetREDIS_PORT := os.LookupEnv("REDIS_PORT")
	_, isSetREDIS_PASSWORD := os.LookupEnv("REDIS_PASSWORD")

	if !isSetREDIS_HOST || !isSetREDIS_PORT || !isSetREDIS_PASSWORD {
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
