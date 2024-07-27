package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any) error
	ID() string
}

// Store 管理 Session 本身
type Store interface {
	// Generate 生成一个新的 Session
	Generate(ctx context.Context, id string) (Session, error)
	// Refresh 这种设计是一直用同一个 id 的
	// 如果想支持 Refresh 换 ID，那么可以重新生成一个，并移除原有的
	//Refresh(ctx context.Context, id string) (Session, error)
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

type Propagator interface {
	// Inject 将 session id 注入到里面
	// Inject 必须是幂等的
	Inject(id string, writer http.ResponseWriter) error
	// Extract 将 session id 从 http.Request 中提取出来
	// 例如从 cookie 中将 session id 提取出来
	Extract(req *http.Request) (string, error)

	// Remove 从 http.ResponseWriter 中删除 session id
	// 例如从 cookie 中删除 session id
	Remove(writer http.ResponseWriter) error
}
