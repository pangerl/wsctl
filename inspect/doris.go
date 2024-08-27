// Package inspect @Author lanpang
// @Date 2024/8/23 下午7:02:00
// @Desc
package inspect

// 查询多行数据
//func selectFailedJob(id int, db *sql.DB) {
//	sqlStr := "select id,name,age from user where id>?"
//	rows, err := db.Query(sqlStr, id)
//	if err != nil {
//		fmt.Println("数据查询失败. err:", err)
//	}
//	defer rows.Close()
//	// 循环读数据
//	for rows.Next() {
//		var u user
//		err := rows.Scan(&u.id, &u.name, &u.age)
//		if err != nil {
//			fmt.Println("数据读取失败. err:", err)
//			return
//		}
//		// 输出数据
//		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
//	}
//}
