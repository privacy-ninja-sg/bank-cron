package main

import (
	"bank-crons/worker"
	"log"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	log.Println("Cron is running...")
	gocron.Every(30).Seconds().Do(worker.GetBayBalance)
	gocron.Every(35).Seconds().Do(worker.IncreaseBAYBalance)

	gocron.Every(30).Second().Do(worker.GetSCBBalance)
	gocron.Every(35).Second().Do(worker.IncreaseSCBBalance)
	<-gocron.Start()
}
