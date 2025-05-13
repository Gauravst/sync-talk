import { MessageProps } from "@/types/messageTypes";
import { useEffect, useRef, useState } from "react";

export const useSocket = (
  roomName: string | null,
  onMessage?: (msg: MessageProps) => void,
) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<NodeJS.Timeout | null>(null);
  const [onlineUsers, setOnlineUsers] = useState<number>(0);
  const SOCKET_URL_ENV = import.meta.env.VITE_REACT_APP_SOCKET_URL;
  const SOCKET_URL = roomName ? `${SOCKET_URL_ENV}/chat/${roomName}` : null;

  useEffect(() => {
    if (!SOCKET_URL) return;

    const connectSocket = () => {
      if (socketRef.current) return;

      const ws = new WebSocket(SOCKET_URL);

      ws.onopen = () => {
        console.log(`Connected to room: ${roomName}`);
        setSocket(ws);
        socketRef.current = ws;
        if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.type === "chat") {
            onMessage?.(data);
          } else {
            setOnlineUsers(data.count);
          }
        } catch (error) {
          console.error("Error parsing WebSocket message:", error);
        }
      };

      ws.onerror = (error) => {
        console.error("WebSocket Error:", error);
      };

      ws.onclose = () => {
        console.warn("WebSocket disconnected. Reconnecting...");
        socketRef.current = null;
        setSocket(null);
        reconnectTimer.current = setTimeout(connectSocket, 3000);
      };

      socketRef.current = ws;
    };

    connectSocket();

    return () => {
      closeSocket();
    };
  }, [roomName]);

  const sendMessage = (message: string) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(message);
    } else {
      console.warn("Cannot send message. WebSocket is not open.");
    }
  };

  const closeSocket = () => {
    if (socketRef.current) {
      console.log("🔌 Manually closing WebSocket...");
      socketRef.current.close();
      socketRef.current = null;
    }
    if (reconnectTimer.current) {
      clearTimeout(reconnectTimer.current);
    }
  };

  return { socket, sendMessage, onlineUsers, closeSocket };
};
