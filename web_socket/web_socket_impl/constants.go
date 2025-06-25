package websocketservice

import "time"

const (
	
	WriteWait      = 10 * time.Second
	PongWait       = 60 * time.Second
	PingInterval   = (PongWait * 9) / 10 
	MaxMessageSize = 1024 * 1024         
	MaxRetryCount  = 3                  

	
	CloseNormalClosure   = 1000
	CloseGoingAway       = 1001
	CloseProtocolError   = 1002
	CloseUnsupportedData = 1003
)
