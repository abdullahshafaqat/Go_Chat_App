package router

func (r *Router) defineWebSocketRouter() {
	r.Engine.GET("/ws/connect", r.StartWebSocketServer)
	r.Engine.GET("/ws/backend", r.SaveBackendConnection)
}
