import api from "./api";
import { ChatRoom, Message } from "@/types/messageTypes";

// Fetch chat rooms from API
export const getChatRooms = async (): Promise<ChatRoom[]> => {
  try {
    const response = await api.get("/room");
    return response.data;
  } catch (error) {
    console.error("Error fetching chat rooms:", error);
    throw error;
  }
};

// to join chat room
export const joinChatRoom = async (roomName: string): Promise<boolean> => {
  try {
    const response = await api.post(`/join/${roomName}`);
    return response.status === 200;
  } catch (error) {
    console.error("Error joining chat room:", error);
    return false;
  }
};

// to get joined chat room by user
export const getJoinedRoom = async (): Promise<ChatRoom[]> => {
  try {
    const response = await api.get(`/join`);
    return response.data;
  } catch (error) {
    console.error("Error joining chat room:", error);
    throw error;
  }
};

// leave joined chat room
export const leaveRoom = async (roomName: string): Promise<boolean> => {
  try {
    const response = await api.delete(`/join/${roomName}`);
    return response.status === 200;
  } catch (error) {
    console.error("Error joining chat room:", error);
    throw error;
  }
};

// get old chats
export const getOldMessage = async (
  roomName: string,
  limit: number,
): Promise<Message[]> => {
  try {
    const response = await api.get(`/chat/${roomName}/${limit}`);
    return response.data;
  } catch (error) {
    console.error("Error geting old chat rooms:", error);
    throw error;
  }
};
