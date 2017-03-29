package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gocarina/gocsv"
)

// Database is the basis of all the systems
type Database struct {
	Name       string
	Tables     map[string]Table
	Server     string
	Type       string
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

// Entry is a single line of a database/table definition
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
	ID        uint64
}

// UID returns a unique identifier for the entry
func (e *Entry) UID() string {
	return e.System + "|" + e.Schema + "|" + e.Table + "|" + e.Column
}

// fourLetterGenerator creates a randomly string of four characters, upper and lower
// case
func fourLetterGenerator() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	char1 := rand.Intn(52)
	char2 := rand.Intn(52)
	char3 := rand.Intn(52)
	char4 := rand.Intn(52)

	value := string(chars[char1]) + string(chars[char2]) + string(chars[char3]) + string(chars[char4])
	return value
}

func findDBSKey(name string, bucketDB *bolt.Bucket) string {
	c := bucketDB.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if string(v) == name {
			return string(k)
		}
	}
	return ""
}

func newDBSKey(name string, bucketDB *bolt.Bucket) string {
	newKey := fourLetterGenerator()

	matched := bucketDB.Get([]byte(newKey))

	keyExists := (matched != nil)
	for keyExists {
		newKey = fourLetterGenerator()
		matched = bucketDB.Get([]byte(newKey))
		keyExists = (matched != nil)
	}

	return newKey
}

func main() {

	// Load Definition file from csv
	definitionFile, err := os.OpenFile("../psqlmetadata.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer definitionFile.Close()

	entries := []*Entry{}

	if err = gocsv.UnmarshalFile(definitionFile, &entries); err != nil { // Load entries from file
		panic(err)
	}

	// Send Entries to Bolt DB

	db, err := bolt.Open("datagrammar.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("ready to load")
	databaseName := entries[0].Database

	err = db.Update(func(tx *bolt.Tx) error {
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		tx.DeleteBucket([]byte(databaseName))

		log.Println("Opening the bucket")

		bucketDB, err := tx.CreateBucketIfNotExists([]byte("DBS"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		log.Println("Datbase name is ", databaseName)
		dbKey := findDBSKey(databaseName, bucketDB)
		if dbKey == "" {
			dbKey = newDBSKey(databaseName, bucketDB)
		}

		log.Println("DB Key is ", dbKey)
		b, err := tx.CreateBucket([]byte(databaseName))
		if err != nil {
			log.Println(err)
			return fmt.Errorf("create bucket: %s", err)
		}

		for _, entry := range entries {
			entry.ID, _ = b.NextSequence()
			log.Println(entry.ID)
			encoded, err := json.Marshal(entry)
			if err != nil {
				return err
			}

			err = b.Put([]byte(entry.UID()), encoded)
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}

		return nil
	})

}
