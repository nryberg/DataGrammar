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
	TableKey     string
	DatabaseKey  string
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
	Key       string
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

			databaseList.Databases[key] = value

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
	if err != nil {
		log.Println("Error:", err)
	}
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
	var column Column
	var entry Entry
	columnKey := r.URL.Path[len("/col/"):]
	boltDBinstance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("column"))

		encoded := b.Get([]byte(columnKey))
		log.Println(string(encoded))
		errUnmarshal := json.Unmarshal(encoded, &entry)
		if errUnmarshal != nil {
			log.Println("Unmarshalling error:", errUnmarshal)
		}

		log.Println("Entry:", entry)
		column.Name = entry.Column
		column.Length = entry.Length
		column.DatabaseName = fetchNameFromKey(columnKey[:4])
		column.DatabaseKey = columnKey[:4]
		column.TableName = fetchNameFromKey(columnKey[:8])
		column.TableKey = columnKey[:8]
		column.Type = entry.Type
		return nil
	})
	log.Println("Column Name:", entry.Column)
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
