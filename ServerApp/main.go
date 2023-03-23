package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const query = "SELECT plant_name, last_watered, watering_interval_hours, picture FROM plants where plant_id = 123"

type DBstart struct {
	Host     string
	Port     int
	User     string
	Password string
}
type DBConfig struct {
	DBstart DBstart
	DBName  string
}

type Plants struct {
	Plant_name              string
	Last_watered            int64
	Watering_interval_hours int64
}

func main() {
	//check for a DB connection, wait and retry.
	// if there is a DB conn create (if it does not exist) a DB schema and close DB connection
	err := createDB()
	if err != nil {
		fmt.Println(err)
	}
	// return a OK or a DB error to frontend?

	http.HandleFunc("/addplant", addplant)
	http.HandleFunc("/getplants", getplants)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func createDB() error {
	dbconf := DBstart{"userver02.lan", 3306, "root", "VMware1!"}
	dbdetails := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port)
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
	return nil
}

func NewDBconn(config *DBConfig) (*sql.DB, error) {
	dbdetails := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DBstart.User, config.DBstart.Password, config.DBstart.Host, config.DBstart.Port, config.DBName)

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

//func DbQuery(db *sql.DB, query string) (*Plants, error) {
//	var p Plants
//
//	err := row.Scan(&p.Plant_name, &p.Last_watered, &p.Watering_interval_hours)
//	if err != nil {
//		return nil, fmt.Errorf("failed to query DB with error: %w", err)
//	}
//	return &p, nil
//}

func addplant(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the input data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal the JSON input into the Plants struct
	var p Plants
	err = json.Unmarshal(body, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//printing request body, uncomment the next 2 lines for debug
	//fmt.Println("Request body:", string(body))
	//fmt.Println(p.Plant_name, p.Watering_interval_hours, p.Last_watered)

	// set the DB parameters:
	dbconn := DBstart{"userver02.lan", 3306, "root", "VMware1!"}
	dbconf := DBConfig{dbconn, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prepare the SQL statement to insert a new plant
	stmt, err := db.Prepare("INSERT INTO plants (plant_name, last_watered, watering_interval_hours, picture) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement to insert a new plant
	res, err := stmt.Exec(p.Plant_name, p.Last_watered, p.Watering_interval_hours, p.Picture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted plant
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID of the newly inserted plant as the response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}
func getplants(w http.ResponseWriter, r *http.Request) {
	// set the DB parameters:
	dbconn := DBstart{"userver02.lan", 3306, "root", "VMware1!"}
	dbconf := DBConfig{dbconn, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Prepare the SQL statement to retrieve all plants
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement to retrieve all plants
	rows, err := stmt.Query()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Loop through the rows and add each plant to a slice of Plants
	var plants []Plants
	for rows.Next() {
		var p Plants
		err = rows.Scan(&p.Plant_name, &p.Last_watered, &p.Watering_interval_hours)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		plants = append(plants, p)
	}

	// Marshal the slice of Plants to JSON and return it as the response
	json.NewEncoder(w).Encode(plants)
}
