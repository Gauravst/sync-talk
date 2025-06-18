import api from "./api";
import { UploadedFileResponseProps } from "@/types/fileTypes";
import { AxiosProgressEvent } from "axios";

// upload file in chat
export const uploadFile = async (
  file: File,
  roomName: string,
  message: string,
  onUploadProgress: (progressEvent: AxiosProgressEvent) => void,
): Promise<UploadedFileResponseProps> => {
  try {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("message", message);

    const response = await api.post(`/chat/upload/${roomName}`, formData, {
      onUploadProgress,
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  } catch (error) {
    console.error("Error uploading file:", error);
    throw error;
  }
};
