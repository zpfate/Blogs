package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 定义返回的数据结构
type Response struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`
}

type ApiResponse struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is home page.")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // 返回 200 状态码
		fmt.Fprint(w, "OK")
	})

	// 路由 3: API 接口，返回 JSON
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res := ApiResponse{
			Path:    "/api/data",
			Message: "这是从 Docker 容器返回的 JSON 数据 支持热重载",
		}

		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头为 JSON 格式
		w.Header().Set("Content-Type", "application/json")

		// 获取参数
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Developer"
		}

		// 创建返回的对象
		data := Response{
			Message:   fmt.Sprintf("Hello world, %s!", name),
			Status:    200,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 将对象编码为 JSON 并写入响应
		json.NewEncoder(w).Encode(data)
	})

	fmt.Println("JSON API Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
