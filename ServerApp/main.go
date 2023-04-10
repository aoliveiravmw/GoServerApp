package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type ConnDB struct {
	Host     string
	Port     int
	User     string
	Password string
}
type DBConfig struct {
	DBstart ConnDB
	DBName  string
}

type Plants struct {
	Plant_name              string
	Plant_alias             string
	Last_watered            int64
	Watering_interval_hours int64
}

type Users struct {
	Name     string
	Surname  string
	Alias    string
	Password string
}

type Pictures struct {
	pic []byte
}

var dbenv ConnDB

func main() {
	//check for a DB connection, wait and retry.
	// if there is a DB conn create (if it does not exist) a DB schema and close DB connection
	envDB(&dbenv)
	err := CreateDB(dbenv)
	if err != nil {
		fmt.Println(err)
	}
	// return a OK or a DB error to frontend?

	http.HandleFunc("/addplant", AddplantHttp)
	//http.HandleFunc("/getplants", GetplantHttp)
	http.HandleFunc("/createuser", AddusHttp)
	http.HandleFunc("/upload", AddPicHttp)
	http.HandleFunc("/download", GetPicHttp)
	//fmt.Println("Server started and listening on port 8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))
	// Serving on HTTPS with TLS
	fmt.Println("Server started and listening on port 4443")
	log.Fatal(http.ListenAndServeTLS(":4443", "server.crt", "server.key", nil))
}
func envDB(v *ConnDB) error {
	v.Host = os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	v.Port, _ = strconv.Atoi(dbPortStr)
	v.User = os.Getenv("DB_USER")
	v.Password = os.Getenv("DB_PASS")
	return nil
}

func CreateDB(dbconf ConnDB) error {
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
	dbconf := DBConfig{dbenv, "plant_watering"}
	db, err := NewDBconn(&dbconf)
	if err != nil {
		return fmt.Errorf("create user failed with error: %w", err)
	}
	defer db.Close()
	//verify if the user alias exists, it has to be unique
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
	insert, _ := db.Prepare("insert into users (name, surname, alias, password) values (?, ?, ?, ?)")
	defer insert.Close()
	_, err = insert.Exec(u.Name, u.Surname, u.Alias, u.Password)
	if err != nil {
		return fmt.Errorf("insert user failed with error: %w", err)
	}
	return nil
}

func AddplantSql(p Plants, u Users) error {
	// set the DB parameters:
	dbconf := DBConfig{dbenv, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		return fmt.Errorf("add plant failed with error: %w", err)
	}
	defer db.Close()
	//insert query
	insert, _ := db.Prepare("insert into plants (plant_name, plant_alias, last_watered, watering_interval_hours, user_alias) values (?, ?, ?, ?, ?)")
	defer insert.Close()
	_, err = insert.Exec(p.Plant_name, p.Plant_alias, p.Last_watered, p.Watering_interval_hours, u.Alias)
	if err != nil {
		return fmt.Errorf("add plant failed with error: %w", err)
	}
	return nil
}

func AddpicSql(pa Plants, p Pictures) error {
	// set the DB parameters:
	dbconf := DBConfig{dbenv, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		return fmt.Errorf("add picture failed with error: %w", err)
	}
	defer db.Close()
	//insert query
	insert, _ := db.Prepare("insert into pictures (picture, plant_alias) values (?, ?)")
	defer insert.Close()
	_, err = insert.Exec(p.pic, pa.Plant_alias)
	if err != nil {
		return fmt.Errorf("add picture failed with error: %w", err)
	}
	return nil
}

func GetPlantssql(u Users) ([]Plants, error) {
	// set the DB parameters:
	dbconf := DBConfig{dbenv, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		return nil, fmt.Errorf("get plants failed with error: %w", err)
	}
	defer db.Close()
	//select query
	stmt, _ := db.Prepare("select plant_name, last_watered, watering_interval_hours, plant_alias from plants where alias_id=?")
	defer stmt.Close()
	rows, err := stmt.Query(u.Alias)
	if err != nil {
		return nil, fmt.Errorf("get plants failed with error: %w", err)
	}
	defer rows.Close()
	var plants []Plants
	for rows.Next() {
		var p Plants
		err = rows.Scan(&p.Plant_name, &p.Last_watered, &p.Watering_interval_hours, &p.Plant_alias)
		if err != nil {
			return nil, fmt.Errorf("get plants failed with error: %w", err)
		}
		plants = append(plants, p)
	}
	return plants, nil
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
	//fmt.Println(p.Plant_name, p.Plant_alias, p.Watering_interval_hours, p.Last_watered)

	//test user
	u := Users{"b", "b", "MC", "1234dsh"}
	err = AddplantSql(p, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("The plant %v was added for user %v", p.Plant_name, u.Alias)
	json.NewEncoder(w).Encode(response)
}

/*
	func GetplantHttp(w http.ResponseWriter, r *http.Request) {
		u := Users{"B", "B", "JC", "gg1234"}
		plants, err := GetPlantssql(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Marshal the slice of Plants to JSON and return it as the response
		json.NewEncoder(w).Encode(plants)
	}
*/
func AddusHttp(w http.ResponseWriter, r *http.Request) {
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

func AddPicHttp(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form file
	err := r.ParseMultipartForm(10 << 20) // 10 MB max form size
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get the picture data from the form
	file, handler, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "Failed to read file from form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded file %+v\n", handler.Filename)
	fmt.Printf("File size %+v\n", handler.Size)
	fmt.Printf("MIME header%+v\n", handler.Header)
	// Read the file data into a byte slice
	picData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read picture data", http.StatusBadRequest)
		return
	}

	// Create a Pictures struct from the file data
	pic := Pictures{pic: picData}

	pa := Plants{"garlic", "Plantero", 8562723, 2984}
	// Store the picture in the database
	err = AddpicSql(pa, pic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func GetPicHttp(w http.ResponseWriter, r *http.Request) {

	// Get the plant alias from the query parameters
	params := r.URL.Query()
	plantAlias := params.Get("plant_alias")

	// set the DB parameters:
	dbconf := DBConfig{dbenv, "plant_watering"}
	// Create a new database connection
	db, err := NewDBconn(&dbconf)
	if err != nil {
		http.Error(w, fmt.Sprintf("get picture failed with error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database for the picture associated with the plant alias
	var pic Pictures
	row := db.QueryRow("select picture from pictures where plant_alias = ?", plantAlias)
	err = row.Scan(&pic.pic)
	if err != nil {
		http.Error(w, fmt.Sprintf("get picture failed with error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Set the content type header to the appropriate image type
	w.Header().Set("Content-Type", "image/png")

	// Serve the image file to the client
	_, err = w.Write(pic.pic)
	if err != nil {
		http.Error(w, fmt.Sprintf("get picture failed with error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
