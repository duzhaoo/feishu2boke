package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

var (
	AppID     = os.Getenv("FEISHU_APP_ID")
	AppSecret = os.Getenv("FEISHU_APP_SECRET")
	AppToken  = os.Getenv("FEISHU_APP_TOKEN")
	TableID   = os.Getenv("FEISHU_TABLE_ID")
)

type BlogPost struct {
	Title   string
	Content string
}

var blogPosts []BlogPost

func getBlogPosts(client *lark.Client) []BlogPost {
	// 构建请求
	req := larkbitable.NewListAppTableRecordReqBuilder().
		AppToken(AppToken).
		TableId(TableID).
		Build()

	// 发送请求
	resp, err := client.Bitable.AppTableRecord.List(context.Background(), req)
	if err != nil {
		fmt.Printf("获取数据失败: %v\n", err)
		return nil
	}

	// 处理响应数据
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

func main() {
	// 初始化飞书客户端
	client := lark.NewClient(AppID, AppSecret)

	// 设置Gin路由
	r := gin.Default()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 设置静态文件路由
	r.Static("/static", "./static")

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		// 获取最新数据
		posts := getBlogPosts(client)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"posts": posts,
		})
	})

	// 博客详情页路由
	r.GET("/post/:id", func(c *gin.Context) {
		// 获取最新数据
		posts := getBlogPosts(client)
		id := c.Param("id")
		if id, err := strconv.Atoi(id); err == nil && id >= 0 && id < len(posts) {
			c.HTML(http.StatusOK, "detail.html", gin.H{
				"post": posts[id],
			})
		} else {
			c.String(http.StatusNotFound, "博客不存在")
		}
	})

	// 启动服务器
	r.Run(":5000")
}

func syncBlogPosts(client *lark.Client) {
	for {
		// 创建请求
		req := larkbitable.NewListAppTableRecordReqBuilder().
			AppToken(AppToken).
			TableId(TableID).
			Build()

		// 发送请求
		resp, err := client.Bitable.AppTableRecord.List(context.Background(), req)
		if err != nil {
			fmt.Printf("同步数据失败: %v\n", err)
		} else {
			// 清空旧数据
			blogPosts = []BlogPost{}

			// 处理响应数据
			for _, record := range resp.Data.Items {
				var title, content string
				if titleField := record.Fields["AI标题"]; titleField != nil {
					title = titleField.(string)
				}
				if contentField := record.Fields["原创内容"]; contentField != nil {
					content = contentField.(string)
				}
				if title != "" && content != "" {
					blogPosts = append(blogPosts, BlogPost{
						Title:   title,
						Content: content,
					})
				}
			}
		}

		// 每5分钟同步一次数据
		time.Sleep(5 * time.Minute)
	}
}