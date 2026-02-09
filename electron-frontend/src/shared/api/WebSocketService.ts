class WebSocketService {
  private onMessageCallback: ((data: any) => void) | null = null;


  /**
   * Registers a callback function to be called whenever a message is received
   * from the WebSocket bridge.
   *
   * The callback function will be called with the processed message data as an
   * argument. The message data will be processed to convert it into a
   * usable format. For example, if the message data is a Uint8Array,
   * it will be decoded into a string using a TextDecoder.
   *
   * @param callback - The callback function to be called whenever a message is received.
   * @returns A function to unsubscribe from the message event.
   */
  onMessage(callback: (data: any) => void): () => void {
    this.onMessageCallback = callback;

    if ((window as any).appAPI?.websocket) {
      const wrappedCallback = (rawData: any) => {
        try {
          let processedData = JSON.parse(rawData);

          if (processedData?.type === 'Buffer' && Array.isArray(processedData.data)) {
            const decoder = new TextDecoder();
            processedData = decoder.decode(new Uint8Array(processedData.data));
          }

          if (processedData instanceof Uint8Array) {
            const decoder = new TextDecoder();
            processedData = decoder.decode(processedData);
          }

          if (typeof processedData === 'string') {
            try {
              processedData = JSON.parse(processedData);
            } catch (e) {}
          }

          if (this.onMessageCallback) {
            this.onMessageCallback(processedData);
          }
        } catch (error) {
          console.error("Error processing WebSocket message data:", error);
        }
      };

      const unsubscribe = (window as any).appAPI.websocket.onMessage(wrappedCallback);
      return typeof unsubscribe === 'function' ? unsubscribe : () => {};
    }
    
    console.warn("WebSocket Bridge not available");
    return () => {};
  }

  /**
   * Establishes a WebSocket connection to the server using the provided user ID and token.
   * This method will only work if the WebSocket Bridge is available.
   *
   * @param userId - The user ID to use for the WebSocket connection.
   * @param token - The token to use for the WebSocket connection.
   */
  connect(userId: string, token: string) {
    if ((window as any).appAPI?.websocket) {
      const baseUrl = (window as any).env?.SERVER_URL || "http://localhost:8080";
      const wsUrl = baseUrl.replace(/^http/, "ws");
      const url = `${wsUrl}/wss?user_id=${userId}&token=${encodeURIComponent(token)}`;
      
      console.log("Connecting to WebSocket via bridge...");
      (window as any).appAPI.websocket.connect(url);
    }
  }


/**
 * Disconnects from the WebSocket server if the WebSocket Bridge is available.
 * This method will do nothing if the WebSocket Bridge is not available.
 */
  disconnect() {
    if ((window as any).appAPI?.websocket) {
      (window as any).appAPI.websocket.disconnect();
      console.log("WebSocket disconnected.");
    }
  }
}

export const webSocketService = new WebSocketService();