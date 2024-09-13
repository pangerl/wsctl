package tenant

//func CheckErr(err error) {
//	if err != nil {
//		//log.Printf("Failed info: %s \n", err)
//		log.Fatalf("Failed info: %s \n", err)
//	}
//}

//func CurrentMessageNum(client *elastic.Client, corpid string, dateNow time.Time) int64 {
//	// 统计今天的会话数
//	startTime := getZeroTime(dateNow).UnixNano() / 1e6
//	endTime := dateNow.UnixNano() / 1e6
//	messagenum, _ := countMessageNum(client, corpid, startTime, endTime)
//	return messagenum
//}

//func getZeroTime(d time.Time) time.Time {
//	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
//}
