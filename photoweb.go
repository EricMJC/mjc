package main

import (
	"net/http"
	"log"
	"html/template"
	"os"
	"io"
	"fmt"
	"database/sql"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"github.com/satori/go.uuid"
)





//lalalal
//bendi second
//lalalalllllll

var buf []byte
var templates map[string]*template.Template = make(map[string]*template.Template)

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, err := template.ParseFiles("upload.html")
		if err != nil {
			fmt.Fprintf(w, "upload.html failed")
			return
		}
		//checkErr(err)
		t.Execute(w, nil)
	} else {
		file, handle, err := r.FormFile("image")
		if err != nil {
			fmt.Fprintf(w, " no file")
			fmt.Println("aaaaa: %+v", err)
			return
		}
		filename := handle.Filename
		defer file.Close()
		//checkErr(err)
		//f, err := os.OpenFile("./test1"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		f, err := os.Create("./test1" + "/" + filename)
		if err != nil {
			//fmt.Fprintf(w, " open file failed")
			fmt.Println("aaaaa: %+v", err)
			http.Redirect(w, r, "/uploadfailed", http.StatusFound)
			return
		}

		defer f.Close()
		io.Copy(f, file)

		//checkErr(err)

		fmt.Println("upload success")
		http.Redirect(w, r, "/uploadsuccess?id="+handle.Filename, http.StatusFound)
	}
}

func uploadsuccess(w http.ResponseWriter, r *http.Request) {
	/*
 imageId := r.FormValue("id")
 imagePath := "./test1"+"/"+imageId

 w.Header().Set("Content-Type","image")
 http.ServeFile(w,r,imagePath)
	*/

	imageId := r.FormValue("id")
	imagePath := "./test1/" + imageId
	/*
	_, err := os.Stat(imagePath)
	if (os.IsNotExist(err)) {
		fmt.Println("file  not exits")
		http.Redirect(w, r, "/uploadfailed", http.StatusFound)
	}
	*/
	t, _ := template.ParseFiles("uploadsuccess.html")
	t.Execute(w, imagePath)
	http.Redirect(w, r, "/view?id="+imageId, http.StatusFound)

}
func uploadfailed(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("uploadfailed.html")
	t.Execute(w, nil)
}
func list1(w http.ResponseWriter, r *http.Request) {
	//t, _ := template.ParseFiles("list1.html")
	//t.Execute(w, nil)
	//文件夹的数组
	/*
	var err1 error
	funcMap := template.FuncMap{"generrateRand":generrateRand}
	t1 := template.New("list1.html").Funcs(funcMap)
	t1,err1= t1.ParseFiles("./list1.html")
	if err1!= nil{
		fmt.Println(err1)
	}
	t1.Execute(w,nil)*/
	fileInfoArr, err := ioutil.ReadDir("./test1")
	if err != nil {
		http.Error(w, err.Error(),
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
	//readerHtml(w, "list", locals);
	t, err := template.ParseFiles("list1.html")
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	t.Execute(w, locals)

}

func readerHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	err := templates[tmpl].Execute(w, locals)
	check(err)
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

var v1 string = "s"

func view(w http.ResponseWriter, r *http.Request) {
	//t, _ := template.ParseFiles("view.html")
	//t.Execute(w, nil)
	//var Img IMG
	//fmt.Println(";lalala")

	imageId := r.FormValue("id")
	fmt.Println(imageId)
	u1, _ := uuid.NewV4()
	//fmt.Println(u1)
	tx, _ := DB.Begin()
	stmt, err2 := tx.Prepare("insert into TokenNum (value,I_id) values (?,?)")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	res, err1 := stmt.Exec(u1, imageId)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	tx.Commit()
	fmt.Println(res.LastInsertId())

	http.Redirect(w, r, "/view1?id="+u1.String(), http.StatusFound)
	//imagePath := "./test1/" + imageId
	/*
	err := DB.QueryRow("select ImgPath from IMG where I_Id= (select I_Id from TokenNum where T_value = ?)", u1).Scan(&Img.ImgPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		fmt.Println("404")
		return
	}//

	fmt.Println(Img.ImgPath)
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, Img.ImgPath)
	*/
}
func view1(w http.ResponseWriter, r *http.Request) {
	var Img IMG
	imgId := r.FormValue("id")
	//Index:=r.FormValue("num")

	err := DB.QueryRow("select ImgPath from IMG where I_Id= (select I_id from TokenNum where value = ?)", imgId).Scan(&Img.ImgPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
if exists := isExists(imagePath); !exists {
	http.NotFound(w, r)
	fmt.Println("404")
	return
}
*/
	fmt.Println(Img.ImgPath)
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, Img.ImgPath)
	tx, err := DB.Begin()

	if err != nil {
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
	res, err := stmt.Exec(imgId)
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
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}
	return os.IsExist(err)
}
//bendicaozuo


/*
func checkErr(err error){
	if err!= nil {
		err.Error()
	}
}
*/
type IMG struct {
	I_id    int
	ImgPath string
}

const (
	userName = "root"
	password = "173167479"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "mysqlTmd"
)

var DB *sql.DB

func init() {

	var err error
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, err = sql.Open("mysql", path)
	if err != nil {
		fmt.Println(err)
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database failed")
		fmt.Println(err)
		return
	}
	fmt.Println("connect success")

}
func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/uploadsuccess", uploadsuccess)
	http.HandleFunc("/uploadfailed", uploadfailed)
	http.HandleFunc("/list1", list1)
	http.HandleFunc("/view", view)
	http.HandleFunc("/view1", view1)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("listenandserve:", err)
	}

}
