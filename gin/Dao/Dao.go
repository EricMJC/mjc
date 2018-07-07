package Dao

import (
	"database/sql"
	"strings"
	"fmt"
	"github.com/satori/go.uuid"
)

const   (
	userName="root"
	password="173167479"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "mysqlTmd"
)
type MySQL struct {
	DB *sql.DB
	tx *sql.Tx
	Table string
	order string
	limit string
	field string
	err  error
	maxConnMaxLifetime int
	maxIdleConns int
}


type IMG struct {
	I_id    int
	ImgPath string
}

type DAO struct {
	DB *sql.DB
	Img IMG
}
func InitDao() (*DAO , error ) {
	var err error
	var dao DAO
	DB := dao.DB
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, err = sql.Open("mysql", path)
	if err != nil{
		fmt.Println(err)
		return nil,err
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database failed")
		fmt.Println(err)
		return nil,err
	}
	fmt.Println("connect success")
	dao.DB=DB
	return &dao,nil
}
func (this *DAO)Insert(dao *DAO,imageId string,u1 uuid.UUID){
	DB := dao.DB
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println(err)
		return
	}
	stmt, err := tx.Prepare("insert into TokenNum (value,I_id) values (?,?)")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(u1.String(), imageId)
	if err != nil{
		fmt.Println(err)
		return
	}
	tx.Commit()
	fmt.Println(res.LastInsertId())
}
func (this *DAO) Delete(id string , dao *DAO){
	DB := dao.DB
	tx, err := DB.Begin()
	if err != nil{
		fmt.Println(err)
		return
	}
	stmt, err := tx.Prepare("delete  from TokenNum where value = ?")
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	fmt.Println(res.LastInsertId())
}
func (this *DAO) Select(id string ,dao *DAO) IMG{
	//DB=dao.DB
	err :=dao.DB.QueryRow("select ImgPath from IMG where I_Id= (select I_id from TokenNum where value = ? )",id).Scan(&dao.Img.ImgPath)
	if err != nil{
		fmt.Println(err)
		//return
	}
	return dao.Img
}
