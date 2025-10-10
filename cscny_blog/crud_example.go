package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CRUDExample 演示如何使用CRUD操作
func CRUDExample() {
	// 数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/crud_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化数据库表
	if err := InitCRUDDatabase(db); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件（可选）
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 设置CRUD路由
	SetupCRUDRoutes(r, db)

	// 添加健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "CRUD API服务运行正常",
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // 使用不同的端口避免与博客系统冲突
	}

	log.Printf("CRUD API服务器启动成功，监听端口: %s", port)
	log.Printf("API文档地址: http://localhost:%s/health", port)
	log.Printf("示例API调用:")
	log.Printf("  GET  http://localhost:%s/api/v1/categories", port)
	log.Printf("  GET  http://localhost:%s/api/v1/products", port)
	log.Printf("  POST http://localhost:%s/api/v1/products", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// ==================== 使用说明 ====================

/*
如何运行CRUD示例：

1. 确保MySQL数据库运行正常
2. 创建数据库：CREATE DATABASE crud_demo;
3. 在main.go中调用CRUDExample()函数，或者创建单独的main函数：

func main() {
    CRUDExample()
}

4. 运行程序：go run *.go

API测试示例（使用curl）：

# 1. 健康检查
curl http://localhost:8081/health

# 2. 获取所有分类
curl http://localhost:8081/api/v1/categories

# 3. 创建新分类
curl -X POST http://localhost:8081/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "新分类"}'

# 4. 获取所有产品
curl http://localhost:8081/api/v1/products

# 5. 创建新产品
curl -X POST http://localhost:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新产品",
    "description": "产品描述",
    "price": 99.99,
    "stock": 10,
    "category_id": 1
  }'

# 6. 获取单个产品
curl http://localhost:8081/api/v1/products/1

# 7. 更新产品
curl -X PUT http://localhost:8081/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新后的产品名",
    "price": 199.99
  }'

# 8. 删除产品
curl -X DELETE http://localhost:8081/api/v1/products/1

CRUD操作特点：

✅ Create（创建）：
- 数据验证
- 外键约束检查
- 错误处理

✅ Read（读取）：
- 单个查询
- 列表查询
- 分页支持
- 关联数据预加载

✅ Update（更新）：
- 部分更新支持
- 数据验证
- 存在性检查

✅ Delete（删除）：
- 软删除（GORM默认）
- 级联删除检查
- 安全性验证

数据库设计特点：

🗄️ 表结构：
- 使用GORM的Model基类（包含ID、CreatedAt、UpdatedAt、DeletedAt）
- 外键关联
- 索引优化
- 数据类型约束

🔗 关联关系：
- 一对多关系（Category -> Products）
- 外键约束
- 预加载支持

🛡️ 安全特性：
- 数据验证
- SQL注入防护
- 错误处理
- 软删除
*/
