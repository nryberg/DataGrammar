package main

import (
	"log"
	"math/rand"
	"time"
)

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

func main() {
	test := fourLetterGenerator()
	log.Println(test)
}
