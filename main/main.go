package main

import (
    "github.com/joho/godotenv"
    "log"
)

func init() {
    if err := godotenv.Load("../.env"); err != nil {
        log.Fatalln("No .env file found")
    }
}

func main() {

}
