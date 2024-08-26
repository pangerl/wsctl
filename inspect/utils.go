package inspect

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"math/rand"
	"strconv"
	"time"
)

//func CheckErr(err error) {
//	if err != nil {
//		//log.Printf("Failed info: %s \n", err)
//		log.Fatalf("Failed info: %s \n", err)
//	}
//}

func CurrentMessageNum(client *elastic.Client, corpid string, dateNow time.Time) int64 {
	// 统计今天的会话数
	startTime := getZeroTime(dateNow).UnixNano() / 1e6
	endTime := dateNow.UnixNano() / 1e6
	messagenum, _ := countMessageNum(client, corpid, startTime, endTime)
	return messagenum
}

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

func GetRandomDuration() time.Duration {
	// 创建一个新的随机数生成器，使用当前时间作为种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// 生成随机数
	randomSeconds := r.Intn(300)
	// 将随机秒数转换为时间.Duration
	duration := time.Duration(randomSeconds) * time.Second
	return duration
}
