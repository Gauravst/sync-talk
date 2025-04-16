export interface MessageProps {
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

export interface ChatRoomProps {
  id: string;
  name: string;
  description: string;
  members: string;
  userId: number;
}
