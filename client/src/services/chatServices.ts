import api from "./api";

export interface Message {
  id: string;
  sender: string;
  text: string;
  timestamp: string;
}

// Get messages from API
export const getMessages = async (chatId: string): Promise<Message[]> => {
  const response = await api.get(`/chats/${chatId}/messages`);
  return response.data;
};

// Send a new message
export const sendMessage = async (
  chatId: string,
  message: string,
): Promise<Message> => {
  const response = await api.post(`/chats/${chatId}/messages`, { message });
  return response.data;
};
