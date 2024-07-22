package web

import (
	"bytes"
	"context"
	"text/template"
)

type TemplateEngine interface {
	// Render 渲染页面
	// data 是渲染页面所需要的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateRender struct {
	T *template.Template
	// 也可以考虑设计为 map[string]*template.Template
	// 但是其实没太大必要，因为 template.Template 本身就提供了按名索引的功能
}

func (r *GoTemplateRender) Render(ctx context.Context,
	tplName string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := r.T.ExecuteTemplate(res, tplName, data)
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
