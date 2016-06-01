package main

import (
	"bufio"
	"io"
	"os"
	"github.com/smtc/glog"
	"runtime"
	"strings"
)

/**
打开文件并处理
创建人:邵炜
创建时间:2016年6月1日09:40:03
输入参数:filePath 文件地址  readOrWrite 读还是写文件  true是读  false写
输出参数:文件对象 错误对象
 */
func openFile(filePath string,readOrWrite bool) (*os.File,error) {
	var (
		fs *os.File
		err error
	)
	if readOrWrite {
		fs,err=os.Open(filePath)
	}else {
		fs,err=os.OpenFile(filePath,os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	}
	if err != nil {
		glog.Error("open file is error! filePath: %s err: %s \n",filePath,err.Error())
		return nil,err
	}
	glog.Info("file open success! filePath: %s \n",filePath)
	return fs,nil
}
/**
逐行写文件
创建人:邵炜
创建时间:2016年6月1日09:45:11
 */
func writeFile() {

}

/**
逐行读文件
创建人:邵炜
创建时间:2016年6月1日09:49:45
输入参数:文件地址  是否写入文件
 */
func readFile(filepath string,isWrite bool) {
	var (
		readAll =false
		readByte []byte
		line []byte
		err error
		write *os.File
		wrap []byte=[]byte("\n")
		sqlStr []byte
	)
	if runtime.GOOS == "windows" {
		wrap=[]byte("\r\n")
	}
	if isWrite {
		write,err=openFile("./afterProcess.txt",false)
		if err != nil {
			return
		}
		defer write.Close()
	}
	read,err:=openFile(filepath,true)
	if err != nil {
		return
	}
	defer read.Close()
	buf:=bufio.NewReader(read)
	for err!=io.EOF {
		if err!=nil {
			glog.Error("read error! err: %s \n",err.Error())
		}
		if readAll {
			readByte,readAll,err=buf.ReadLine()
			line=append(line,readByte...)
		}else{
			readByte,readAll,err=buf.ReadLine()
			line=append(line,readByte...)
			if len(strings.TrimSpace(string(line)))==0 {
				continue
			}
			sqlStr=append(sqlStr,[]byte("INSERT INTO `dspadmin`.`url`(`Url`,`UrlGroupId`,`Name`,`IsDelete`)VALUES('")...)
			sqlStr=append(sqlStr,line...)
			sqlStr=append(sqlStr,[]byte("',1,'',0);")...)
			//line=append(line,[]byte(",随机网址")...)
			if err != io.EOF {
				sqlStr=append(sqlStr,wrap...)
			}
			//写文件
			if isWrite {
				write.Write(sqlStr)
			}
			line=line[:0]
			sqlStr=sqlStr[:0]
		}
	}
}