package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

func hello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	bucketList, err := Buckets("datagrammar.db")
	if err != nil {
		fmt.Println(err)
	}

	head := `<DOCTYPE html>
<html>
<head>
    <title></title>
      </head>
      <body>
          Hello World!`
	tail := `
    </body>
  </html>`
	io.WriteString(
		res,
		head+bucketList[0]+tail,
	)
}
func main() {
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":9000", nil)
}
