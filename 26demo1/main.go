package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book 结构体定义图书信息
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

// 全局变量存储图书数据
var books = make(map[string]Book)

// AddBook 添加图书
// @Summary 添加书籍
// @Description 传入书籍信息新增书籍
// @Tags 图书
// @Accept json
// @Produce json
// @Param book body Book true "书籍信息"
// @Success 200 {object} Response{data=Book} "书籍信息"
// @Failure 400 {object} Response "参数解析失败"
// @Failure 409 {object} Response "书籍ID已存在"
// @Router /book/add [post]
func AddBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查ID是否已存在
	if _, exists := books[book.ID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Book ID already exists"})
		return
	}

	books[book.ID] = book
	c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully", "book": book})
}

// DeleteBook 删除图书
// @Summary 删除书籍
// @Description 根据书籍 ID 删除书籍
// @Tags 图书
// @Accept json
// @Produce json
// @Param id path string true "书籍ID"
// @Success 200 {object} Response "书籍删除成功"
// @Router /book/delete/{id} [delete]
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
// UpdateBook 更新书籍
// @Summary 更新书籍
// @Description 根据 ID 更新书籍信息
// @Tags 图书
// @Accept json
// @Produce json
// @Param id path string true "书籍ID"
// @Param book body Book true "更新后的书籍信息"
// @Success 200 {object} Response{data=Book} "书籍更新成功"
// @Failure 400 {object} Response "参数绑定失败"
// @Failure 400 {object} Response "路径ID与请求体ID不一致"
// @Router /book/update/{id} [put]
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
// @Success 200 {object} map[string]Book
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

	// 设置路由
	r.POST("/books", AddBook)          // 添加图书
	r.DELETE("/books/:id", DeleteBook) // 删除图书
	r.PUT("/books/:id", UpdateBook)    // 更新图书
	r.GET("/books", SearchAllBook)     // 查询所有图书
	r.GET("/books/:id", GetBook)       // 根据ID查询图书

	r.Run(":8080")
}
