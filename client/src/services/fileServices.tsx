import api from "./api";
import { UploadedFileProps } from "@/types/fileTypes";
import { AxiosProgressEvent } from "axios";

// upload file in chat
export const uploadFile = async (
  file: File,
  roomName: string,
  onUploadProgress: (progressEvent: AxiosProgressEvent) => void,
): Promise<UploadedFileProps> => {
  try {
    const formData = new FormData();
    formData.append("file", file);

    const response = await api.post(`/chat/upload/${roomName}`, formData, {
      onUploadProgress,
    });
    return response.data;
  } catch (error) {
    console.error("Error uploading file:", error);
    throw error;
  }
};
