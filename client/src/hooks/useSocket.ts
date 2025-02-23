import { useEffect, useRef, useState } from "react";

export const useSocket = (
  roomName: string | null,
  onMessage?: (msg: any) => void,
) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<NodeJS.Timeout | null>(null);
  const SOCKET_URL = roomName ? `ws://localhost:8080/chat/${roomName}` : null;

  useEffect(() => {
    if (!SOCKET_URL) return;

    const connectSocket = () => {
      if (socketRef.current) {
        console.warn("⚠️ WebSocket already exists. Not reconnecting.");
        return;
      }

      console.log(`🔗 Connecting to WebSocket: ${SOCKET_URL}`);

      // ✅ WebSocket will send cookies automatically (No need to pass extra headers)
      const ws = new WebSocket(SOCKET_URL);

      ws.onopen = () => {
        console.log(`✅ Connected to room: ${roomName}`);
        setSocket(ws);
        socketRef.current = ws;

        if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      };

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          console.log(`📩 New message:`, message);
          if (onMessage) onMessage(message);
        } catch (error) {
          console.error("❌ Error parsing WebSocket message:", error);
        }
      };

      ws.onerror = (error) => {
        console.error("❌ WebSocket Error:", error);
      };

      ws.onclose = () => {
        console.warn("⚠️ WebSocket Disconnected. Attempting to reconnect...");

        socketRef.current = null;
        setSocket(null);

        reconnectTimer.current = setTimeout(connectSocket, 3000);
      };

      socketRef.current = ws;
    };

    connectSocket();

    return () => {
      if (socketRef.current) {
        console.log("🔌 Closing WebSocket connection...");
        socketRef.current.close();
        socketRef.current = null;
      }
      if (reconnectTimer.current) {
        clearTimeout(reconnectTimer.current);
      }
    };
  }, [roomName]); // 🔥 Depend only on `roomName` to avoid unnecessary reconnects

  const sendMessage = (message: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(message);
    } else {
      console.warn("⚠️ Cannot send message. WebSocket is not open.");
    }
  };

  return { socket, sendMessage };
};
