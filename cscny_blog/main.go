package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// === 全局配置 ===
var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))                             // JWT密钥（环境变量获取）
	logger    = log.New(os.Stdout, "[Blog] ", log.LstdFlags|log.Lshortfile) // 日志实例
)

// === 数据模型 ===
// User 用户模型
type User struct {
	gorm.Model
	Username string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 用户名（唯一）
	Password string    `gorm:"not null" json:"-"`                                     // 密码哈希（json:-表示不返回）
	Posts    []Post    `gorm:"foreignKey:UserID" json:"-"`                            // 关联文章
	Comments []Comment `gorm:"foreignKey:UserID" json:"-"`                            // 关联评论
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"type:varchar(100);not null" json:"title"` // 标题
	Content  string    `gorm:"type:text;not null" json:"content"`       // 内容
	UserID   uint      `gorm:"not null" json:"user_id"`                 // 作者ID（外键）
	User     User      `gorm:"foreignKey:UserID" json:"author"`         // 作者信息（关联用户）
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments"`       // 关联评论
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"` // 评论内容
	UserID  uint   `gorm:"not null" json:"user_id"`           // 评论者ID（外键）
	PostID  uint   `gorm:"not null" json:"post_id"`           // 文章ID（外键）
	User    User   `gorm:"foreignKey:UserID" json:"author"`   // 评论者信息（关联用户）
	Post    Post   `gorm:"foreignKey:PostID" json:"-"`        // 关联文章（不返回）
}

// === 认证中间件（已实现，直接复用） ===
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Authorization格式（需为Bearer <token>）"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("无效的签名算法")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			logger.Printf("JWT验证失败: %v", err) // 记录错误日志
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token（已过期或被伪造）"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token载荷"})
			c.Abort()
			return
		}

		userId := uint(claims["sub"].(float64))
		c.Set("userId", userId)
		c.Next()
	}
}

// === 文章管理功能 ===
// 创建文章（需认证）
func createPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("userId") // 从上下文获取当前用户ID

		var input struct {
			Title   string `json:"title" binding:"required,min=1,max=100"` // 标题必填，1-100字
			Content string `json:"content" binding:"required,min=10"`      // 内容必填，至少10字
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post := Post{
			Title:   input.Title,
			Content: input.Content,
			UserID:  userId.(uint), // 关联当前用户为作者
		}
		if err := db.Create(&post).Error; err != nil {
			logger.Printf("创建文章失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
			return
		}

		// 关联查询作者信息（返回给客户端）
		db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username") // 只返回作者ID和用户名，避免敏感信息
		}).First(&post)

		c.JSON(http.StatusCreated, gin.H{"data": post})
	}
}

// 获取所有文章列表（无需认证）
func listPostsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var posts []Post
		// 预加载作者信息（只返回ID和用户名）
		if err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username")
		}).Find(&posts).Error; err != nil {
			logger.Printf("查询文章列表失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": posts})
	}
}

// 获取单篇文章详情（无需认证）
func getPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id") // 从URL参数获取文章ID

		var post Post
		// 预加载作者信息和评论列表（评论需预加载评论者信息）
		if err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username")
		}).Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username")
		}).First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
				return
			}
			logger.Printf("查询文章详情失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": post})
	}
}

// 更新文章（仅作者可操作）
func updatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("userId")
		postID := c.Param("id")

		var post Post
		if err := db.First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
				return
			}
			logger.Printf("查询文章失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询文章失败"})
			return
		}

		// 验证当前用户是否为作者
		if post.UserID != userId.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此文章"})
			return
		}

		var input struct {
			Title   string `json:"title" binding:"omitempty,min=1,max=100"` // 可选更新，1-100字
			Content string `json:"content" binding:"omitempty,min=10"`      // 可选更新，至少10字
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 只更新非空字段
		if input.Title != "" {
			post.Title = input.Title
		}
		if input.Content != "" {
			post.Content = input.Content
		}

		if err := db.Save(&post).Error; err != nil {
			logger.Printf("更新文章失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
			return
		}

		// 关联作者信息返回
		db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username")
		}).First(&post, postID)

		c.JSON(http.StatusOK, gin.H{"data": post})
	}
}

