package main

import (
	"fmt"
	"os"
	"strconv"
)

type DBstart struct {
	Host     string
	Port     int
	User     string
	Password string
}

var bdstart DBstart

func main() {
	envDB(&bdstart)
	println(bdstart.Host)
	fmt.Println("now printing the function")
	printX()
}

func envDB(v *DBstart) error {
	v.Host = os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	v.Port, _ = strconv.Atoi(dbPortStr)
	v.User = os.Getenv("DB_USER")
	v.Password = os.Getenv("DB_PASS")
	return nil
}

func printDB() {
	println(bdstart.Port)
}

func printX() {
	printDB()
}
