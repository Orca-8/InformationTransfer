package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

var (
	text    string
	fileArr []string
)

func main() {
	engine := gin.Default()
	engine.Static("/static", "./statics")
	engine.Static("/files", "./wrapFiles")
	engine.LoadHTMLGlob("statics/*")

	// 首页
	engine.GET("/index", func(context *gin.Context) {
		context.HTML(200, "index.html", nil)
	})

	// 上传文本/文件
	engine.POST("/upload", func(context *gin.Context) {
		// 清空text和wrapFiles内容
		text = ""
		err := os.RemoveAll("./wrapFiles/")
		if err != nil {
			text = err.Error()
			context.HTML(200, "get.html", nil)
			return
		}
		// 重新创建wrapFiles文件夹
		err = os.Mkdir("./wrapFiles", 0777)
		if err != nil {
			text = err.Error()
			context.HTML(200, "get.html", nil)
			return
		}
		text = context.PostForm("text")
		// Multipart form
		multipartForm, err := context.MultipartForm()
		if err != nil {
			text = err.Error()
			context.HTML(200, "get.html", nil)
			return
		}
		form, _ := multipartForm, err
		files := form.File["files"]

		fileArr = fileArr[:0]
		for _, file := range files {
			fileArr = append(fileArr, file.Filename)

			// Upload the file to specific dst.
			err := context.SaveUploadedFile(file, "./wrapFiles/"+file.Filename)
			if err != nil {
				text = err.Error()
				context.HTML(200, "get.html", nil)
				return
			}
		}
		context.HTML(200, "get.html", nil)
	})

	// 获取文本/文件
	engine.GET("/get", func(context *gin.Context) {
		context.HTML(200, "get.html", nil)
	})

	// 获取文本
	engine.GET("/getText", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"msg": text,
		})
	})

	// 获取文件
	engine.GET("/getFiles", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"files": fileArr,
		})
	})

	engine.Run(":12345")
}
