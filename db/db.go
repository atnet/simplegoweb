package db

import (
	"database/sql"
	"fmt"
	_ "mysql"
	_ "reflect"
)

//数据库连接器
type Connector struct {
	DriverName		string		//驱动名称
	DriverSource	string		//驱动连接地址
}

//获取数据库连接
func (connector *Connector) Open() *sql.DB {
	return connector.open()
}

func (connector *Connector) open() *sql.DB {
	
	db, err := sql.Open(connector.DriverName, connector.DriverSource)
	if err != nil {
		defer db.Close()
		//连接成功 err一定是nil否则就是报错
		panic(err.Error()) //抛出异常
		//fmt.Println(err.Error())//仅仅是显示异常
	}
	return db
}

func (connector *Connector) Query(sql string, f func(*sql.Rows), arg ...interface{}) {

	db := connector.open()
	defer db.Close()
	stmt, _ := db.Prepare(sql)

	rows, err := stmt.Query(arg...)
	if err != nil {
		fmt.Println("SQL:" + sql)
		fmt.Println("ERR:" + err.Error())
	}
	if f != nil {
		f(rows)
	}
}

//查询Rows
func  (connector *Connector) QueryRow(sql string, f func(*sql.Row), arg ...interface{}) {
	db := connector.open()
	defer db.Close()
	stmt, _ := db.Prepare(sql)
	row := stmt.QueryRow(arg...)

	if f != nil && row != nil {
		f(row)
	}
}

//执行
func  (connector *Connector) Exec(sql string,args ...interface{})(rows int,lastInsertId int,err error){
	db := connector.open()
	defer db.Close()
	stmt,err := db.Prepare(sql)
	if err!= nil {
		return 0,-1,err
	}
	result,err := stmt.Exec(args...)
	if err != nil {
		return 0,-1,err
	}
	lastId,_ := result.LastInsertId()
	arows,_ := result.RowsAffected()
	
	return int(arows),int(lastId),nil
}

//转换为字典数组
//参考：http://my.oschina.net/nowayout/blog/143278
func ConvRowsToMapForJson(rows *sql.Rows) (rowsMap []map[string]interface{}){
	rowsMap = [](map[string]interface{}){} 		  //数据切片
	var tmpInt int = 0                            //序列
	columns, _ := rows.Columns()                  //列名
	
	//定义数组，数组的类型为[]byte
	var values []interface{} = make([]interface{}, len(columns))
	var rawBytes [][]byte = make([][]byte,len(values))
	
	for v := range values {
		values[v] = &rawBytes[v]
	}
	
	for rows.Next() {
		rows.Scan(values...)

		if len(rowsMap) == tmpInt {
			rowsMap = append(rowsMap, make(map[string]interface{}))
		}
			
		for i, v := range columns {
			rowsMap[tmpInt][v] = string(rawBytes[i])
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
		tmpInt++
	}
	return rowsMap
}

func ConvRowsToMap(rows *sql.Rows) (rowsMap []map[string][]byte){
	rowsMap = [](map[string][]byte){} 		  //数据切片
	var tmpInt int = 0                            //序列
	columns, _ := rows.Columns()                  //列名
	
	//定义数组，数组的类型为[]byte
	var values []interface{} = make([]interface{}, len(columns))
	var rawBytes [][]byte = make([][]byte,len(values))
	
	for v := range values {
		values[v] = &rawBytes[v]
	}
	
	for rows.Next() {
		rows.Scan(values...)

		if len(rowsMap) == tmpInt {
			rowsMap = append(rowsMap, make(map[string][]byte))
		}
			
		for i, v := range columns {
			rowsMap[tmpInt][v] =rawBytes[i]
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
		tmpInt++
	}
	return rowsMap
}

func ConvSqlRowToMap(rows *sql.Rows)(map[string][]byte){
	rowMap := make(map[string][]byte)
	columns, _ := rows.Columns()              //列名
	if rows.Next() {
		row := rows
		 	  //数据
		//定义数组，数组的类型为[]byte
		var values []interface{} = make([]interface{}, len(columns))
		var rawBytes [][]byte = make([][]byte,len(values))
		for v := range values {
			values[v] = &rawBytes[v]
		}
		row.Scan(values...)
		for i, v := range columns {
			rowMap[v] = rawBytes[i]
			//fmt.Println(v + "===>" + string(rawBytes[i]))
		}
	}
	return rowMap
}
