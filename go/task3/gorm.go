package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Name      string `gorm:"size:100;not null"`
	Age       uint8
	Email     string `gorm:"size:100;"`
	PostCount uint8
	Posts     []Post
	gorm.Model
}

type Post struct {
	UserID  uint
	Title   string `gorm:"size:255;"`
	Content string `gorm:"type:text;"`
	// 评论状态
	CommentStatus string `gorm:"size:255;"`
	Comments      []Comment
	gorm.Model
}

type Comment struct {
	Content string `gorm:"type:text"`
	PostID  uint
	gorm.Model
}

func (u *Post) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&User{}).Where("id = ?", u.UserID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error
}

func (u *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var count int64
	err = tx.Model(&Comment{}).Where("post_id=?", u.PostID).Count(&count).Error
	if count == 0 && err == nil {
		err = tx.Model(&Post{}).Where("id = ?", u.PostID).Update("comment_status", "无评论").Error
	}
	return err
}

func main() {
	// 配置日志级别为 Info（会打印 SQL）
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值（超过 1s 视为慢查询）
			LogLevel:      logger.Info, // 打印所有 SQL
			Colorful:      true,        // 彩色输出
		},
	)
	// MySQL 连接信息
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	queryData(db)

}

func queryData(db *gorm.DB) {
	// 创建用户
	user := User{Name: "Alice", Email: "alice@example.com"}
	db.Create(&user)

	// 用户发布文章
	post1 := Post{Title: "GORM GuIDe1", Content: "How to use GORM1...", UserID: user.ID}
	post2 := Post{Title: "GORM GuIDe2", Content: "How to use GORM2...", UserID: user.ID}
	db.Create(&post1)
	db.Create(&post2)

	// 文章添加评论
	comment1 := Comment{Content: "Great post1!", PostID: post1.ID}
	comment2 := Comment{Content: "Great post2!", PostID: post2.ID}
	db.Create(&comment1)
	db.Create(&comment2)

	// 查询用户的所有文章（预加载）
	var userWithPosts User
	db.Preload("Posts.Comments").First(&userWithPosts, user.ID)

	// 查询评论数量最多的文章信息
	p := &Post{}
	qurey := db.Model(&Comment{}).Select("post_id, count(id) as num").Group("post_id").Order("num desc").Limit(1)
	db.Joins("join (?) c on c.post_id=posts.id", qurey).First(p)
	log.Println(*p)

	log.Printf("delete post:%v comment\n", post1.ID)
	db.Where("post_id=?", post1.ID).Delete(&Comment{PostID: post1.ID})
}
