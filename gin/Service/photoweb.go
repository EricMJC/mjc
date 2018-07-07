package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"os"
	"io"
	"html/template"
	"net/http"
	"io/ioutil"
	"github.com/satori/go.uuid"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"gin/Dao"
)

type SERVICE struct {
	dao *Dao.DAO
}

var DB *sql.DB
var Img Dao.IMG
var dao *Dao.DAO
var service SERVICE

func Init() {
	dao, _ = Dao.InitDao()
	//service.DB=dao.DB
	//service.Img=dao.Img
	service.dao = dao
}
func (this *SERVICE) upload(c *gin.Context) {

	c.Request.ParseForm()
	if c.Request.Method == "GET" {
		t, err := template.ParseFiles("./view/upload.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		t.Execute(c.Writer, nil)
	} else {
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return
		}
		/*
		if err != nil {
			//fmt.Fprintf(w, " no file")
			fmt.Println("aaaaa: %+v", err)
			return
		}*/
		filename := header.Filename
		defer file.Close()
		f, err := os.Create("./test1" + "/" + filename)
		/*
		if err != nil {
			//fmt.Fprintf(w, " open file failed")
			fmt.Println("aaaaa: %+v", err)
			//http.Redirect(w, c.Request, "/uploadfailed", http.StatusFound)
			return
		}*/
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Println("upload success")
		http.Redirect(c.Writer, c.Request, "/uploadsuccess?id="+filename, http.StatusFound)
	}
}
func (this *SERVICE) uploadsuccess(c *gin.Context) {
	imageId := c.Request.FormValue("id")
	imagePath := "./test1/" + imageId
	t, err := template.ParseFiles("./view/uploadsuccess.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(c.Writer, imagePath)
	http.Redirect(c.Writer, c.Request, "/view?id="+imageId, http.StatusFound)
}
func (this *SERVICE) list1(c *gin.Context) {
	fileInfoArr, err := ioutil.ReadDir("./test1")
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		fmt.Println("failed listhander")
		return
	}
	locals := make(map[string]interface{})
	images := [] string{}
	for _, fileinInfo := range fileInfoArr {

		images = append(images, fileinInfo.Name())

	}
	locals["images"] = images
	t, err := template.ParseFiles("./view/list1.html")
	if err != nil {
		http.Error(c.Writer, err.Error(),
			http.StatusInternalServerError)
		return
	}
	t.Execute(c.Writer, locals)
}

func (this *SERVICE) view(c *gin.Context) {
	imageId := c.Request.FormValue("id")
	fmt.Println(imageId)
	u1, _ := uuid.NewV4()
	dao.Insert(service.dao, imageId, u1)
	//Dao.Insert(DB,imageId,u1)
	http.Redirect(c.Writer, c.Request, "/view1?id="+u1.String(), http.StatusFound)
}

func (this *SERVICE) view1(c *gin.Context) {
	imgId := c.Request.FormValue("id")
	//Img = dao.Select(imgId,DB,Img)
	Img = dao.Select(imgId, service.dao)
	fmt.Println(Img.ImgPath)
	c.Writer.Header().Set("Content-Type", "image")
	http.ServeFile(c.Writer, c.Request, Img.ImgPath)
	//Dao.Delete(imgId,DB)
	dao.Delete(imgId, service.dao)
}

func router(urlPath string, f gin.HandlerFunc) {
	router := gin.Default()
	router.Any(urlPath, f)
	router.Run()
}

func main() {
	router := gin.Default()
	router.Any("/upload", service.upload)
	router.Any("/uploadsuccess", service.uploadsuccess)
	router.Any("/uist1", service.list1)
	router.Any("/view", service.view)
	router.Any("/view1", service.view1)
	router.Run()
}
