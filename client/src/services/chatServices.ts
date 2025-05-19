import api from "./api";
import {
  ChatRoomProps,
  MessageProps,
  PrivateChatRoomProps,
} from "@/types/messageTypes";

// Fetch chat rooms from API
export const getChatRooms = async (): Promise<ChatRoomProps[]> => {
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

// to join chat room
export const joinPrivateChatRoom = async (
  roomName: string,
  roomCode: string,
): Promise<boolean> => {
  try {
    const response = await api.post(`/join/${roomName}/${roomCode}`);
    return response.status === 200;
  } catch (error) {
    console.error("Error joining private chat room:", error);
    return false;
  }
};

// to get joined chat room by user
export const getJoinedRoom = async (): Promise<ChatRoomProps[]> => {
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
): Promise<MessageProps[]> => {
  try {
    const response = await api.get(`/chat/${roomName}/${limit}`);
    console.log(response.data);
    return response.data;
  } catch (error) {
    console.error("Error geting old chat rooms:", error);
    throw error;
  }
};

export const createNewRoom = async (
  username: string,
  description: string,
): Promise<ChatRoomProps> => {
  try {
    const response = await api.post("/room", { name: username, description });
    return response.data;
  } catch (error) {
    console.error("Error creating chat rooms:", error);
    throw error;
  }
};

// to get PrivateChatRoom data using code
export const getPrivateChatRoom = async (
  code: string,
): Promise<PrivateChatRoomProps | null> => {
  try {
    const response = await api.get(`/room/private/${code}`);
    return response.data || null;
  } catch (error) {
    console.error("Error geting Private chat room data :", error);
    throw error;
  }
};

// to join private chat room
export const joinPrivateRoom = async (code: string): Promise<boolean> => {
  try {
    const response = await api.post(`/join/private/${code}`);
    return response.status === 200;
  } catch (error) {
    console.error("Error joining chat room:", error);
    throw error;
  }
};
