package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gocarina/gocsv"
)

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
}

// UID returns a unique identifier for the entry
func (e *Entry) UID() string {
	return e.System + "|" + e.Schema + "|" + e.Table + "|" + e.Column
}

func main() {

	// Load Definition file from csv
	definitionFile, err := os.OpenFile("psqlmetadata.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer definitionFile.Close()

	entries := []*Entry{}

	if err = gocsv.UnmarshalFile(definitionFile, &entries); err != nil { // Load entries from file
		panic(err)
	}
	for _, entry := range entries {
		fmt.Println("Hello", entry.UID())
	}

	// Send Entries to Bolt DB

	db, err := bolt.Open("datagrammar.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db_name := entries[0].Database

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(db_name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		for _, entry := range entries {
			encoded, err := json.Marshal(entry)
			if err != nil {
				return err
			}

			err = b.Put([]byte(entry.UID()), encoded)
		}

		return nil
	})

}
