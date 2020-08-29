package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Tweetモデルの宣言
// モデルとはDBのテーブル構造をGoの構造体で表したもの
type Tweet struct {
	gorm.Model
	Content string
}

// DBの初期化
func dbInit() {
	// MySQLだと文字コードの問題で"?parseTime=true"を末尾につける必要がある
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		panic("DBが開けません(dbInit())")
	}
	// コネクション解放
	defer db.Close()
	db.AutoMigrate(&Tweet{}) // 構造体に基づいてテーブルを作成
}

// データインサート処理
func dbInsert(content string) {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		panic("You can't open DB(dbInsert())")
	}
	defer db.Close()
	// Insert処理
	db.Create(&Tweet{Content: content})
}

// 全件処理
func dbGetAll() []Tweet {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		panic("You can't open DB(dbGetAll())")
	}
	defer db.Close()
	var tweets []Tweet
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&tweets)

	return tweets
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("static/*html")

	// インデックスページのルーティング
	router.GET("/", func(c *gin.Context) {
		tweets := dbGetAll()
		c.HTML(200, "index.html", gin.H{"tweets": tweets})
	})

	// POSTデータを受け取ってDBに登録する
	router.POST("/new", func(c *gin.Context) {
		content := c.PostForm("content")
		dbInsert(content)
		c.Redirect(302, "/")
	})

	router.Run()
}
