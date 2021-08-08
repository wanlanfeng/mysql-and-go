package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//Db数据库连接池
var db *sql.DB

//提前设置常量
const(
	//数据库的用户名
	userName = "root"
	//密码
	password = "1925688295WLFwlf"
	//ip地址
	ip = "127.0.0.1"
	//端口
	port = "3306"
	//数据库名称
	dbName = "lanshan_task"
)

//定义仓库结构体
type Warehouses struct {
	warehouse_code string
	warehouse_capacity   int
}

//定义服装结构体
type Costumes struct{
	id string
	size string
	price int
	typ string
}

func Initdb(){
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库，返回db
	db, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	if err := db.Ping();err != nil{
		fmt.Println("open database fail")
		return
	}
	fmt.Println("open success")
	return
}

//定义供应商结构体
type Providers struct {
	id string
	name  string
}

//定义供应情况结构体
type Provide_case struct {
	provider_id string
	costume_id string
	costume_level string
}
//查询
func Query(){

	//查询服装尺码为'S'且销售价格在100以下的服装信息。
	fmt.Println("查询服装尺码为'S'且销售价格在100以下的服装信息。")
	var cos Costumes
	rows, err := db.Query("select * from Costumes where price < 100 and size = 'S';")
	if err != nil{
		fmt.Println(errors.New("query occur error"))
	}
	for rows.Next(){
		err := rows.Scan(&cos.id,&cos.size,&cos.price,&cos.typ)
		if err != nil{
			fmt.Printf("scan falied")
		}
		fmt.Printf("%#v\n",cos)
	}
	rows.Close()

	//查询仓库容量最大的仓库信息
	fmt.Println("查询仓库容量最大的仓库信息")
	var w Warehouses
	rows, err = db.Query("select * from warehouses where warehouse_capacity = (select max(warehouse_capacity) from warehouses);")
	if err != nil{
		fmt.Println(errors.New("query occur error"))
	}
	for rows.Next(){
		err := rows.Scan(&w.warehouse_code, &w.warehouse_capacity)
		if err != nil{
			fmt.Printf("scan failed")
		}
		fmt.Printf("%#v\n", w)
	}
	rows.Close()

	//查询服装编码以‘A’开始开头的服装。
	fmt.Println("查询服装编码以‘A’开始开头的服装。")
	var cos2 Costumes
	rows, err = db.Query("select * from costumes where id like 'A%';")
	if err != nil{
		fmt.Println(errors.New("query occur error"))
	}
	for rows.Next(){
		err := rows.Scan(&cos2.id, &cos2.size,&cos2.price, &cos2.typ)
		if err != nil{
			fmt.Printf("scan failed")
		}
		fmt.Printf("%#v\n", cos2)
	}
	rows.Close()

	//查询服装质量等级有不合格的供应商信息。
	fmt.Println("查询服装质量等级有不合格的供应商信息。")
	var p Provide_case
	rows, err = db.Query("select provider_id from provide_case where costume_level = 'reject';")
	if err != nil{
		fmt.Println(errors.New("query occur error"))
	}
	for rows.Next(){
		err = rows.Scan(&p.provider_id)
		if err != nil{
			fmt.Printf("scan failed")
		}
		fmt.Printf("%#v\n", p.provider_id)
	}
	rows.Close()

	//查询A类服装的库存总量。
	var Anum int
	fmt.Println("查询A类服装的库存总量。 ")
	rows, err = db.Query("select count(*) as Anum from costumes where id like 'A%';")
	if err != nil{
		fmt.Println(errors.New("query occur error"))
	}
	for rows.Next(){
		err := rows.Scan(&Anum)
		if err != nil{
			fmt.Printf("scan failed")
		}
		fmt.Printf("%#v\n", Anum)
	}
	rows.Close()
}

//删除所有服装质量等级不合格的供应情况
func Delete(){
	fmt.Println("删除所有服装质量等级不合格的供应情况")
	stmt, err := db.Prepare("delete from provide_case where costume_level = ?")
	if err != nil{
		fmt.Println("delete failed")
		return
	}
	res, err := stmt.Exec("reject")
	if err != nil{
		fmt.Println("delete failed")
		return
	}
	num, err := res.RowsAffected()
	if err != nil{
		fmt.Println("delete failed")
		return
	}
	fmt.Printf("删除的行数为：%d\n", num)
	return
}

//向每张表插入一条记录
func Insert(){
	//向warehouses插入数据
	fmt.Println("向warehouses插入数据:")
	res1, err := db.Exec("insert into warehouses values('wlf3',30);" )
	if err != nil {
		fmt.Printf("insert failed %v\n", err)
		return
	}
	code, err := res1.LastInsertId()
	if err != nil{
		fmt.Println("get lastinsertid failed")
		return
	}
	fmt.Println("insert success, the code is: ", code)

	//向Costumes插入数据
	fmt.Println("向Costumes插入数据:")
	res2, err := db.Exec("insert into costumes values('C1','L',110,'B');" )
	if err != nil {
		fmt.Printf("insert failed %v\n", err)
		return
	}
	id, err := res2.LastInsertId()
	if err != nil{
		fmt.Println("get lastinsertid failed")
		return
	}
	fmt.Println("insert success, the code is: ", id)

	//向Providers插入数据
	fmt.Println("向Providers插入数据:")
	res3, err := db.Exec("insert into Providers values('2001','taobao');" )
	if err != nil {
		fmt.Printf("insert failed %v\n", err)
		return
	}
	id2, err := res3.LastInsertId()
	if err != nil{
		fmt.Println("get lastinsertid failed")
		return
	}
	fmt.Println("insert success, the code is: ", id2)
	//向Provide_case插入数据
	fmt.Println("向Provide_case插入数据:")
	res4, err := db.Exec("insert into Provide_case values('1','qiaodan','qualify');" )
	if err != nil {
		fmt.Printf("insert failed %v\n", err)
		return
	}
	costume_id, err := res4.LastInsertId()
	if err != nil{
		fmt.Println("get lastinsertid failed")
		return
	}
	fmt.Println("insert success, the code is: ", costume_id)
	return
}

//更新
func Update(){
	//把服装尺寸为'S'的服装的销售价格均在原来基础上提高10%。
	fmt.Println("把服装尺寸为'S'的服装的销售价格均在原来基础上提高10%。")
	res, err := db.Exec("update costumes set price = price * 1.1 where size = 'S';")
	if err != nil{
		println("update failed")
		return
	}
	num, err := res.RowsAffected()
	fmt.Println("更新的行数为：", num)
	return
}

//主函数
func main() {
	Initdb()
	Query()
	Delete()
	Insert()
	Update()
	return
}