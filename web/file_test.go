package web

import (
	"log"
	"mime/multipart"
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	path := "testdata/my_file.txt"
	// 打开一个文件，只能读
	file, err := os.Open(path)
	defer file.Close()
	assert.NoError(t, err)
	assert.Equal(t, "testdata/my_file.txt", file.Name())

	data := make([]byte, 1024)
	n, err := file.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, 18, n)
	assert.Equal(t, "这是我的文件", string(data[:n]))
}

func TestCreateFile(t *testing.T) {
	path := "testdata/my_file_1.txt"

	// 创建一个文件，如果文件已经存在，会清空原有内容
	file, err := os.Create(path)
	assert.NoError(t, err)
	defer file.Close()
	assert.Equal(t, "testdata/my_file_1.txt", file.Name())
}

func TestUpload(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	assert.NoError(t, err)

	engine := &GoTemplateRender{T: tpl}
	s := NewHTTPServer(ServerWithTemplateEngine(engine))

	s.Get("/upload", func(ctx *Context) {
		err = ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})

	fu := FileUploader{
		FileField: "myfile",
		DstPathFunc: func(file *multipart.FileHeader) string {
			return path.Join("testdata", "upload", file.Filename)
		},
	}
	s.Post("/upload", fu.Handle())
	s.Start(":8081")
}

func TestFileDownloader(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/download", (&FileDownloader{
		// 下载的文件所在目录
		Dir: "testdata/upload"}).Handle())

	s.Start(":8081")
}

func TestStaticResourceHandler(t *testing.T) {
	s := NewHTTPServer()

	handler, err := NewStaticResourceHandler("./testdata/img")
	assert.NoError(t, err)
	s.Get("/img/:file", handler.Handle)
	// 在浏览器中访问:http://localhost:8081/img/come_on_baby.png
	s.Start(":8081")
}
