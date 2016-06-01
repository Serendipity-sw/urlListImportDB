package main

import (
	"database/sql"
	"github.com/guotie/config"
	"fmt"
	"github.com/smtc/glog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guotie/deferinit"
	"runtime"
	"github.com/go-sql-driver/mysql"
)

var(
	dbs       *sql.DB
)

func init() {
	deferinit.AddInit(sqlConntion, sqlClose, 999)
}
/**
数据库连接打开
创建人:邵炜
创建时间:2015年12月29日17:25:57
输入参数:无
输出参数:无
*/
func sqlConntion() {
	var (
		err error
	)
	dbuser := config.GetStringMust("navdbuser")
	dbhost := config.GetStringMust("navdbhost")
	dbport := config.GetIntDefault("navdbport", 3306)
	dbpass := config.GetStringMust("navdbpassword")
	dbname := config.GetStringDefault("navdbname", "mdc")

	dbclause := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", dbuser, dbpass, dbhost, dbport, dbname)

	dbs, err = sql.Open("mysql", dbclause)
	if err != nil {
		glog.Error("mysql can't connection %s \n", err.Error())
		return //nil, err
	}

	err = dbs.Ping()

	if err != nil {
		glog.Error("mysql can't ping %s \n", err.Error())
		return //nil, err
	}
}

/**
数据库连接关闭
创建人:邵炜
创建时间:2015年12月30日11:08:17
输入参数:无
输出参数:无
*/
func sqlClose() {
	dbs.Close()
}

func loadFileDB(filePath string) {
	wrap:="\n"
	if runtime.GOOS == "windows" {
		wrap="\r\n"
	}
	mysql.RegisterLocalFile(filePath)
	loadSql:="load data LOW_PRIORITY local infile '%s' into table url FIELDS TERMINATED BY ','  LINES TERMINATED BY '%s' (Url,`Name`)  set UrlGroupId=1"
	loadSql=fmt.Sprintf(loadSql,filePath,wrap)
	_,err:=dbs.Exec(loadSql)
	if err != nil {
		glog.Error("sql is warn, sql: %s err: %s \n",loadSql,err.Error())
		sqlClose()
		sqlConntion()
		return
	}
	glog.Info("db run success! filePath: %s \n",filePath)
}