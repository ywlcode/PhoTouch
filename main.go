package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Key struct {
	Email_user string `json:"email_user"`
	Email_key  string `json:"email_key"`
	Mysql_user string `json:"mysql_user"`
	Mysql_key  string `json:"mysql_key"`
	Mysql_db   string `json:"mysql_db"`
	Token_key  string `json:"token"`
}

type bu struct {
	Id     string `json:"id"`
	Imgurl string `json:"url"`
}

type buimg struct {
	Data bu `json:"data"`
}

type uu struct {
	Url_du   string
	Minurl   string
	Year     int
	Month    int
	Day      int
	Userid   int
	Imgbuid  string
	Imgminid string
}

type imgpublic struct {
	Url string `json:"minurl"`
}

var Config Key
var db *sql.DB

func init() { // 初始化config 和 mysql数据库
	keys, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("open config.json failed")
		os.Exit(1)
	}
	err = json.Unmarshal(keys, &Config)
	if err != nil {
		fmt.Println("config init failed")
		os.Exit(1)
	} else {
		fmt.Println("config init succeeded")
	}
	db, err = sql.Open("mysql", Config.Mysql_user+":"+Config.Mysql_key+"@tcp(127.0.0.1:3306)/"+Config.Mysql_db+"?charset=utf8")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		os.Exit(1)
	}
	fmt.Println("Connected database succeeded")
}

func minimg(path string, imagename string) {
	imgData, _ := ioutil.ReadFile(path + imagename)
	buf := bytes.NewBuffer(imgData)
	image, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal(err)
		return
	}
	image = imaging.Fill(image, 400, 400, imaging.Center, imaging.Lanczos)
	err = imaging.Save(image, path+"temp.jpg")
	if err != nil {
		log.Fatal(err)
	}
}

func uploadimg(path string, name string) buimg {
	url := "https://7bu.top/api/upload"
	file, err := os.Open(path + name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", name)
	_, _ = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	httpRequest, _ := http.NewRequest("POST", url, body)
	httpRequest.Header.Add("token", Config.Token_key) // 请求头自定义参数
	httpRequest.Header.Add("Content-Type", writer.FormDataContentType())
	httpclient := &http.Client{} // 创建指向Client 结构体类型的指针
	resp, err := httpclient.Do(httpRequest)
	if err != nil {
		fmt.Println("Failed to post image")
	}
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	var ans buimg
	err = json.Unmarshal(respbody, &ans)
	if err != nil {
		log.Fatal(err)
	}
	return ans
}

func Addimgtodb(alb uu, table string) {
	result, err := db.Exec("INSERT INTO "+table+" (url,minurl,year,month,day,userid,imgbuid,imgminid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", alb.Url_du, alb.Minurl, alb.Year, alb.Month, alb.Day, alb.Userid, alb.Imgbuid, alb.Imgminid)
	if err != nil {
		log.Fatal(err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
}

func PathExists(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func AddimgFromuser(IDuser int) {
	pathtemp := "tempimg" // tempimg文件夹存上传的图片
	file, err := ioutil.ReadDir(pathtemp)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, ff := range file {
		name := ff.Name()
		kk := uploadimg(pathtemp+"/", name)
		timeimg := time.Now()
		var imgdb uu
		imgdb.Url_du = kk.Data.Imgurl
		imgdb.Userid = IDuser
		imgdb.Imgbuid = kk.Data.Id
		imgdb.Year = timeimg.Year()
		imgdb.Month = int(timeimg.Month())
		imgdb.Day = timeimg.Day()

		year := strconv.Itoa(imgdb.Year)
		month := strconv.Itoa(imgdb.Month)
		day := strconv.Itoa(imgdb.Day)
		pathnew := "base/" + year + "/" + month + "/" + day

		PathExists(pathnew)
		minimg(pathtemp+"/", name)
		kk = uploadimg(pathtemp+"/", "temp.jpg")
		imgdb.Minurl = kk.Data.Imgurl
		imgdb.Imgminid = kk.Data.Id
		os.Remove(pathtemp + "/temp.jpg")
		_ = os.Rename(pathtemp+"/"+name, pathnew+"/"+path.Base(kk.Data.Imgurl))
		Addimgtodb(imgdb, "imguser") // 上传图片默认保存在私人相册
	}
}

func randminimg16() []imgpublic {
	var albums []imgpublic
	rows, err := db.Query("SELECT minurl FROM imgpublic ORDER BY RAND() LIMIT 16")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var alb imgpublic
		if err := rows.Scan(&alb.Url); err != nil {
			log.Fatal(err)
			return nil
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return albums
}

func bigimgurl(table string, minurl string) string {
	var alb imgpublic
	rows, err := db.Query("SELECT url FROM "+table+" WHERE minurl = ?", minurl)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&alb.Url)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return alb.Url
}

/* // 打包静态文件 这里导入html模板 不适合开发阶段
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	sum := 0
	for _, name := range bindata.AssetNames() {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		asset, err := bindata.Asset(name)
		if err != nil {
			continue
		}
		name := strings.Replace(name, "assets/templates/", "", 1)
		t, err = t.New(name).Parse(string(asset))
		if err != nil {
			log.Fatal(err)
		}
		sum++
		fmt.Println(sum, ":", name)
	}
	if sum == 4 {
		fmt.Println(time.Now(), "html templates init succeeded ! ! ! ! ! !")
	}
	return t, nil
}
*/

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("assets/templates/*")
	/*  // 打包静态文件 不适合开发阶段
	fs := assetfs.AssetFS{
		Asset:     bindata.Asset,
		AssetDir:  bindata.AssetDir,
		AssetInfo: nil,
		Prefix:    "assets",
	}
	router.StaticFS("/assets", &fs)
	t, err := loadTemplate()
	if err != nil {
		log.Fatal(err)
	}
	router.SetHTMLTemplate(t)
	*/
	PathExists("tempimg")

	router.GET("/", index)
	router.GET("/login", login)
	router.GET("/signup", signup)
	router.GET("/wj", wj)
	router.POST("/user/login")
	//router.POST("/user/login")
	//router.POST("/signup/email")
	//router.POST("/signup/up")
	router.GET("/img/rand", randimgpublic)
	router.POST("/img/big", bigimgpublic)

	router.POST("/upload", uploadimgfromuser)
	router.Run(":8000")
}

func index(c *gin.Context) {
	con := randminimg16()
	data := make(map[string][]imgpublic) // 注意这里只能是 !!! map
	data["imgsrc"] = con                 // 如何随机推荐图片
	c.HTML(http.StatusOK, "index.html", data)
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", "")
}

func signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", "")
}

func wj(c *gin.Context) {
	c.HTML(http.StatusOK, "wj.html", "")
}

func uploadimgfromuser(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	for _, file := range files {
		c.SaveUploadedFile(file, "tempimg/"+file.Filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	ID := 2 // 用户记录ID :: 如何快速获取ID
	AddimgFromuser(ID)
}

func randimgpublic(c *gin.Context) {
	con := randminimg16()
	c.JSON(http.StatusOK, con)
}

func bigimgpublic(c *gin.Context) {
	kkurl := c.PostForm("minurl") // 提取参数
	ans := bigimgurl("imgpublic", kkurl)
	c.String(http.StatusOK, ans)
}
