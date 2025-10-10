package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// === 数据模型 ===
// Product 产品模型 - 用于演示CRUD操作
type Product struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(100);not null" json:"name" binding:"required"`        // 产品名称
	Description string  `gorm:"type:text" json:"description"`                                     // 产品描述
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price" binding:"required,gt=0"` // 产品价格
	Stock       int     `gorm:"default:0" json:"stock"`                                           // 库存数量
	CategoryID  uint    `gorm:"not null" json:"category_id" binding:"required"`                   // 分类ID
}

// Category 分类模型
type Category struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"name" binding:"required"` // 分类名称
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`                      // 关联产品
}

// === CRUD操作 ===

// ==================== CREATE 创建操作 ====================

// CreateProduct 创建产品
func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product Product
		
		// 绑定JSON数据到结构体
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "数据格式错误",
				"details": err.Error(),
			})
			return
		}

		// 验证分类是否存在
		var category Category
		if err := db.First(&category, product.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "指定的分类不存在",
			})
			return
		}

		// 创建产品
		if err := db.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "创建产品失败",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "产品创建成功",
			"data":    product,
		})
	}
}

// CreateCategory 创建分类
func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var category Category
		
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "数据格式错误",
				"details": err.Error(),
			})
			return
		}

		if err := db.Create(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "创建分类失败",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "分类创建成功",
			"data":    category,
		})
	}
}

// ==================== READ 读取操作 ====================

// GetAllProducts 获取所有产品列表
func GetAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []Product
		
		// 预加载分类信息，支持分页
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		if err := db.Preload("Category").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取产品列表失败",
			})
			return
		}

		// 获取总数
		var total int64
		db.Model(&Product{}).Count(&total)

		c.JSON(http.StatusOK, gin.H{
			"data": products,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	}
}

// GetProductByID 根据ID获取单个产品
func GetProductByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product Product

		// 预加载分类信息
		if err := db.Preload("Category").First(&product, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "产品不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取产品失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": product,
		})
	}
}

// GetAllCategories 获取所有分类
func GetAllCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []Category
		
		// 可选择是否包含产品信息
		includeProducts := c.Query("include_products") == "true"
		
		query := db
		if includeProducts {
			query = query.Preload("Products")
		}
		
		if err := query.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "获取分类列表失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": categories,
		})
	}
}

// ==================== UPDATE 更新操作 ====================

// UpdateProduct 更新产品信息
func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product Product

		// 检查产品是否存在
		if err := db.First(&product, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "产品不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询产品失败",
			})
			return
		}

		// 绑定更新数据
		var updateData Product
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "数据格式错误",
				"details": err.Error(),
			})
			return
		}

		// 如果更新了分类ID，验证分类是否存在
		if updateData.CategoryID != 0 && updateData.CategoryID != product.CategoryID {
			var category Category
			if err := db.First(&category, updateData.CategoryID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "指定的分类不存在",
				})
				return
			}
		}

		// 更新产品
		if err := db.Model(&product).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "更新产品失败",
			})
			return
		}

		// 重新查询更新后的产品信息
		db.Preload("Category").First(&product, id)

		c.JSON(http.StatusOK, gin.H{
			"message": "产品更新成功",
			"data":    product,
		})
	}
}

// ==================== DELETE 删除操作 ====================

// DeleteProduct 删除产品
func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product Product

		// 检查产品是否存在
		if err := db.First(&product, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "产品不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询产品失败",
			})
			return
		}

		// 软删除产品
		if err := db.Delete(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "删除产品失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "产品删除成功",
		})
	}
}

// DeleteCategory 删除分类
func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category Category

		// 检查分类是否存在
		if err := db.First(&category, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "分类不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询分类失败",
			})
			return
		}

		// 检查分类下是否还有产品
		var productCount int64
		db.Model(&Product{}).Where("category_id = ?", id).Count(&productCount)
		if productCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "该分类下还有产品，无法删除",
			})
			return
		}

		// 删除分类
		if err := db.Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "删除分类失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "分类删除成功",
		})
	}
}
