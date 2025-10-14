class WebSocketService {
  private ws: WebSocket | null = null;
  private onMessageCallback: ((data: any) => void) | null = null;
  private readonly reconnectInterval: number = 5000;
  private userId: string | null = null;
  private token: string | null = null;

  connect(userId: string, token: string) {
    if (this.ws) {
      this.ws.close();
    }

    this.userId = userId;
    this.token = token;
    const url = `${window.env.SERVER_URL}/ws?token=${this.token}`;

    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log("WebSocket connection established.");
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (this.onMessageCallback) {
          this.onMessageCallback(data);
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    };

    this.ws.onclose = () => {
      console.log("WebSocket connection closed. Attempting to reconnect...");
      setTimeout(
        () => this.connect(this.userId!, this.token!),
        this.reconnectInterval
      );
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      this.ws?.close();
    };
  }

  onMessage(callback: (data: any) => void) {
    this.onMessageCallback = callback;
  }

  disconnect() {
    if (this.ws) {
      this.ws.onclose = null;
      this.ws.close();
      this.ws = null;
      console.log("WebSocket disconnected.");
    }
  }
}

export const webSocketService = new WebSocketService();
