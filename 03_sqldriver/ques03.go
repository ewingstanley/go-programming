package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 使用一致的命名规范
type UserModel struct {
	gorm.Model
	Name     string
	Posts    []ArticleModel `gorm:"foreignKey:UserID"` // Articles → Posts
	Comments []CommentModel `gorm:"foreignKey:UserID"`
}

type ArticleModel struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     UserModel      `gorm:"foreignKey:UserID"`
	Comments []CommentModel `gorm:"foreignKey:PostID"`
}

type CommentModel struct {
	gorm.Model
	Comment string
	UserID  uint
	PostID  uint
	User    UserModel    `gorm:"foreignKey:UserID"`
	Post    ArticleModel `gorm:"foreignKey:PostID"` // Articles → Post
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/dbtest?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&UserModel{}, &ArticleModel{}, &CommentModel{})
	if err != nil {
		log.Fatal("迁移失败：", err)
	}

	// 1. 查询用户所有文章及评论
	fmt.Println("=== 查询用户所有文章及评论 ===")
	err = getUserPostsWithComments(db, 1)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}

	// 2. 查询评论数量最多的文章
	fmt.Println("\n=== 查询评论数量最多的文章 ===")
	mostCommentedPost, err := getMostCommentedPost(db)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("评论最多的文章: ID=%d, 标题=%s, 评论数=%d\n",
			mostCommentedPost.ID, mostCommentedPost.Title, len(mostCommentedPost.Comments))
	}
}

func getUserPostsWithComments(db *gorm.DB, userID uint) error {
	var user UserModel

	result := db.
		Preload("Posts.Comments"). // Articles → Posts
		Preload("Posts.Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(&user, userID)

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("用户 %s 的文章列表:\n", user.Name)
	for _, post := range user.Posts { // Articles → Posts
		fmt.Printf("  文章: %s (ID: %d)\n", post.Title, post.ID)
		fmt.Printf("    评论数: %d\n", len(post.Comments))
		for _, comment := range post.Comments {
			fmt.Printf("    - %s (评论者: %s)\n", comment.Comment, comment.User.Name)
		}
		fmt.Println()
	}

	return nil
}

func getMostCommentedPost(db *gorm.DB) (ArticleModel, error) {
	var post ArticleModel
	var postID uint

	// 使用更简洁的查询方式
	err := db.Model(&CommentModel{}).
		Select("post_id, COUNT(*) as comment_count").
		Group("post_id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&postID).Error

	if err != nil {
		return ArticleModel{}, err
	}

	// 加载文章详情
	err = db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("Comments").
		Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(&post, postID).Error

	return post, err
}