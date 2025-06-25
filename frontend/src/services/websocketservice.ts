// src/services/websocketservice.ts
type Message = {
  sender_id: number;
  receiver_id: number;
  message: string;
  timestamp: string;
};

type MessageHandler = (message: Message) => void;

class WebSocketService {
  private ws: WebSocket | null = null;
  private onMessageCallback?: MessageHandler;

  connect(token: string, onMessageCallback: MessageHandler) {
    this.onMessageCallback = onMessageCallback;
    this.ws = new WebSocket(`ws://localhost:8004/protected/ws?token=${token}`);

    this.ws.onmessage = (evt) => {
      try {
        const data: Message = JSON.parse(evt.data);
        if (this.onMessageCallback) {
          this.onMessageCallback(data); // Push into UI
        }
      } catch (e) {
        console.error("Error parsing WS message:", e);
      }
    };

    this.ws.onerror = (err) => console.error("WebSocket error:", err);
    this.ws.onclose = () => console.info("WebSocket closed");
  }

  sendMessage(payload: { receiver_id: number; message: string }) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(payload));
    } else {
      console.error("WebSocket not connected.");
    }
  }

  disconnect() {
    this.ws?.close();
    this.ws = null;
  }
}

export const webSocketService = new WebSocketService();
