package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

// Entry is a single line of a database/table definition
type Entry struct { // Our example struct, you can use "-" to ignore a field
	System string `csv:"system"`
	Schema string `csv:"schema"`
	Table  string `csv:"table"`
	Column string `csv:"column"`
}

// UID returns a unique identifier for the entry
func (e *Entry) UID() string {
	return e.System + "|" + e.Schema + "|" + e.Table + "|" + e.Column
}

func main() {
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

	if _, err = definitionFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}

}
