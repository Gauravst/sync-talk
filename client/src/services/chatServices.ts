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

export const joinChatRoom = async (roomName: string): Promise<boolean> => {
  try {
    const response = await api.post(`/join/${roomName}`);

    // Return true if request is successful (status 200)
    return response.status === 200;
  } catch (error) {
    console.error("Error joining chat room:", error);
    return false; // Return false in case of failure
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
