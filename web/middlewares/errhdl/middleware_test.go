package errhdl

import (
	"awesomeProject/web"
	"bytes"
	"net/http"
	"testing"
	"text/template"
)

func TestMiddlewareHandlesRegisteredError(t *testing.T) {
	page := `
<html>
	<body>
		<h1>404 NOT FOUND</h1>
	</body>
</html>
`
	tpl, err := template.New("404").Parse(page)
	if err != nil {
		t.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, nil)
	if err != nil {
		t.Fatal(err)
	}
	m := NewMiddlewareBuilder().
		RegisterError(404, buffer.Bytes()).
		Build()

	s := web.NewHTTPServer(web.ServerWithMiddleware(m))
	s.Get("/user", func(ctx *web.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	s.Start(":8080")
}
