import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { UserProps } from "@/types";
import { MessageProps } from "@/types/messageTypes";
import { Button } from "@/components/ui/button";
import { File, Send, Plus } from "lucide-react";
import { useRef, useState } from "react";
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
  const [file, setFile] = useState<File | null>(null);
  const [fileCloseButton, setFileCloseButton] = useState(false);

  const handleSendMessage = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!file) {
      if (!message.trim()) return;
    }

    const newMessageData: MessageProps = {
      userId: user?.id ?? 0,
      username: user?.username ?? "Unknown",
      roomName: name!,
      content: message,
      time: Date.now(),
    };

    // if (previewUrl) {
    //   URL.revokeObjectURL(previewUrl);
    // }

    if (file) {
      newMessageData.content = "";
      newMessageData.file = {
        secureUrl: previewUrl,
      };
    }

    console.info("content--------");
    console.log(newMessageData.content);

    if (!file) {
      sendMessage(JSON.stringify(newMessageData));
    }
    setMessages((prev: MessageProps[]) => {
      return [...prev, newMessageData];
    });
    setMessage("");

    setPreviewPopup(false);
    setFileCloseButton(false);

    //image upload here
    setIsUploading(true);
    if (file) {
      const response = await uploadFile(
        file!,
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
  };

  const handleFileButtonClick = () => {
    if (file) {
      setPreviewUrl(null);
      setFile(null);
    } else {
      fileInputRef.current?.click();
    }
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setFileCloseButton(true);
      setIsUploading(true);
      setPreviewUrl(URL.createObjectURL(file));
      setFile(file);
      setPreviewPopup(true);
      console.log("Selected file:", file);
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
