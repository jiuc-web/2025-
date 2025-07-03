package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// 数据库连接参数
var (
	dbServer   = getEnv("DB_SERVER", "db")
	dbPort, _  = strconv.Atoi(getEnv("DB_PORT", "1433"))
	dbUser     = getEnv("DB_USER", "sa")
	dbPassword = getEnv("DB_PASSWORD", "YourStrong!Passw0rd")
	dbName     = getEnv("DB_NAME", "MUXI")
)

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

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

// 自动创建数据库
func ensureDatabaseExists() error {
	// 先连 master 数据库
	masterDsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=master&encrypt=disable",
		dbUser, dbPassword, dbServer, dbPort)
	masterDb, err := gorm.Open(sqlserver.Open(masterDsn), &gorm.Config{})
	if err != nil {
		return err
	}
	var exists int
	row := masterDb.Raw("SELECT COUNT(*) FROM sys.databases WHERE name = ?", dbName).Row()
	row.Scan(&exists)
	if exists == 0 {
		if err := masterDb.Exec("CREATE DATABASE [" + dbName + "]").Error; err != nil {
			return err
		}
	}
	sqlDB, _ := masterDb.DB()
	sqlDB.Close()
	return nil
}

// 生成Token
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

// JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"message": "未认证"})
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
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

// 用户注册
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

// 用户登录
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

// 添加图书
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

// 删除图书
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

// 更新图书
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

// 查询所有图书
func SearchAllBook(c *gin.Context) {
	var bookList []Book
	if err := db.Find(&bookList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": bookList})
}

// 根据ID查询图书
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
	// 等待数据库服务可用并确保数据库存在
	for i := 1; i <= 20; i++ {
		err := ensureDatabaseExists()
		if err == nil {
			break
		}
		fmt.Printf("数据库创建检测失败，第%d次重试，错误信息：%v\n", i, err)
		time.Sleep(3 * time.Second)
	}

	// 连接业务数据库
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
		dbUser, dbPassword, dbServer, dbPort, dbName)

	var err error
	for i := 1; i <= 20; i++ {
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("数据库连接失败，第%d次重试，错误信息：%v\n", i, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	_ = db.AutoMigrate(&Book{}, &User{})

	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)

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
