import api from "./api";

export interface Message {
  id: string;
  sender: string;
  text: string;
  timestamp: string;
}

export interface ChatRoom {
  id: string;
  name: string;
  profilePic: string;
  userId: number;
}

// Get messages from API
// export const getMessages = async (chatId: string): Promise<Message[]> => {
//   const response = await api.get(`/chats/${chatId}/messages`);
//   return response.data;
// };

// Send a new message
export const sendMessage = async (
  chatId: string,
  message: string,
): Promise<Message> => {
  const response = await api.post(`/chats/${chatId}/messages`, { message });
  return response.data;
};

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

export const joinChatRoom = async (roomName: string): Promise<void> => {
  try {
    const response = await api.post(`/join/${roomName}`);
    console.log("Joined room successfully:", response.data);
  } catch (error) {
    console.error("Error joining chat room:", error);
    throw error;
  }
};

export const getJoinedRoom = async (): Promise<ChatRoom[]> => {
  try {
    const response = await api.get(`/join`);
    return response.data;
  } catch (error) {
    console.error("Error joining chat room:", error);
    throw error;
  }
};
