import { useEffect, useState } from "react";

export const useSocket = (roomName: string | null) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const SOCKET_URL = roomName ? `ws://localhost:8080/chat/${roomName}` : null;

  useEffect(() => {
    if (!SOCKET_URL) return;

    const ws = new WebSocket(SOCKET_URL);
    setSocket(ws);

    ws.onopen = () => {
      console.log(`Connected to room: ${roomName}`);
    };

    ws.onmessage = (event) => {
      console.log(`New message: ${event.data}`);
    };

    ws.onerror = (error) => {
      console.error("WebSocket Error:", error);
    };

    ws.onclose = () => {
      console.log("WebSocket Disconnected");
    };

    return () => {
      ws.close();
    };
  }, [SOCKET_URL]);

  const sendMessage = (message: string) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(message);
    }
  };

  return { socket, sendMessage };
};
