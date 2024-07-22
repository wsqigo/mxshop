package web

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"text/template"
)

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}
	tpl := template.New("hello-world")
	tpl, err := tpl.Parse("Hello, {{.Name}}!")
	assert.NoError(t, err)

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, User{Name: "Jerry"})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Jerry!", bs.String())
}

func TestMapData(t *testing.T) {
	tpl := template.New("map-data")
	tpl, err := tpl.Parse("Hello {{.Name}}")
	assert.NoError(t, err)

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, map[string]string{"Name": "Jerry"})
	assert.NoError(t, err)
	assert.Equal(t, "Hello Jerry", bs.String())
}

func TestSliceData(t *testing.T) {
	tpl := template.New("slice-data")
	tpl, err := tpl.Parse("Hello, {{index . 0}}")
	assert.NoError(t, err)

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, []string{"Tom", "Jerry"})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Tom", bs.String())
}

const serviceTpl = `
{{- $service := .GenName -}}
type {{ $service }} struct {
	Endpoint string
	Path string
	Client http.Client
}
`

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(firstName, lastName string) string {
	return fmt.Sprintf("Hello, %s %s", firstName, lastName)
}

func TestFuncCall(t *testing.T) {
	tpl := template.New("func-call")
	tpl, err := tpl.Parse(`
切片长度： {{ len .Slice }}
say Hello: {{ .Hello "Tom" "Jerry" }}
打印数字： {{ printf "%.2f" 1.234 }}
`)
	assert.NoError(t, err)

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, FuncCall{Slice: []string{"Tom", "Jerry"}})
	assert.NoError(t, err)
	assert.Equal(t, `
切片长度： 2
say Hello: Hello, Tom Jerry
打印数字： 1.23
`, bs.String())
}

func TestForLoop(t *testing.T) {
	// 用一点小技巧来实现 for i 循环
	tpl := template.New("for-loop")
	tpl, err := tpl.Parse(`
{{ range $idx, $elem := . -}}
下标 {{ $idx }} 元素 {{ $elem }}
{{ end -}}
`)
	assert.NoError(t, err)

	bs := &bytes.Buffer{}
	// 假设我们要从 0 迭代到 100，即 [0, 100)
	// 这里的切片可以是任意类型，比如 []bool, []byte 都可以
	err = tpl.Execute(bs, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	assert.NoError(t, err)
	assert.Equal(t, `
下标 0 元素 0
下标 1 元素 1
下标 2 元素 2
下标 3 元素 3
下标 4 元素 4
下标 5 元素 5
下标 6 元素 6
下标 7 元素 7
下标 8 元素 8
下标 9 元素 9
下标 10 元素 10
`, bs.String())
}

func TestIfElseBlock(t *testing.T) {
	tpl := template.New("if-else-block")
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6) -}}
儿童 0<age<=6
{{- else if and (gt .Age 6) (le .Age 18) -}}
少年 6<age<=18
{{- else if and (gt .Age 18) (le .Age 60) }}
壮年 18<age<=60
{{ else }}
老年 age>60
{{ end -}}
`)
	assert.NoError(t, err)

	bs := &bytes.Buffer{}
	err = tpl.Execute(bs, map[string]interface{}{"Age": 10})
	assert.NoError(t, err)
	assert.Equal(t, "少年 6<age<=18", bs.String())
}

func TestLoginPage(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	assert.NoError(t, err)
	engine := &GoTemplateRender{
		T: tpl,
	}

	s := NewHTTPServer(ServerWithTemplateEngine(engine))
	s.Get("/login", func(ctx *Context) {
		err := ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Println(engine)
		}
	})
	s.Start(":8080")
}
