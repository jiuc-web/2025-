package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// 测试GET /book
	getBook("三体")

	// 测试POST /comment
	postComment("小李", "这本书真棒！")
}

func getBook(title string) {
	url := fmt.Sprintf("http://localhost:8080/book?title=%s", title)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("GET请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("GET /book 响应:", string(body))
}

func postComment(user, comment string) {
	requestData := map[string]string{
		"user":    user,
		"comment": comment,
	}
	jsonData, _ := json.Marshal(requestData)

	resp, err := http.Post(
		"http://localhost:8080/comment",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("POST请求失败:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("POST /comment 响应:", result)
}