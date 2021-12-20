package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path"
	"photouch/bindata"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jordan-wright/email"
)

// 运行密钥,从config.json中读取
type Key struct {
	Email_user string `json:"email_user"`
	Email_key  string `json:"email_key"`
	Mysql_user string `json:"mysql_user"`
	Mysql_key  string `json:"mysql_key"`
	Mysql_db   string `json:"mysql_db"`
	Token_key  string `json:"token"`
}

// 7bu 图床上传返回信息读取结构体
type bu struct {
	Id     string `json:"id"`
	Imgurl string `json:"url"`
}

type buimg struct {
	Data bu `json:"data"`
}

// 数据库图片信息表数据行对应结构体
type imgrow struct {
	Url_du   string
	Minurl   string
	Year     int
	Month    int
	Day      int
	Userid   int
	Imgbuid  string
	Imgminid string
}

// 图片 一般信息结构体
type imginformation struct {
	Url string `json:"minurl"`
}

type userdb struct {
	userid int
	email  string
	pwd    string
}

// 点赞结构体
type gooddb struct {
	userid int
	minurl string
}

// 数据库句柄
var db *sql.DB
var Config Key

// 此时随机值
var todayrand string

// mapcookie结构
type usermanage struct {
	userid int
	last   time.Time
}

// 管理用户登录
var live map[string]usermanage

// 注册验证码
var code map[string]string

// 初始化config , mysql数据库句柄 和 随机值
func init() {
	fmt.Println("初始化")
	fmt.Println("----------------------------------")
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
	live = make(map[string]usermanage)
	code = make(map[string]string)
	todayrand = strconv.Itoa(rand.Intn(100000))
	fmt.Println("----------------------------------")
}

//将指定图片压缩生成缩略图保存到"tempimg/temp.jpg"
func minimg(path string, imagename string) {
	imgData, _ := ioutil.ReadFile(path + imagename)
	buf := bytes.NewBuffer(imgData)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	image = imaging.Fill(image, 400, 400, imaging.Center, imaging.Lanczos)
	err = imaging.Save(image, path+"temp.jpg")
	if err != nil {
		fmt.Println(err)
	}
}

//上传指定图片到7bu图床
func uploadimg(path string, name string) buimg {
	url := "https://7bu.top/api/upload"
	file, err := os.Open(path + name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", name)
	_, _ = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}
	return ans
}

//添加新图片到指定数据表
func Addimgtodb(alb imgrow, table string) {
	result, err := db.Exec("INSERT INTO "+table+" (url,minurl,year,month,day,userid,imgbuid,imgminid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", alb.Url_du, alb.Minurl, alb.Year, alb.Month, alb.Day, alb.Userid, alb.Imgbuid, alb.Imgminid)
	if err != nil {
		fmt.Println(err)
	}
	_, err = result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
}

//查看目录是否创建,未创建则创建
func PathExists(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}

// 处理用户上传的图片,发送到7bu图床
func AddimgFromuser(IDuser int) {
	pathtemp := "tempimg" // tempimg文件夹存上传的图片
	file, err := ioutil.ReadDir(pathtemp)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ff := range file {
		name := ff.Name()
		kk := uploadimg(pathtemp+"/", name)
		timeimg := time.Now()
		var imgdb imgrow
		imgdb.Url_du = kk.Data.Imgurl
		imgdb.Userid = IDuser
		imgdb.Imgbuid = kk.Data.Id
		imgdb.Year = timeimg.Year()
		imgdb.Month = int(timeimg.Month())
		imgdb.Day = timeimg.Day()
		bigurl := kk.Data.Imgurl
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
		_ = os.Rename(pathtemp+"/"+name, pathnew+"/"+path.Base(bigurl))
		Addimgtodb(imgdb, "imguser") // 上传图片默认保存在私人相册
	}
}

// 检索随机32张分享图片的minurl,返回imgpublic
func randminimg32() []imginformation {
	var albums []imginformation
	rows, err := db.Query("SELECT minurl FROM imgpublic ORDER BY RAND() LIMIT 32")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var alb imginformation
		if err := rows.Scan(&alb.Url); err != nil {
			fmt.Println(err)
			return nil
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil
	}
	return albums
}

// 通过minurl和数据表名检索对应大图url并返回
func bigimgurl(table string, minurl string) string {
	var alb imginformation
	rows, err := db.Query("SELECT url FROM "+table+" WHERE minurl = ?", minurl)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&alb.Url)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return alb.Url
}

// 打包静态文件 这里导入html模板 不适合开发阶段
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

