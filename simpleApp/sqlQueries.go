package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type dbenv struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}
type pictures struct {
	pic []byte
}

func envDB(v *dbenv) error {
	v.Host = os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	v.Port, _ = strconv.Atoi(dbPortStr)
	v.User = os.Getenv("DB_USER")
	v.Password = os.Getenv("DB_PASS")
	v.DBname = os.Getenv("DB_NAME")
	return nil
}
func vcapDB(v *dbenv, cfMysql MySQLService) error {
	v.Host = cfMysql.Credentials.Hostname
	v.Port = cfMysql.Credentials.Port
	v.User = cfMysql.Credentials.Username
	v.Password = cfMysql.Credentials.Password
	v.DBname = os.Getenv("DB_NAME")
	return nil
}
func createDB(config dbenv) {
	dbdetails := fmt.Sprintf("%s:%s@tcp(%s:%d)/", config.User, config.Password, config.Host, config.Port)
	db, err := sql.Open("mysql", dbdetails)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlQueryBytes, err := ioutil.ReadFile("createDB.sql")
	if err != nil {
		log.Fatal(err)
	}
	// Remove whitespace from SQL queries
	sqlQueryStr := strings.TrimSpace(string(sqlQueryBytes))
	// Split SQL queries into separate statements
	sqlStatements := strings.Split(sqlQueryStr, ";")
	// Execute each SQL statement
	for _, sqlStatement := range sqlStatements {
		if len(strings.TrimSpace(sqlStatement)) == 0 {
			continue
		}
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
	//For debugging:
	//fmt.Println("Database and tables created successfully!")
}
func connectDB(config dbenv) (*sql.DB, error) {
	dbdetails := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.DBname)
	db, err := sql.Open("mysql", dbdetails)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
func initDB() (*sql.DB, error) {
	var dbenv dbenv
	var cfMysql MySQLService
	err := vcapSqlService(&cfMysql)
	if err != nil {
		fmt.Println(err)
		envDB(&dbenv)
	} else {
		vcapDB(&dbenv, cfMysql)
	}
	createDB(dbenv)
	db, err := connectDB(dbenv)
	if err != nil {
		return nil, fmt.Errorf("db initialization has failed: %w", err)
	}
	return db, nil
}
func getCountrows(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(picture_id) FROM pictures").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (p *pictures) addFileSql(db *sql.DB) error {
	insert, _ := db.Prepare("insert into pictures (picture) values ( ?)")
	defer insert.Close()
	_, err := insert.Exec(p.pic)
	if err != nil {
		return fmt.Errorf("add picture failed with error: %w", err)
	}
	return nil
}
func (p *pictures) getFileSql(db *sql.DB) error {
	return db.QueryRow("SELECT picture FROM pictures").Scan(&p.pic)
}

func deleteRows(db *sql.DB) error {
	_, err := db.Exec("delete from pictures")
	if err != nil {
		return err
	}
	return nil
}
