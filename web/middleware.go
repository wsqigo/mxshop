package web

type Middleware func(next HandlerFunc) HandlerFunc

//type Net struct {
//	handlerChains []*HandlerChain
//}
//
//func (n Net) Run(ctx *Context) {
//	wg := sync.WaitGroup{}
//	for _, chain := range n.handlerChains {
//		c := chain
//		if c.concurrent {
//			wg.Add(1)
//			go func() {
//				defer wg.Done()
//				c.Run(ctx)
//			}()
//		} else {
//			c.Run(ctx)
//		}
//	}
//	wg.Done()
//}
//
//type HandlerChain struct {
//	concurrent bool
//	handlers   []HandlerFunc
//}
//
//func (c *HandlerChain) Run(ctx *Context) {
//	for _, handler := range c.handlers {
//		handler(ctx)
//	}
//}