// 通过minurl 改变私有或私有状态
func changeimgquan(minurl string, oldtable string, newtable string) {
	var alb imgrow
	rows, err := db.Query("SELECT * FROM "+oldtable+" WHERE minurl = ?", minurl)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var id int
	for rows.Next() {
		if err := rows.Scan(&id, &alb.Url_du, &alb.Minurl, &alb.Year, &alb.Month, &alb.Day, &alb.Userid, &alb.Imgbuid, &alb.Imgminid); err != nil {
			fmt.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	Addimgtodb(alb, newtable)
	_, _ = db.Exec("delete from "+oldtable+" WHERE minurl = ?", minurl)
}

/*
// 检索出用户的某年某月的照片的minurl列表, 年, 月, 用户id 返回 imgpublic切片
func SelectByMonth(year string, month string, userid int) []imginformation {
	var albums []imginformation
	rows, err := db.Query("SELECT minurl FROM imguser WHERE year = ? AND month = ? AND userid = ?", year, month, userid)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb imginformation
		if err := rows.Scan(&alb.Url); err != nil {
			fmt.Println(err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return albums
}
*/

// 检索出用户所有分享的照片的minurl列表, 按年月日排序,返回 imgpublic切片
func Selectuserpublic(userid int) []imginformation {
	var albums []imginformation
	rows, err := db.Query("SELECT minurl FROM imgpublic WHERE userid = ? ORDER BY year DESC,month DESC,day DESC", userid)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb imginformation
		if err := rows.Scan(&alb.Url); err != nil {
			fmt.Println(err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return albums
}

// 检索出用户所有私有的照片的minurl列表, 按年月日排序,返回 imgpublic切片
func Selectuser(userid int) []imginformation {
	var albums []imginformation
	rows, err := db.Query("SELECT minurl FROM imguser WHERE userid = ? ORDER BY year DESC,month DESC,day DESC", userid)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb imginformation
		if err := rows.Scan(&alb.Url); err != nil {
			fmt.Println(err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return albums
}

// 检索出用户点赞的照片的minurl列表,返回 imgpublic切片
func Selectuserlike(userid int) []imginformation {
	var albums []imginformation
	rows, err := db.Query("SELECT minurl FROM good WHERE userid = ?", userid)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb imginformation
		if err := rows.Scan(&alb.Url); err != nil {
			fmt.Println(err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	return albums
}

// base64编码
func encodebase64(data string) string {
	// Base64 Standard Encoding
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}

// 自定义中间件 验证是否登录
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Cookie("user")
		tt, ok := live[val]
		if ok {
			now := time.Now()
			if now.Before(tt.last) {
				c.Set("Auth", "YES")
				c.Set("id", tt.userid)
			} else {
				c.Set("Auth", "NO")
				c.Set("id", 0)
			}
		} else {
			c.Set("Auth", "NO")
			c.Set("id", 0)
		}
		c.Next()
	}
}

// 发送邮件
func SendMail(toemail string, emailcode string) {
	e := email.NewEmail()
	e.From = "PhoTouch APP Registration verification code" + "<1589292300@qq.com>"
	e.To = []string{toemail}
	e.Subject = "Registration verification code"
	e.HTML = []byte("<h1>you code is " + emailcode + "</h1>")
	auth := smtp.PlainAuth("", "1589292300@qq.com", "kwtrsuisbscdbaac", "smtp.qq.com")
	err := e.Send("smtp.qq.com:25", auth)
	if err != nil {
		fmt.Println(err)
	}
}

// 主程序
func main() {
	go timeToTime() // 定时器
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(Auth())
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

	// 四个页面
	router.GET("/", index)
	router.GET("/login", login)
	router.GET("/signup", signup)
	router.GET("/wj", wj)

	//不需要登录认证的
	router.POST("/user/login", loginkk)
	router.POST("/signup/email", signupsend)
	router.POST("/signup/up", signupup)
	router.GET("/download/:a/:b/:c/:name", downloadimg)
	router.GET("/img/rand", randimgpublic)
	router.POST("/img/big", bigimg)

	// 需要登录认证的
	router.GET("/img/user/like", userlike)
	router.GET("/img/user/all", userallimg)
	router.GET("/img/user/share", usershare)
	//router.GET("/img/user/time")
	//router.GET("/img/user/month")
	router.POST("/upload", uploadimgfromuser)
	router.POST("/change", changeimg)
	router.POST("/good", goodgood)
	router.Run(":8000")
}

// 主页
func index(c *gin.Context) {
	con := randminimg32()
	data := make(map[string][]imginformation) // 注意这里只能是 !!! map
	data["imgsrc"] = con                      // 如何随机推荐图片
	c.HTML(http.StatusOK, "index.html", data)
}

// 登录
func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", "")
}

// 注册
func signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", "")
}

// 忘记密码
func wj(c *gin.Context) {
	c.HTML(http.StatusOK, "wj.html", "")
}

// 登录验证
func loginkk(c *gin.Context) {
	email := c.PostForm("name")
	pwd := c.PostForm("pwd")
	var albums userdb
	rows, err := db.Query("SELECT * FROM user WHERE email = ?", email)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&albums.userid, &albums.email, &albums.pwd); err != nil {
			fmt.Println(err)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	if pwd == albums.pwd {
		now := time.Now()
		timechu := strconv.Itoa(int(now.Unix()))
		cookiejar := strconv.Itoa(albums.userid) + todayrand + timechu
		cookie := encodebase64(cookiejar)
		c.SetCookie("user", cookie, 1800, "/", "", false, false)
		//参数 1.key 2.对应的值 3.过期时间,单位秒 4.cookie 所在的目录
		// 5.cookie 作用范围 5.是否只能通过 https 访问 6.是否对js隐藏(js不能操作)
		ma := usermanage{userid: albums.userid, last: now.Add(30 * time.Minute)}
		live[cookie] = ma
		c.JSON(200, gin.H{"ss": "200"})
	} else {
		c.JSON(201, gin.H{"ss": "400"})
	}
}

// 用户上传图片
func uploadimgfromuser(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	for _, file := range files {
		c.SaveUploadedFile(file, "tempimg/"+file.Filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	fmt.Printf("%d files uploaded!", len(files))
	ID := c.MustGet("id").(int)
	AddimgFromuser(ID)
}

// 热门随机32张图片
func randimgpublic(c *gin.Context) {
	con := randminimg32()
	c.JSON(http.StatusOK, con)
}

// 获取大图url
func bigimg(c *gin.Context) {
	kkurl := c.PostForm("minurl") // 提取参数
	ans := bigimgurl("imgpublic", kkurl)
	if ans == "" {
		ans = bigimgurl("imguser", kkurl)
	}
	c.String(http.StatusOK, ans)
}

// 改变公有私有权限, 分享或取消分享
func changeimg(c *gin.Context) {
	minurl := c.PostForm("minurl") // 提取参数
	old := c.PostForm("old")
	new := c.PostForm("new")
	changeimgquan(minurl, old, new)
	c.String(http.StatusOK, "YES")
}

// 用户点赞
func goodgood(c *gin.Context) {
	minurl := c.PostForm("minurl")
	ID := c.MustGet("id").(int)
	var albums []gooddb
	rows, err := db.Query("SELECT * FROM good WHERE userid = ? AND minurl = ?", ID, minurl)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb gooddb
		if err := rows.Scan(&alb.userid, &alb.minurl); err != nil {
			fmt.Println(err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}
	if len(albums) >= 1 {
		_, err = db.Exec("delete from good WHERE minurl = ?", minurl)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = db.Exec("INSERT INTO good (userid,minurl) VALUES (?, ?)", ID, minurl)
		if err != nil {
			fmt.Println(err)
		}
	}
	c.String(http.StatusOK, "ok")
}

// 发送验证码
func signupsend(c *gin.Context) {
	emailname := c.PostForm("email")
	randcode := strconv.Itoa(rand.Intn(100000))
	code[randcode] = emailname
	SendMail(emailname, randcode)
	c.String(http.StatusOK, "YES")
}

// 检验验证码正确性 添加用户到DB
func signupup(c *gin.Context) {
	emailcode := c.PostForm("code")
	emailname := c.PostForm("emailname")
	password := c.PostForm("pwd")
	kk, ok := code[emailcode]
	if ok {
		if kk == emailname {
			result, err := db.Exec("INSERT INTO user (email,pwd) VALUES (?, ?)", emailname, password)
			if err != nil {
				fmt.Println(err)
			}
			_, err = result.LastInsertId()
			if err != nil {
				fmt.Println(err)
			}
			delete(code, emailcode)
			c.String(http.StatusOK, "YES")
		} else {
			c.String(http.StatusOK, "NO")
		}
	} else {
		c.String(http.StatusOK, "NO")
	}
}

//	定时器
func timeToTime() {
	ticker := time.Tick(10 * time.Minute) //定义一个10分钟间隔的定时器
	for tt := range ticker {
		for iid, n := range live {
			if n.last.Before(time.Now()) {
				delete(live, iid)
				fmt.Println(tt, " user "+iid+" out")
			}
		}
	}
}

// 用户所有私有图片
func userallimg(c *gin.Context) {
	ID := c.MustGet("id").(int)
	data := Selectuser(ID)
	c.JSON(http.StatusOK, data)
}

// 用户所有分享的图片
func usershare(c *gin.Context) {
	ID := c.MustGet("id").(int)
	data := Selectuserpublic(ID)
	c.JSON(http.StatusOK, data)
}

// 用户所有点赞的图片
func userlike(c *gin.Context) {
	ID := c.MustGet("id").(int)
	data := Selectuserlike(ID)
	c.JSON(http.StatusOK, data)
}

// 用户下载大图
func downloadimg(c *gin.Context) {
	year := c.Param("a")
	month := c.Param("b")
	day := c.Param("c")
	name := c.Param("name")
	path := year + "/" + month + "/" + day + "/" + name
	c.Header("content-disposition", `attachment; filename=`+name)
	imgData, err := ioutil.ReadFile("base/" + path)
	if err != nil {
		fmt.Println(err)
	}
	ContentType := http.DetectContentType(imgData)
	c.Data(200, ContentType, imgData)
}
