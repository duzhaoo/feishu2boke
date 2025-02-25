# Vercel 部署规范文档

## 项目结构

```
项目根目录/
├── api/                 # Serverless Functions 目录
│   └── index.go        # API 入口文件
├── public/             # 静态文件目录
│   ├── index.html      # 前端页面
│   └── assets/         # 静态资源
├── go.mod              # Go 依赖管理
├── go.sum              # 依赖版本锁定
└── vercel.json         # Vercel 配置
```

## 代码规范

### 后端代码 (api/index.go)

1. 入口函数规范
```go
// 必须导出 Handler 函数
func Handler(w http.ResponseWriter, r *http.Request) {
    // 处理请求
}
```

2. API 响应规范
```go
// 设置响应头
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")

// 返回数据
json.NewEncoder(w).Encode(data)
```

3. 错误处理规范
```go
if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": err.Error(),
    })
    return
}
```

### 前端代码 (public/index.html)

1. API 调用规范
```javascript
async function fetchData() {
    try {
        const response = await fetch('/api/endpoint');
        const data = await response.json();
        // 处理数据
    } catch (error) {
        console.error('Error:', error);
        // 错误处理
    }
}
```

2. 错误处理规范
```javascript
// 添加加载状态
const loading = document.getElementById('loading');
loading.style.display = 'block';

try {
    // API 调用
} catch (error) {
    // 显示错误信息
} finally {
    loading.style.display = 'none';
}
```

## 部署配置

### vercel.json 配置规范

```json
{
  "version": 2,
  "builds": [
    {
      "src": "api/*.go",
      "use": "@vercel/go"
    },
    {
      "src": "public/**",
      "use": "@vercel/static"
    }
  ],
  "routes": [
    {
      "src": "/api/(.*)",
      "dest": "/api/index.go"
    },
    {
      "src": "/(.*)",
      "dest": "/public/$1"
    }
  ]
}
```

### 环境变量配置

1. 本地开发：创建 `.env` 文件
2. Vercel 部署：使用 Vercel Dashboard 或命令行添加
```bash
vercel env add VARIABLE_NAME
```

