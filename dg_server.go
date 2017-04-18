package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
)

const dbFilePath string = "datagrammar.db"

var name2key *bolt.Bucket
var key2name *bolt.Bucket
var column *bolt.Bucket

var boltDBinstance *bolt.DB
var database Database

// DatabaseList contains a list of the databases
type DatabaseList struct {
	Databases map[string]string
}

// TableList contains a list of the tables
type TableList struct {
	Tables       map[string]string
	DatabaseName string
	DatabaseKey  string
}

// ColumnList contains a table's worth of columns
type ColumnList struct {
	Columns      map[string]string
	DatabaseName string
	DatabaseKey  string
	TableName    string
	TableKey     string
}

// Database contains a list of schemas
type Database struct {
	Name       string
	Tables     map[string]Table
	Server     string
	TableNames []string
}

// NewDatabase initializes the Database struct with a map
func NewDatabase(name string) *Database {
	db := Database{
		Tables: make(map[string]Table),
		Name:   name,
		Server: "",
	}

	return &db
}

// AddorGetTable adds a table if it's not there, otherwise returns it
func AddorGetTable(name string, database *Database) *Table {
	var table Table
	return &table
}

// Table is the tables in the systems
type Table struct {
	Name         string
	Columns      map[uint64]Column
	Schema       string
	ColumnNames  []string
	DatabaseName string
}

// NewTable initializes the Database struct with a map
func NewTable(name, schema string, dbName string) Table {
	tb := Table{
		Columns:      make(map[uint64]Column),
		Name:         name,
		Schema:       schema,
		DatabaseName: dbName,
	}

	return tb
}

// Column is the base type here
type Column struct {
	ID           uint64
	Name         string
	Ordinal      int
	Type         string
	Length       int
	Precision    int
	Scale        int
	TableName    string
	DatabaseName string
}

// Entry is the single database entry
type Entry struct { // Our example struct, you can use "-" to ignore a field
	Database  string
	System    string
	Schema    string
	Table     string
	Column    string
	Ordinal   int
	Type      string
	Length    int
	Precision int
	Scale     int
	ID        uint64
}

func fetchNameFromKey(key string) string {
	var Name string
	boltDBinstance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("key2name"))
		Name = string(b.Get([]byte(key)))

		return nil
	})
	return Name[4:len(Name)]
}

//loadEntries will pull the tables in a bucket
func loadEntries(bucket string) (Database, error) {
	var entry Entry
	database := NewDatabase(bucket)

	err := boltDBinstance.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {

			// Load the entry from binary
			jsonErr := json.Unmarshal(v, &entry)
			if jsonErr != nil {
				fmt.Println(jsonErr, "unMarshalling")
			}
			// Load em up cowboy
			table, exists := database.Tables[entry.Table]

			var column Column

			column.Name = entry.Column
			column.Type = entry.Type
			column.Ordinal = entry.Ordinal
			column.Length = entry.Length
			column.Precision = entry.Precision
			column.Scale = entry.Scale
			column.TableName = table.Name
			column.DatabaseName = database.Name
			column.ID = entry.ID

			fmt.Println(column.ID)
			if !exists {
				table = NewTable(entry.Table, entry.Schema, database.Name)

			}
			table.Columns[column.ID] = column

			database.Tables[entry.Table] = table

		}
		return nil
	})

	if err != nil {
		log.Println(err, "loadEntries")
	}

	return *database, err
}

// Templates setup
func listDBhandler(w http.ResponseWriter, r *http.Request) {
	var databaseList DatabaseList
	databaseList.Databases = make(map[string]string)
	err := boltDBinstance.View(func(tx *bolt.Tx) error {
		names := tx.Bucket([]byte("name2key")).Cursor()
		prefix := []byte("dbs:")
		for k, v := names.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = names.Next() {
			key := strings.Split(string(k), ":")[1]
			value := string(v)
			log.Println("Key:", key)
			log.Println("Value:", value)

			databaseList.Databases[key] = value
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})
	if err != nil {
		log.Println(err, "loadEntries")
	}

	templates := template.Must(template.ParseFiles("templates/databases.html", "templates/header.html", "templates/footer.html"))

	err = templates.Execute(w, databaseList)
	if err != nil {
		fmt.Println(err)
	}
}

func singleTBLhandler(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/singleTable.html", "templates/header.html", "templates/footer.html"))
	var err error
	var columnList ColumnList
	// var columnName []string
	tableKey := r.URL.Path[len("/tbl/"):]
	columnList.TableKey = tableKey
	columnList.TableName = fetchNameFromKey(tableKey)
	columnList.DatabaseKey = tableKey[:4]
	columnList.DatabaseName = fetchNameFromKey(tableKey[:4])
	columnList.Columns = make(map[string]string)
	err = boltDBinstance.View(func(tx *bolt.Tx) error {
		columns := tx.Bucket([]byte("key2name")).Cursor()
		prefix := []byte(tableKey)
		for k, _ := columns.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = columns.Next() {
			if len(string(k)) == 12 {
				columnKey := string(k)
				columnName := fetchNameFromKey(columnKey)
				columnList.Columns[columnName] = columnKey
			}
		}
		return nil
	})
	err = templates.Execute(w, columnList)
	if err != nil {
		fmt.Println(err)
	}
}

func singleDBShandler(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/singleDatabase.html", "templates/header.html", "templates/footer.html"))
	dbsKey := r.URL.Path[len("/dbs/"):]
	var tableList TableList
	tableList.DatabaseName = fetchNameFromKey(dbsKey)
	tableList.Tables = make(map[string]string)
	err := boltDBinstance.View(func(tx *bolt.Tx) error {
		columns := tx.Bucket([]byte("key2name")).Cursor()
		prefix := []byte(dbsKey)
		for k, _ := columns.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = columns.Next() {
			if len(string(k)) == 8 {
				tableKey := string(k)
				tableName := fetchNameFromKey(tableKey)
				tableList.Tables[tableName] = tableKey
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err, "loadEntries")
	}
	log.Println("Table Count:", len(tableList.Tables))
	templates.Execute(w, tableList)
}

func singleColhandler(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/singleColumn.html", "templates/header.html", "templates/footer.html"))

	columnName := r.URL.Path[len("/cl/"):]
	log.Println(columnName)
	column := database.Tables["fred"].Columns[2]

	err := templates.Execute(w, column)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	var err error
	if _, err = os.Stat(dbFilePath); os.IsNotExist(err) {
		fmt.Println(err)
	}

	boltDBinstance, err = bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		log.Println(err)
	}

	defer boltDBinstance.Close()

	http.HandleFunc("/", listDBhandler)

	http.HandleFunc("/col/", singleColhandler)
	http.HandleFunc("/dbs/", singleDBShandler)
	http.HandleFunc("/tbl/", singleTBLhandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":3001", nil)
}
