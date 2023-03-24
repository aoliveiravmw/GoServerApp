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
	AliasID                 string
}

type Users struct {
	Name    string
	Surname string
	Alias   string
}

type Pictures struct {
	pic byte
}

func main() {
	//check for a DB connection, wait and retry.
	// if there is a DB conn create (if it does not exist) a DB schema and close DB connection
	err := CreateDB()
	if err != nil {
		fmt.Println(err)
	}
	// return a OK or a DB error to frontend?

	http.HandleFunc("/addplant", AddplantHttp)
	//http.HandleFunc("/getplants", getplants)
	http.HandleFunc("/createuser", Addushttp)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func CreateDB() error {
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

func Createusql(u Users) error {
	dbconn := DBstart{"userver02.lan", 3306, "root", "VMware1!"}
	dbconf := DBConfig{dbconn, "plant_watering"}
	db, err := NewDBconn(&dbconf)
	if err != nil {
		return fmt.Errorf("create user failed with error: %w", err)
	}
	defer db.Close()
	//verify if the user alias exists
	query, _ := db.Prepare("select alias from users where alias=?")
	defer query.Close()
	resQuery := query.QueryRow(u.Alias)
	var alias string
	err = resQuery.Scan(&alias)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("query failed with error: %w", err)
	}
	if alias != "" {
		return fmt.Errorf("alias already exists")
	}
	//insert the user if alias is new
	insert, _ := db.Prepare("insert into users (name, surname, alias) values (?, ?, ?)")
	defer insert.Close()
	_, err = insert.Exec(u.Name, u.Surname, u.Alias)
	if err != nil {
		return fmt.Errorf("insert user failed with error: %w", err)
	}
	return nil
}

func AddplantSql(p Plants) error {
	// set the DB parameters:
	dbconn := DBstart{"userver02.lan", 3306, "root", "VMware1!"}
	dbconf := DBConfig{dbconn, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		return fmt.Errorf("add plant failed with error: %w", err)
	}
	defer db.Close()
	//insert query
	insert, _ := db.Prepare("insert into plants (plant_name, last_watered, watering_interval_hours, alias_id) values (?, ?, ?, ?)")
	defer insert.Close()
	_, err = insert.Exec(p.Plant_name, p.Last_watered, p.Watering_interval_hours, p.AliasID)
	if err != nil {
		return fmt.Errorf("add plant failed with error: %w", err)
	}
	return nil
}

func AddplantHttp(w http.ResponseWriter, r *http.Request) {
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
	err = AddplantSql(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("The plant %v was added for user %v", p.Plant_name, p.AliasID)
	json.NewEncoder(w).Encode(response)
}

/*
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
*/
func Addushttp(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the input data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Unmarshal the JSON input into the Plants struct
	var nwuser Users
	err = json.Unmarshal(body, &nwuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = Createusql(nwuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
