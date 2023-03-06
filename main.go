package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

type Settings struct {
	Openai string `json:"openai"`
}

var settings Settings
var Logger *log.Logger

func init() {
	// 初始化日志
	fmt.Printf("初始化日志文件\n")
	logFile, err := os.OpenFile("./log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	// 初始化配置
	fmt.Printf("初始化配置\n")
	filePtr, err := os.Open("./settings.json")
	if err != nil {
		fmt.Printf("文件打开失败 [Err:%s]\n", err.Error())
		Logger.Printf("文件打开失败 [Err:%s]\n", err.Error())
		panic(err)
	}
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&settings)
	if err != nil {
		fmt.Println("解码失败", err.Error())
		Logger.Printf("解码失败 [Err:%s]\n", err.Error())
		panic(err)
	} else {
		fmt.Println("解码成功")
		Logger.Println("解码成功")
		fmt.Println(settings)
	}
	defer filePtr.Close()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		content := c.Query("content")
		if content == "" {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		} else {
			data, err := chat(content)
			if err != nil {
				c.JSON(200, gin.H{
					"message": data,
				})
			} else {
				c.JSON(400, gin.H{
					"message": data,
				})
			}
		}

	})
	r.Run("0.0.0.0:8012") // 监听并在 0.0.0.0:8080 上启动服务
}

func chat(content string) (string, error) {
	var data = map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": content,
			},
		},
	}
	b, err := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(b))
	// 比如说设置个token
	req.Header.Set("Authorization", "Bearer "+settings.Openai)
	// 再设置个json
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	Logger.Println(string(respBody))
	return string(respBody), nil
}
