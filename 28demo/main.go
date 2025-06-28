package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 全局变量存储图书数据
var books = make(map[string]Book)

// 简单用户存储（用户名:密码）
var users = map[string]string{}

const SecretKey = "1234567"

// Claims 定义一个结构体用于存储信息
type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

// Book 定义图书信息
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Stock  string `json:"stock"`
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
// @Param user body map[string]string true "用户名和密码"
// @Success 200 {object} Response "注册成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 409 {object} Response "用户已存在"
// @Router /register [post]
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "用户已存在"})
		return
	}
	users[req.Username] = req.Password
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body map[string]string true "用户名和密码"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "用户名或密码错误"
// @Router /login [post]
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if pwd, ok := users[req.Username]; !ok || pwd != req.Password {
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
	if _, exists := books[book.ID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Book ID already exists"})
		return
	}
	books[book.ID] = book
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

	// 检查图书是否存在
	if _, exists := books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	delete(books, id)
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

	// 检查图书是否存在
	if _, exists := books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保ID不被修改
	updatedBook.ID = id
	books[id] = updatedBook
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
	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No books available", "books": []Book{}})
		return
	}

	// 将map转换为slice
	var bookList []Book
	for _, book := range books {
		bookList = append(bookList, book)
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

	// 检查图书是否存在
	if book, exists := books[id]; exists {
		c.JSON(http.StatusOK, book)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

// @title 图书管理系统
// @version 1.0
// @description 实现对图书的增删改查的图书管理系统
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

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
