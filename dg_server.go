package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
)

// Buckets prints a list of all buckets.
func Buckets(path string) ([]string, error) {

	bucketList := []string{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return nil, err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			bucketList = append(bucketList, string(name))
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return bucketList, err
}

func databaseHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/databases.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, r)
}

func main() {
	http.HandleFunc("/databases", databaseHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":9000", nil)

	bucketList, err := Buckets("datagrammar.db")
	log.Println(len(bucketList))
	if err != nil {
		fmt.Println(err)
	}

}
