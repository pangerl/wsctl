package inspect

import (
	"log"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		//log.Printf("Failed info: %s \n", err)
		log.Fatalf("Failed info: %s \n", err)

	}
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
