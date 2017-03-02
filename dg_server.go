package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
)

const dbFilePath string = "datagrammar.db"

var boltDBinstance bolt.DB
var database Database

// BucketList contains a list of the buckets
type BucketList struct {
	Buckets []string
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
	Columns      map[string]Column
	Schema       string
	ColumnNames  []string
	DatabaseName string
}

// NewTable initializes the Database struct with a map
func NewTable(name, schema string, dbName string) Table {
	tb := Table{
		Columns:      make(map[string]Column),
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

func openDB(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return err
	}

	boltDBinstance, err := bolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer boltDBinstance.Close()
	return err
}

// Buckets loads a list of all buckets.
func Buckets() (BucketList, error) {
	var bucketList BucketList
	boltDBinstance, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return bucketList, err
	}
	defer boltDBinstance.Close()
	if err != nil {
		fmt.Println(err)
		return bucketList, err
	}
	err = boltDBinstance.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			bucketList.Buckets = append(bucketList.Buckets, string(name))
			return nil
		})
	})

	return bucketList, err
}

// buildDBtree converts majorKeys to a tree of entries
func buildDBtree() error {
	var err error
	return err
}

//loadEntries will pull the tables in a bucket
func loadEntries(bucket string) (Database, error) {
	var entry Entry
	database := NewDatabase(bucket)

	boltDBinstance, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		fmt.Println(err, "open Bolt DB")
	}

	defer boltDBinstance.Close()

	err = boltDBinstance.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {

			// Load the entry from binary
			jsonErr := json.Unmarshal(v, &entry)
			if jsonErr != nil {
				fmt.Println(err, "unMarshalling")
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
			table.Columns[column.Name] = column

			database.Tables[entry.Table] = table

		}
		return nil
	})

	if err != nil {
		fmt.Println(err, "loadEntries")
	}

	return *database, err
}

// Templates setup

func listDBhandler(w http.ResponseWriter, r *http.Request) {
	bucketList, err := Buckets()
	if err != nil {
		fmt.Println(err, "listDBHandler")
	}

	templates := template.Must(template.ParseFiles("templates/databases.html", "templates/header.html", "templates/footer.html"))

	err = templates.Execute(w, bucketList)
	if err != nil {
		fmt.Println(err)
	}
}

func singleTBhandler(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/singleTable.html", "templates/header.html", "templates/footer.html"))

	tableName := r.URL.Path[len("/tb/"):]
	table := database.Tables[tableName]

	var columnName []string

	for k := range table.Columns {
		columnName = append(columnName, k)

	}

	sort.Strings(columnName)
	table.ColumnNames = columnName

	err := templates.Execute(w, table)
	if err != nil {
		fmt.Println(err)
	}

}

func singleColhandler(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/singleColumn.html", "templates/header.html", "templates/footer.html"))

	columnName := r.URL.Path[len("/cl/"):]
	column := database.Tables["fred"].Columns[columnName]

	err := templates.Execute(w, column)
	if err != nil {
		fmt.Println(err)
	}

}
func singleDBhandler(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/singleDatabase.html", "templates/header.html", "templates/footer.html"))
	dbName := r.URL.Path[len("/db/"):]
	database, _ = loadEntries(dbName)

	for k := range database.Tables {
		database.TableNames = append(database.TableNames, k)
	}

	sort.Strings(database.TableNames)

	templates.Execute(w, database)

}

func main() {

	http.HandleFunc("/", listDBhandler)

	http.HandleFunc("/cl/", singleColhandler)
	http.HandleFunc("/db/", singleDBhandler)
	http.HandleFunc("/tb/", singleTBhandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":3001", nil)
}
