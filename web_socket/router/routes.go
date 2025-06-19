package wsrouter

import (
	"github.com/abdullahshafaqat/Go_Chat_App.git/middelwares"
)

func (r *Router) defineWebSocketRouter() {
	protected := r.Engine.Group("/protected")
	protected.Use(middelwares.WSMiddleware())
	{
		protected.GET("/ws", r.HandleWebSocket())
	}
}
