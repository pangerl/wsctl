package inspect

import (
	"math/rand"
	"time"
)

//func CheckErr(err error) {
//	if err != nil {
//		//log.Printf("Failed info: %s \n", err)
//		log.Fatalf("Failed info: %s \n", err)
//	}
//}

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
