export interface Message {
  id?: number;
  userId: number;
  username: string;
  roomName: string;
  content: string;
  profilePic?: string;
  time?: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface ChatRoom {
  id: string;
  name: string;
  description: string;
  members: string;
  userId: number;
}
