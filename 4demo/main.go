// @title 图书管理系统
// @version 1.0
// @description 实现对图书的增删改查的图书管理系统
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// 数据库连接参数
const (
	dbServer   = "localhost"
	dbPort     = 1433
	dbUser     = "sa"
	dbPassword = "123456"
	dbName     = "MUXI"
)

var db *gorm.DB

const SecretKey = "1234567"

// Claims 定义一个结构体用于存储信息
type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

// Book 定义图书信息
type Book struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Stock  string `json:"stock"`
}

// User 用户表
type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"password"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GenerateToken 生成Token
func GenerateToken(userId string) (string, error) {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

// AuthMiddleware 创建一个中间件来验证 JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"message": "未认证"})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
		if err != nil {
			c.JSON(401, gin.H{"message": "Invalid token", "error": err.Error()})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("user_id", claims.UserId)
			c.Next()
			return
		}
		c.JSON(401, gin.H{"message": "Invalid or expired token"})
		c.Abort()
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body User true "用户名和密码"
// @Success 200 {object} Response "注册成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 409 {object} Response "用户已存在"
// @Router /register [post]
func Register(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var user User
	if err := db.First(&user, "username = ?", req.Username).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户已存在"})
		return
	}
	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body User true "用户名和密码"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "用户名或密码错误"
// @Router /login [post]
func Login(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	var user User
	if err := db.First(&user, "username = ?", req.Username).Error; err != nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	token, _ := GenerateToken(req.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AddBook 添加图书
// @Summary 添加书籍
// @Description 传入书籍信息新增书籍
// @Tags 图书
// @Accept json
// @Produce json
// @Param book body Book true "书籍信息"
// @Success 201 {object} map[string]interface{} "添加成功"
// @Failure 400 {object} Response "参数解析失败"
// @Failure 409 {object} Response "书籍ID已存在"
// @Security ApiKeyAuth
// @Router /books [post]
func AddBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exist Book
	if err := db.First(&exist, "id = ?", book.ID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book ID already exists"})
		return
	}
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully", "book": book})
}

// DeleteBook 删除图书
// @Summary 删除书籍
// @Description 根据书籍ID删除书籍
// @Tags 图书
// @Produce json
// @Param id path string true "书籍ID"
// @Success 200 {object} Response "书籍删除成功"
// @Failure 404 {object} Response "书籍不存在"
// @Security ApiKeyAuth
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// UpdateBook 更新图书信息
// @Summary 更新书籍
// @Description 根据ID更新书籍信息
// @Tags 图书
// @Accept json
// @Produce json
// @Param id path string true "书籍ID"
// @Param book body Book true "更新后的书籍信息"
// @Success 200 {object} map[string]interface{} "书籍更新成功"
// @Failure 400 {object} Response "参数绑定失败"
// @Failure 404 {object} Response "书籍不存在"
// @Security ApiKeyAuth
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedBook.ID = id
	if err := db.Model(&book).Updates(updatedBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "book": updatedBook})
}

// SearchAllBook 查询所有图书
// @Summary 获取所有书籍
// @Description 获取所有图书信息
// @Tags 图书
// @Produce json
// @Success 200 {object} map[string][]Book "书籍信息"
// @Security ApiKeyAuth
// @Router /books [get]
func SearchAllBook(c *gin.Context) {
	var bookList []Book
	if err := db.Find(&bookList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	if len(bookList) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No books available", "books": []Book{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": bookList})
}

// GetBook 根据ID查询图书
// @Summary 根据ID查询书籍
// @Description 获取指定ID的图书信息
// @Tags 图书
// @Produce json
// @Param id path string true "书籍ID"
// @Success 200 {object} Book "书籍信息"
// @Failure 404 {object} Response "书籍不存在"
// @Security ApiKeyAuth
// @Router /books/{id} [get]
func GetBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := db.First(&book, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func main() {
	// 构建连接字符串
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		dbUser, dbPassword, dbServer, dbPort, dbName)
	var err error
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移表结构
	_ = db.AutoMigrate(&Book{}, &User{})

	r := gin.Default()
	r.Use(cors.Default())

	// 静态文件服务，假设 index.html 在当前目录
	r.StaticFile("/", "./index.html")

	r.POST("/register", Register)
	r.POST("/login", Login)

	// 图书相关接口需要认证
	book := r.Group("/")
	book.Use(AuthMiddleware())
	{
		book.POST("/books", AddBook)
		book.DELETE("/books/:id", DeleteBook)
		book.PUT("/books/:id", UpdateBook)
		book.GET("/books", SearchAllBook)
		book.GET("/books/:id", GetBook)
	}

	r.Run()
}