// 删除文章（仅作者可操作）
func deletePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("userId")
		postID := c.Param("id")

		var post Post
		if err := db.First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
				return
			}
			logger.Printf("查询文章失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}

		// 验证权限
		if post.UserID != userId.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
			return
		}

		// 级联删除评论（或在数据库设置外键级联删除）
		if err := db.Where("post_id = ?", postID).Delete(&Comment{}).Error; err != nil {
			logger.Printf("删除评论失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}

		// 删除文章
		if err := db.Delete(&post).Error; err != nil {
			logger.Printf("删除文章失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
	}
}

// === 评论功能 ===
// 创建评论（需认证）
func createCommentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("userId")
		postID := c.Param("post_id") // 从URL参数获取文章ID

		// 验证文章是否存在
		var post Post
		if err := db.First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
				return
			}
			logger.Printf("查询文章失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
			return
		}

		var input struct {
			Content string `json:"content" binding:"required,min=1,max=500"` // 评论内容，1-500字
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		comment := Comment{
			Content: input.Content,
			UserID:  userId.(uint),
			PostID:  post.ID,
		}
		if err := db.Create(&comment).Error; err != nil {
			logger.Printf("创建评论失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
			return
		}

		// 关联评论者信息
		db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username")
		}).First(&comment)

		c.JSON(http.StatusCreated, gin.H{"data": comment})
	}
}

// 获取文章评论列表（无需认证）
func listCommentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("post_id")

		var comments []Comment
		if err := db.Where("post_id = ?", postID).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Username")
		}).Find(&comments).Error; err != nil {
			logger.Printf("查询评论列表失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询评论失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": comments})
	}
}

// === 错误处理与日志记录 ===
// 全局错误处理中间件（记录所有请求错误）
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先执行后续Handler

		// 检查是否有错误发生
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Printf("请求错误: %s %s - %v", c.Request.Method, c.Request.URL.Path, err.Error())
			// 默认返回500错误，除非Handler已设置状态码
			if c.Writer.Status() == http.StatusOK {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			}
		}
	}
}

// === 路由设置 ===
func setupRoutes(r *gin.Engine, db *gorm.DB) {
	r.Use(errorHandler()) // 全局错误处理中间件

	// 公开路由
	public := r.Group("/api/public")
	{
		// 认证相关
		public.POST("/auth/register", registerHandler(db))
		public.POST("/auth/login", loginHandler(db))
		// 文章相关（无需认证）
		public.GET("/posts", listPostsHandler(db))                      // 所有文章列表
		public.GET("/posts/:id", getPostHandler(db))                    // 单篇文章详情
		public.GET("/posts/:post_id/comments", listCommentsHandler(db)) // 文章评论列表
	}

	// 保护路由（需认证）
	protected := r.Group("/api/protected")
	protected.Use(authMiddleware())
	{
		// 文章相关
		protected.POST("/posts", createPostHandler(db))       // 创建文章
		protected.PUT("/posts/:id", updatePostHandler(db))    // 更新文章
		protected.DELETE("/posts/:id", deletePostHandler(db)) // 删除文章
		// 评论相关
		protected.POST("/posts/:post_id/comments", createCommentHandler(db)) // 创建评论
	}
}

// === 原有注册/登录Handler（复用并优化） ===
func registerHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required,min=3,max=20"` // 用户名3-20字
			Password string `json:"password" binding:"required,min=6,max=32"` // 密码6-32字
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查用户名是否已存在
		var existingUser User
		if err := db.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}

		hashpassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Printf("密码加密失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}

		user := User{
			Username: input.Username,
			Password: string(hashpassword),
		}
		if err := db.Create(&user).Error; err != nil {
			logger.Printf("用户创建失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
	}
}

func loginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		claims := jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			logger.Printf("JWT生成失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败，请重试"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"token":   tokenString,
		})
	}
}

// === 主函数 ===
func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/dbtest?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移表结构（添加外键约束）
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		logger.Fatalf("自动迁移失败: %v", err)
	}
	logger.Println("数据库表迁移成功")

	r := gin.Default()
	setupRoutes(r, db)

	logger.Println("服务器启动成功，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		logger.Fatalf("服务器启动失败: %v", err)
	}
}