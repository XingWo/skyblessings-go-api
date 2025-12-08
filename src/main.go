package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

// Config 配置结构
type Config struct {
	Server struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		LogLevel string `toml:"log_level"`
	} `toml:"server"`

	Image struct {
		Width     int    `toml:"width"`
		Height    int    `toml:"height"`
		FontSize  int    `toml:"font_size"`
		AssetsDir string `toml:"assets_dir"` // 资源文件夹的绝对路径
	} `toml:"image"`
}

var config Config

func main() {
	// 切换到可执行文件所在目录
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		os.Chdir(exeDir)
		log.Printf("📁 工作目录: %s", exeDir)
	}

	// 加载配置
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	log.Printf("📂 资源目录: %s", config.Image.AssetsDir)

	// 初始化数据
	initDrawData()

	// 设置 Gin 为 release 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// API 路由
	r.GET("/blessing", handleGetBlessing)
	r.GET("/", handleIndex)

	// 启动服务
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("🌟 祈福签 API 服务启动成功！")
	log.Printf("📍 访问地址: http://localhost:%d/blessing", config.Server.Port)
	log.Printf("🎨 图片尺寸: %dx%d", config.Image.Width, config.Image.Height)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadConfig 加载配置文件，如果不存在则创建默认配置
func loadConfig() error {
	configPath := "config.toml"

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("配置文件不存在，正在创建默认配置...")
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("创建默认配置失败: %w", err)
		}
		log.Printf("✅ 已创建默认配置文件: %s", configPath)
	}

	// 读取配置
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(path string) error {
	defaultConfig := `# 祈福签 API 配置文件

[server]
host = "0.0.0.0"
port = 51205
log_level = "info" # 日志级别 (info, debug)

[image]
width = 1240
height = 620
font_size = 40
# 资源文件夹路径（绝对路径或相对路径）
# 目录结构要求：
# assets/
#   ├── font/
#   │   └── LXGWWenKaiMono-Medium.ttf
#   └── image/
#       ├── background.png
#       └── ...
assets_dir = "../assets"
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

// handleIndex 首页
func handleIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "光遇祈福签 API",
		"version": "1.0.0",
		"endpoints": gin.H{
			"GET /blessing": "获取随机祈福签图片",
		},
	})
}

// handleGetBlessing 处理获取祈福签请求
func handleGetBlessing(c *gin.Context) {
	// 生成图片
	imgData, err := generateBlessingImage()
	if err != nil {
		log.Printf("生成图片失败: %v", err)
		c.JSON(500, gin.H{"error": "生成图片失败"})
		return
	}

	// 返回图片
	c.Data(200, "image/png", imgData)
}

// getAssetPath 获取资源文件的路径
func getAssetPath(subPath string) string {
	return filepath.Join(config.Image.AssetsDir, subPath)
}
