package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/template"
	"github.com/boltdb/bolt"
)

// BucketList contains a list of the buckets
type BucketList struct {
	Buckets []string
}

// Buckets prints a list of all buckets.
func Buckets(path string) (BucketList, error) {
	var bucketList BucketList

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return bucketList, err
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return bucketList, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			bucketList.Buckets = append(bucketList.Buckets, string(name))
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
		return bucketList, err
	}

	return bucketList, err
}

// Templates setup

func databaseHandler(w http.ResponseWriter, r *http.Request) {
	bucketList, err := Buckets("datagrammar.db")
	templates := template.Must(template.ParseFiles("templates/databases.html", "templates/header.html", "templates/footer.html"))

	if err != nil {
		fmt.Println(err)
	}

	err = templates.Execute(w, bucketList)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/databases", databaseHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":9000", nil)
}
