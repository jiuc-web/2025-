package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Comment struct {
	User    string `json:"user"`
	Comment string `json:"comment"`
}

type Response struct {
	Message string `json:"message"`
	User    string `json:"user"`
	Comment string `json:"comment"`
}

// get/book处理函数
func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "缺少title参数", http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("您正在查询图书: 《%s》", title)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, response)
}

// post/comment处理函数
func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "无效的数据", http.StatusBadRequest)
		return
	}
	response := Response{
		Message: "评论提交成功",
		User:    comment.User,
		Comment: comment.Comment,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/comment", commentHandler)
	fmt.Println("服务器正在监听8080端口...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}