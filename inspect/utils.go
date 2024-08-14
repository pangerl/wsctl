package inspect

import (
	"fmt"
	"strconv"
	"time"
)

//func CheckErr(err error) {
//	if err != nil {
//		//log.Printf("Failed info: %s \n", err)
//		log.Fatalf("Failed info: %s \n", err)
//	}
//}

func getZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func calluser(users []string) string {
	var result string
	if len(users) == 0 {
		return result
	}
	for _, user := range users {
		result += fmt.Sprintf("<@%s>", user)
	}
	return result
}

func getRole(role string) string {
	if role == "0" {
		return "Master"
	}
	return "Slave"
}

func convertAndCalculate(str1, str2 string) (int, error) {
	num1, err := strconv.Atoi(str1)
	if err != nil {
		return 0, err
	}

	num2, err := strconv.Atoi(str2)
	if err != nil {
		return 0, err
	}

	return num1 - num2, nil
}
