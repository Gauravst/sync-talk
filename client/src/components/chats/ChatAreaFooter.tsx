import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { UserProps } from "@/types";
import { MessageProps } from "@/types/messageTypes";
import { Button } from "@/components/ui/button";
import { File, Send, Plus } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { uploadFile } from "@/services/fileServices";

type ChatAreaFooterProps = {
  user: UserProps;
  sendMessage: (data: string) => void;
  setMessages: React.Dispatch<React.SetStateAction<MessageProps[]>>;
  previewUrl: string;
  setPreviewUrl: (data: string | null) => void;
  name: string;
  setUploadProgress: (data: number) => void;
  setIsUploading: (data: boolean) => void;
  setPreviewPopup: (data: boolean) => void;
};

export const ChatAreaFooter = ({
  user,
  sendMessage,
  setMessages,
  previewUrl,
  setPreviewUrl,
  name,
  setUploadProgress,
  setIsUploading,
  setPreviewPopup,
}: ChatAreaFooterProps) => {
  const [message, setMessage] = useState("");
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [file2, setFile] = useState<File | null>(null);
  const [fileCloseButton, setFileCloseButton] = useState(!!previewUrl);

  useEffect(() => {
    setFileCloseButton(!!previewUrl);
  }, [previewUrl]);

  const handleSendMessage = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!file2) {
      if (!message.trim()) return;
    }

    const newMessageData: MessageProps = {
      userId: user?.id ?? 0,
      username: user?.username ?? "Unknown",
      roomName: name!,
      content: message,
      time: Date.now(),
    };

    if (file2) {
      newMessageData.file = {
        secureUrl: previewUrl,
      };
    }

    if (!file2) {
      sendMessage(JSON.stringify(newMessageData));
    }
    setMessages((prev: MessageProps[]) => {
      return [...prev, newMessageData];
    });
    setMessage("");

    setPreviewPopup(false);
    setFileCloseButton(false);

    setIsUploading(true);
    if (file2) {
      const response = await uploadFile(
        file2!,
        name!,
        newMessageData.content,
        (progressEvent) => {
          if (progressEvent.total) {
            const percent = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total,
            );
            setUploadProgress(percent);
          }
        },
      );
      if (response.secureUrl) {
        setIsUploading(false);
      }
    }

    console.log("new message");
    console.log(newMessageData);
  };

  const handleFileButtonClick = () => {
    if (fileCloseButton) {
      setPreviewUrl(null);
      setFile(null);
      setFileCloseButton(false);
    } else {
      fileInputRef.current?.click();
    }
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setFileCloseButton(true);
      setPreviewUrl(URL.createObjectURL(file));
      setFile(file);
      setPreviewPopup(true);
    } else {
      setFileCloseButton(false);
    }
  };

  return (
    <div className="flex flex-col w-full">
      <input
        accept="image/*"
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        hidden
      />

      <Separator className="" />
      <form onSubmit={handleSendMessage} className="p-4 w-full flex space-x-2">
        <Button size="icon" type="button" onClick={handleFileButtonClick}>
          {fileCloseButton ? (
            <Plus className="h-4 w-4 rotate-45" />
          ) : (
            <File className="h-4 w-4" />
          )}
        </Button>
        <Input
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type your message..."
          className="h-9 my-0 flex-1"
        />
        <Button type="submit" size="icon">
          <Send className="h-4 w-4" />
          <span className="sr-only">Send</span>
        </Button>
      </form>
    </div>
  );
};
