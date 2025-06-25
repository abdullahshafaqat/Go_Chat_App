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

  /**
   * Connect to the WebSocket server
   */
  connect(token: string, onMessageCallback?: MessageHandler) {
    if (this.ws) return; // Avoid multiple connections
    if (onMessageCallback) {
      this.onMessageCallback = onMessageCallback;
    }

    this.ws = new WebSocket(`ws://localhost:8004/protected/ws?token=${token}`);

    this.ws.onopen = () => {
      console.info("[WS] Connected ✅");
    };

    this.ws.onerror = (err) => {
      console.error("[WS] Error:", err);
    };

    this.ws.onclose = () => {
      console.info("[WS] Closed 🔒");
      this.ws = null;
      this.onMessageCallback = undefined;
    };

    this.ws.onmessage = (evt) => {
      try {
        const data: Message = JSON.parse(evt.data);
        if (this.onMessageCallback) {
          this.onMessageCallback(data); // Dispatch to UI
        }
      } catch (e) {
        console.error("[WS] Error parsing message:", e);
      }
    };
  }

  /**
   * Set the message handler after connection
   */
  setOnMessageCallback(callback: MessageHandler) {
    this.onMessageCallback = callback;
  }

  /**
   * Send a new message
   */
  sendMessage(payload: { receiver_id: number; message: string }) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(payload));
    } else {
      console.error("[WS] Not connected — cannot send message.");
    }
  }

  /**
   * Cleanly close the connection
   */
  disconnect() {
    this.ws?.close();
    this.ws = null;
    this.onMessageCallback = undefined;
  }
}

export const webSocketService = new WebSocketService();
