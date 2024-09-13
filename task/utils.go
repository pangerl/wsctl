// Package task @Author lanpang
// @Date 2024/9/13 下午3:44:00
// @Desc
package task

import (
	"fmt"
	"time"
)

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func CallUser(users []string) string {
	var result string
	if len(users) == 0 {
		return result
	}
	for _, user := range users {
		result += fmt.Sprintf("<@%s>", user)
	}
	return result
}
