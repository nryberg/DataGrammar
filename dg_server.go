package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
)

const dbFilePath string = "datagrammar.db"

var boltDBinstance bolt.DB

// BucketList contains a list of the buckets
type BucketList struct {
	Buckets []string
}

// Database contains a list of schemas
type Database struct {
	name   string
	tables map[string]Table
	server string
}

// NewDatabase initializes the Database struct with a map
func NewDatabase(name string) *Database {
	db := Database{
		tables: make(map[string]Table),
		name:   name,
		server: "",
	}

	return &db
}

// Table is the tables in the systems
type Table struct {
	name    string
	columns map[string]Column
	schema  string
}

// NewTable initializes the Database struct with a map
func NewTable() *Table {
	tb := Table{
		columns: make(map[string]Column),
		name:    "",
		schema:  "",
	}

	return &tb
}

// Column is the base type here
type Column struct {
	name      string
	Ordinal   int
	Type      string
	Length    int
	Precision int
	Scale     int
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
			//majorKeys.Items = append(majorKeys.Items, string(k))

			jsonErr := json.Unmarshal(v, &entry)
			if jsonErr != nil {
				fmt.Println(err, "unMarshalling")
			}

			// Load em up cowboy

			table := database.tables[entry.Table]
			log.Println(table.name)

		}
		return nil
	})

	if err != nil {
		fmt.Println(err, "loadMajorKeys")
	}

	return *database, err
}

// Templates setup

func listDBhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("In listDBhandler")
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

func singleDBhandler(w http.ResponseWriter, r *http.Request) {
	log.Println("In single DB Handler")
	templates := template.Must(template.ParseFiles("templates/singleDatabase.html", "templates/header.html", "templates/footer.html"))

	dbName := r.URL.Path[len("/db/"):]
	database, err := loadEntries(dbName)

	if err != nil {
		fmt.Println(err)
	}

	err = templates.Execute(w, database)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	http.HandleFunc("/", listDBhandler)
	http.HandleFunc("/db/", singleDBhandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":9000", nil)
}
