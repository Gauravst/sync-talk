import { UploadedFileProps } from "./fileTypes";

export interface MessageProps {
  id?: number;
  userId: number;
  username: string;
  roomName: string;
  content: string;
  profilePic?: string;
  time?: number;
  file? : UploadedFileProps
  createdAt?: string;
  updatedAt?: string;
}

export interface ChatRoomProps {
  id: string;
  name: string;
  description: string;
  members: string;
  userId: number;
  private: boolean;
  code?: string;
}

export interface PrivateChatRoomProps {
  id: string;
  name: string;
  description: string;
  members: number;
  userId: number;
  private: boolean;
  code?: string;
  isMember: boolean;
}
