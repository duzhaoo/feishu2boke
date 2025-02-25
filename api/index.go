package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

const (
	AppID     = "cli_a73f58b52cfb100d"
	AppSecret = "YZAtETOOlJBvfq6lnzMQLF3eGX6Bxk6Z"
	AppToken  = "SnIKbaEXlaAeIDsB9Krc26IYnPh"
	TableID   = "tblfXbKMjzjjk47N"
)

type BlogPost struct {
	Title   string
	Content string
}

func getBlogPosts(client *lark.Client) []BlogPost {
	req := larkbitable.NewListAppTableRecordReqBuilder().
		AppToken(AppToken).
		TableId(TableID).
		Build()

	resp, err := client.Bitable.AppTableRecord.List(context.Background(), req)
	if err != nil {
		fmt.Printf("获取数据失败: %v\n", err)
		return nil
	}

	var posts []BlogPost
	for _, record := range resp.Data.Items {
		var title, content string
		if titleField := record.Fields["AI标题"]; titleField != nil {
			title = titleField.(string)
		}
		if contentField := record.Fields["原创内容"]; contentField != nil {
			content = contentField.(string)
		}
		if title != "" && content != "" {
			posts = append(posts, BlogPost{
				Title:   title,
				Content: content,
			})
		}
	}

	return posts
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// 初始化飞书客户端
	client := lark.NewClient(AppID, AppSecret)

	// 获取博客文章数据
	posts := getBlogPosts(client)

	// 根据路径返回不同的响应
	switch r.URL.Path {
	case "/":
		// 返回JSON格式的博客列表
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"posts": posts,
		})

	case "/api/posts":
		// 返回JSON格式的博客列表
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	default:
		// 检查是否是获取单个博客文章的请求
		if strings.HasPrefix(r.URL.Path, "/post/") {
			idStr := strings.TrimPrefix(r.URL.Path, "/post/")
			id, err := strconv.Atoi(idStr)
			if err == nil && id >= 0 && id < len(posts) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(posts[id])
				return
			}
		}
		
		// 404 Not Found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "页面不存在",
		})
	}
}
