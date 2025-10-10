package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupCRUDRoutes 设置CRUD操作的路由
func SetupCRUDRoutes(r *gin.Engine, db *gorm.DB) {
	// API版本分组
	api := r.Group("/api/v1")
	
	// ==================== 产品相关路由 ====================
	products := api.Group("/products")
	{
		// CREATE - 创建产品
		products.POST("", CreateProduct(db))
		
		// READ - 读取操作
		products.GET("", GetAllProducts(db))        // 获取所有产品（支持分页）
		products.GET("/:id", GetProductByID(db))    // 根据ID获取单个产品
		
		// UPDATE - 更新产品
		products.PUT("/:id", UpdateProduct(db))     // 更新产品信息
		
		// DELETE - 删除产品
		products.DELETE("/:id", DeleteProduct(db))  // 删除产品
	}
	
	// ==================== 分类相关路由 ====================
	categories := api.Group("/categories")
	{
		// CREATE - 创建分类
		categories.POST("", CreateCategory(db))
		
		// READ - 读取操作
		categories.GET("", GetAllCategories(db))    // 获取所有分类
		
		// DELETE - 删除分类
		categories.DELETE("/:id", DeleteCategory(db)) // 删除分类
	}
}

// ==================== 数据库初始化 ====================

// InitCRUDDatabase 初始化CRUD操作相关的数据库表
func InitCRUDDatabase(db *gorm.DB) error {
	// 自动迁移表结构
	err := db.AutoMigrate(&Category{}, &Product{})
	if err != nil {
		return err
	}
	
	// 创建一些示例数据（可选）
	CreateSampleData(db)
	
	return nil
}

// CreateSampleData 创建示例数据
func CreateSampleData(db *gorm.DB) {
	// 检查是否已有数据
	var categoryCount int64
	db.Model(&Category{}).Count(&categoryCount)
	if categoryCount > 0 {
		return // 已有数据，不重复创建
	}
	
	// 创建示例分类
	categories := []Category{
		{Name: "电子产品"},
		{Name: "服装"},
		{Name: "图书"},
		{Name: "家居用品"},
	}
	
	for _, category := range categories {
		db.Create(&category)
	}
	
	// 创建示例产品
	products := []Product{
		{
			Name:        "iPhone 15",
			Description: "苹果最新款智能手机",
			Price:       7999.00,
			Stock:       50,
			CategoryID:  1, // 电子产品
		},
		{
			Name:        "MacBook Pro",
			Description: "苹果专业级笔记本电脑",
			Price:       15999.00,
			Stock:       20,
			CategoryID:  1, // 电子产品
		},
		{
			Name:        "Nike运动鞋",
			Description: "舒适透气的运动鞋",
			Price:       599.00,
			Stock:       100,
			CategoryID:  2, // 服装
		},
		{
			Name:        "Go语言编程",
			Description: "Go语言学习指南",
			Price:       89.00,
			Stock:       200,
			CategoryID:  3, // 图书
		},
	}
	
	for _, product := range products {
		db.Create(&product)
	}
}

// ==================== API使用示例 ====================

/*
CRUD操作API使用示例：

1. 创建分类
POST /api/v1/categories
{
    "name": "电子产品"
}

2. 获取所有分类
GET /api/v1/categories

3. 创建产品
POST /api/v1/products
{
    "name": "iPhone 15",
    "description": "苹果最新款智能手机",
    "price": 7999.00,
    "stock": 50,
    "category_id": 1
}

4. 获取所有产品（支持分页）
GET /api/v1/products?page=1&limit=10

5. 获取单个产品
GET /api/v1/products/1

6. 更新产品
PUT /api/v1/products/1
{
    "name": "iPhone 15 Pro",
    "price": 8999.00,
    "stock": 30
}

7. 删除产品
DELETE /api/v1/products/1

8. 删除分类
DELETE /api/v1/categories/1

响应格式示例：

成功响应：
{
    "message": "操作成功",
    "data": { ... }
}

错误响应：
{
    "error": "错误信息",
    "details": "详细错误信息"
}

分页响应：
{
    "data": [ ... ],
    "pagination": {
        "page": 1,
        "limit": 10,
        "total": 100
    }
}
*/
