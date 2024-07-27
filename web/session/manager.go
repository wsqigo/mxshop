package session

import (
	"awesomeProject/web"
	"github.com/google/uuid"
)

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

// GetSession 将会尝试从 ctx 中拿到 Session.
// 如果成功了，那么它会将 Session 实例缓存到 ctx 的 UserValues 里面
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}

	val, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return val.(Session), nil
	}
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

	sess, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	ctx.UserValues[m.SessCtxKey] = sess
	return sess, nil
}

// RefreshSession 刷新 session
func (m *Manager) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	// 刷新 session
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

// InitSession 初始化一个 session，并且注入到 http response 里面
func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	// 注入进去 HTTP 响应里面
	err = m.Inject(id, ctx.Resp)
	return sess, err
}

// RemoveSession 删除 session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	// 从 store 中删除
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}

	// 从 http.ResponseWriter 中删除
	return m.Propagator.Remove(ctx.Resp)
}
