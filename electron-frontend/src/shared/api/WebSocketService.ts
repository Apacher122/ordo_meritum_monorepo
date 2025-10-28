/**
 * @class WebSocketService
 * Manages a persistent WebSocket connection, including automatic reconnection logic.
 * This service is designed to be a singleton.
 */
class WebSocketService {
  private ws: WebSocket | null = null;
  private onMessageCallback: ((data: any) => void) | null = null;
  private readonly reconnectInterval: number = 5000;
  private userId: string | null = null;
  private token: string | null = null;

  /**
   * Establishes a WebSocket connection to the server. If an existing connection is active,
   * it will be closed before the new one is established. The connection will automatically
   * attempt to reconnect on close.
   * @param {string} userId The user's unique identifier.
   * @param {string} token The authentication token for the user.
   */
  connect(userId: string, token: string) {
    if (this.ws) {
      this.ws.close();
    }

    this.userId = userId;
    this.token = token;

    const wsUrl = window.env.SERVER_URL.replace(/^http/, "ws");
    const url = `${wsUrl}/ws?user_id=${this.userId}&token=${this.token}`;

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

  /**
   * Registers a callback function to be executed when a message is received
   * from the WebSocket server.
   * @param {(data: any) => void} callback The function to execute on message receipt.
   */
  onMessage(callback: (data: any) => void) {
    this.onMessageCallback = callback;
  }

  /**
   * Closes the WebSocket connection and prevents automatic reconnection.
   */
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
