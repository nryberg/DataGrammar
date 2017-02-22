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
var majorKeys Branch
var system System

// BucketList contains a list of the buckets
type BucketList struct {
	Buckets []string
}

// System contains a list of the systems and a count of tables
type System struct {
	name   string
	tables []Table
}

// Table is the tables in the systems
type Table struct {
	name    string
	columns []Column
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

// Branch is a general type for a tree form
type Branch struct {
	Name  string
	Items []string
}

// Entry is the single database entry
type Entry struct { // Our example struct, you can use "-" to ignore a field
	Database  string `csv:"database"`
	System    string `csv:"system"`
	Schema    string `csv:"schema"`
	Table     string `csv:"table"`
	Column    string `csv:"column"`
	Ordinal   int    `csv:"Ordinal"`
	Type      string `csv:"Type"`
	Length    int    `csv:"Length"`
	Precision int    `csv:"Precision"`
	Scale     int    `csv:"Scale"`
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

// Buckets prints a list of all buckets.
func Buckets() (Branch, error) {
	var bucketList Branch
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
			bucketList.Items = append(bucketList.Items, string(name))
			return nil
		})
	})

	return bucketList, err
}

//MajorKeys will pull the systems in a bucket
func loadMajorKeys(bucket string) error {
	var entry Entry
	majorKeys.Name = bucket
	boltDBinstance, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		fmt.Println(err, "open Bolt DB")
	}

	defer boltDBinstance.Close()
	err = boltDBinstance.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			majorKeys.Items = append(majorKeys.Items, string(k))

			jsonErr := json.Unmarshal(v, &entry)
			if jsonErr != nil {
				fmt.Println(err, "unMarshalling")
			}

		}
		return nil
	})

	if err != nil {
		fmt.Println(err, "loadMajorKeys")
	}

	return err
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
	err := loadMajorKeys(dbName)

	if err != nil {
		fmt.Println(err)
	}

	err = templates.Execute(w, majorKeys)
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
