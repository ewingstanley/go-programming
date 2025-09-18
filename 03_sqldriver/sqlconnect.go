package main

import ("database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)


func main() {
	dsn := "root:123456@tcp(localhost:3306)/dbtest"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("dsn格式错误：", err)
	}
	
	 err = db.Ping()
	if err != nil{
		fmt.Println("数据库链接失败：", err)
	}

	defer db.Close()

	// res,_ :=db.Query("select 标题,制定机关 from fabaodata limit 10")

	// defer res.Close()

	// for res.Next(){
	// 	var title,organ string
	// 	res.Scan(&title,&organ)
	// 	fmt.Println(title,organ)
	// }

	var title1,organ1 string
	db.QueryRow("select 标题,制定机关 from fabaodata limit 10").Scan(&title1,&organ1)
	fmt.Println(title1,organ1)

}