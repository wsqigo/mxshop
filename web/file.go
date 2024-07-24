package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type FileUploader struct {
	// FileField 对应于文件表单中的字段名
	FileField string

	// DstPathFunc 用于文件上传后的保存路径
	DstPathFunc func(file *multipart.FileHeader) string
}

func NewFileUploader(fileField string) *FileUploader {
	return &FileUploader{}
}

func (fu *FileUploader) Handle() HandlerFunc {
	return func(ctx *Context) {
		// 从请求中获取文件
		file, fileHeader, err := ctx.Req.FormFile(fu.FileField)
		if err != nil {
			ctx.RespData = []byte("request form file error: " + err.Error())
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		defer file.Close()

		// 保存文件
		dstPath := fu.DstPathFunc(fileHeader)
		err = os.MkdirAll(path.Dir(dstPath), 0755)
		if err != nil {
			ctx.RespData = []byte("create dir error: " + err.Error())
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			ctx.RespData = []byte("create file error: " + err.Error())
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		defer dstFile.Close()

		// 写入文件
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespData = []byte("write to file error: " + err.Error())
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}

		ctx.RespData = []byte("upload success")
		ctx.RespStatusCode = http.StatusOK
	}
}

type FileDownloader struct {
	Dir string
}

func NewFileDownloader(dir string) *FileDownloader {
	return &FileDownloader{
		Dir: dir,
	}
}

func (f *FileDownloader) Handle() HandlerFunc {
	return func(ctx *Context) {
		req, err := ctx.QueryValue("file").String()
		if err != nil {
			ctx.RespData = []byte("query file error: " + err.Error())
			ctx.RespStatusCode = http.StatusBadRequest
			return
		}
		dst := filepath.Join(f.Dir, filepath.Clean(req))
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment; filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}

// -------------------------------静态资源下载-------------------------------
