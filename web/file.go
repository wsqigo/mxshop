package web

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

type StaticResourceHandlerOption func(*StaticResourceHandler)

// WithFileCache 静态文件将会被缓存
// maxFileSizeThreshold 超过这个大小的文件将不会被缓存
// maxCacheFileCnt 最多缓存的文件数
func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(s *StaticResourceHandler) {
		c, err := lru.New(maxCacheFileCnt)
		if err != nil {
			log.Printf("could not create lru cache: %v", err)
			return
		}
		s.maxFileSize = maxFileSizeThreshold
		s.cache = c
	}
}

type StaticResourceHandler struct {
	Dir               string
	extContentTypeMap map[string]string

	cache       *lru.Cache // 缓存文件内容
	maxFileSize int
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	cache, err := lru.New(100)
	if err != nil {
		return nil, err
	}
	res := &StaticResourceHandler{
		Dir:         dir,
		cache:       cache,
		maxFileSize: 10 * 1024 * 1024,
		extContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"css":  "text/css",
			"js":   "application/javascript",
			"html": "text/html",
			"pdf":  "image/pdf",
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	req, err := ctx.PathValue("file").String()
	if err != nil {
		ctx.RespData = []byte("path file error: " + err.Error())
		ctx.RespStatusCode = http.StatusBadRequest
		return
	}

	if item, ok := s.readFileFromData(req); ok {
		log.Printf("从缓存中读取文件: %s", req)
		s.writeItemAsResponse(ctx, item)
		return
	}

	log.Printf("从磁盘读取文件: %s", req)
	file, err := os.Open(req)
	if err != nil {
		ctx.RespData = []byte("open file error: " + err.Error())
		ctx.RespStatusCode = http.StatusInternalServerError
		return
	}
	defer file.Close()

	// 读取数据返回
	ext := filepath.Ext(req)
	dst := filepath.Join(s.Dir, filepath.Clean(req))
	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespData = []byte("read file error: " + err.Error())
		ctx.RespStatusCode = http.StatusInternalServerError
		return
	}

	t, ok := s.extContentTypeMap[ext]
	if !ok {
		ctx.RespData = []byte("unknown content type: " + ext)
		ctx.RespStatusCode = http.StatusInternalServerError
		return
	}
	item := &fileCacheItem{
		fileName:    req,
		fileSize:    len(data),
		data:        data,
		contentType: t,
	}

	s.cacheFile(item)
	s.writeItemAsResponse(ctx, item)
}

func (s *StaticResourceHandler) readFileFromData(req string) (*fileCacheItem, bool) {
	// 从缓存中获取文件
	if s.cache != nil {
		if item, ok := s.cache.Get(req); ok {
			return item.(*fileCacheItem), true
		}
	}
	return nil, false
}

func (s *StaticResourceHandler) writeItemAsResponse(ctx *Context, item *fileCacheItem) {
	ctx.RespStatusCode = http.StatusOK
	ctx.Resp.Header().Set("Content-Type", item.contentType)
	ctx.Resp.Header().Set("Content-Length", strconv.Itoa(item.fileSize))
	ctx.RespData = item.data
}

// 文件缓存项
func (s *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if s.cache != nil && item.fileSize < s.maxFileSize {
		s.cache.Add(item.fileName, item)
	}
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	data        []byte
	contentType string
}
