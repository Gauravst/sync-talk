import { useEffect, useRef, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { File, Hash, MessageCircle, Send, Users } from "lucide-react";

import { MessageProps } from "@/types/messageTypes";
import { uploadFile } from "@/services/fileServices";
import { getOldMessage } from "@/services/chatServices";

import { useSocket } from "@/hooks/useSocket";
import { useAuth } from "@/context/AuthContext";
import { useNavigate } from "react-router-dom";
import UploadImagePreview from "./UploadImagePreview";

type ChatAreaProps = {
  name: string;
  isJoined: boolean;
  setIsJoined: (value: boolean) => void;
};

const ChatArea = ({ name, isJoined, setIsJoined }: ChatAreaProps) => {
  const { user } = useAuth();
  const navigate = useNavigate();

  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState<boolean>(true);
  const [messages, setMessages] = useState<MessageProps[]>([]);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewUrl, setPreviewUrl] = useState<string>("");
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [isUploading, setIsUploading] = useState(false);
  const [initialized, setInitialized] = useState<boolean>(false);
  console.log(previewUrl)
  console.log(loading);

  const { sendMessage, onlineUsers } = useSocket(
    name!,
    (newMessageOrHistory) => {
      if (!initialized && Array.isArray(newMessageOrHistory)) {
        setMessages(newMessageOrHistory);
        setInitialized(true);
      } else if (
        typeof newMessageOrHistory === "object" &&
        !Array.isArray(newMessageOrHistory)
      ) {
        setMessages((prev) => [...prev, newMessageOrHistory]);
      }
    },
  );

  useEffect(() => {
    if (!name || name.trim() === "") return;

    const fetchOldMessages = async () => {
      setLoading(true);
      try {
        const data = await getOldMessage(name, 20);
        if (Array.isArray(data)) {
          setMessages(data);
        } else {
          console.error("Expected an array but got:", data);
          setMessages([]);
        }
      } catch (error) {
        console.error("Failed to load chat rooms", error);
      } finally {
        setLoading(false);
      }
    };

    fetchOldMessages();
  }, [name]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSendMessage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!message.trim()) return;

    const newMessageData: MessageProps = {
      userId: user?.id ?? 0,
      username: user?.username ?? "Unknown",
      roomName: name!,
      content: message,
      time: Date.now(),
    };

    sendMessage(JSON.stringify(newMessageData));
    setMessages((prev) => [...prev, newMessageData]);
    setMessage("");
  };

  const handleJoinRoom = () => {
    setIsJoined(true);
  };

  const handleFindNewRooms = () => {
    navigate("/rooms");
  };

  const handleFileButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setIsUploading(true);
      setPreviewUrl(URL.createObjectURL(file));
      const response = await uploadFile(file, name!, (progressEvent) => {
        if (progressEvent.total) {
          const percent = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total,
          );
          setUploadProgress(percent);
        }
      });
      if (response.secureUrl) {
        setIsUploading(false);
      }

      //send message here
      console.log(response.secureUrl);
      console.log("Selected file:", file);
    }
  };
  return (
    <div className="lg:col-span-2">
      <Card className="h-full border-2 border-muted">
        <CardContent className="p-0">
          {name ? (
            <>
              <div className="flex items-center justify-between border-b px-4 py-3">
                <div className="flex items-center gap-2">
                  <Hash className="h-5 w-5 text-muted-foreground" />
                  <h3 className="font-semibold">{name}</h3>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className="gap-1">
                    <Users className="h-3 w-3" />
                    <span>{onlineUsers} online</span>
                  </Badge>
                </div>
              </div>

              <div className="flex h-[calc(100vh-8rem)] flex-col">
                <ScrollArea className="flex-1 p-4">
                  {isJoined ? (
                    <div className="space-y-4">
                      {messages.length > 0 ? (
                        messages.map((msg, index) => (
                          <div
                            key={index}
                            className={`flex ${
                              msg?.userId === user?.id
                                ? "justify-end"
                                : "justify-start"
                            } my-2`}
                          >
                            {msg?.userId !== user?.id && (
                              <Avatar className="mr-2">
                                <AvatarFallback>
                                  {msg?.username?.slice(0, 2) || "U"}
                                </AvatarFallback>
                              </Avatar>
                            )}
                            <div
                              className={`p-3 max-w-[75%] rounded-lg ${
                                msg?.userId === user?.id
                                  ? "bg-primary text-primary-foreground"
                                  : "bg-muted text-foreground"
                              }`}
                            >
                              {msg?.userId !== user?.id && (
                                <p className="text-xs font-medium mb-1">
                                  {msg?.username || "User"}
                                </p>
                              )}
                              {msg?.file && (
                                <UploadImagePreview
                                  file={msg.file}
                                  isUploading={isUploading}
                                  progress={uploadProgress}
                                />
                              )}
                              <p className="text-sm">{msg?.content}</p>
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="flex flex-col items-center justify-center h-full text-center p-6">
                          <MessageCircle className="h-12 w-12 text-muted-foreground mb-4" />
                          <h3 className="text-xl font-bold mb-2">
                            No messages yet
                          </h3>
                          <p className="text-muted-foreground">
                            Be the first to start the conversation in this room!
                          </p>
                        </div>
                      )}
                      <div ref={messagesEndRef} />
                    </div>
                  ) : (
                    <div className="flex flex-col items-center justify-center h-full text-center p-6">
                      <Hash className="h-12 w-12 text-muted-foreground mb-4" />
                      <h3 className="text-xl font-bold mb-2">Join {name}</h3>
                      <p className="text-muted-foreground mb-6">
                        You need to join this room to see messages and
                        participate in the conversation.
                      </p>
                      <Button onClick={handleJoinRoom}>Join Room</Button>
                    </div>
                  )}
                </ScrollArea>

                {isJoined && (
                  <>
                    <Separator />
                    <form
                      onSubmit={handleSendMessage}
                      className="p-4 flex space-x-2"
                    >
                      <input
                        type="file"
                        ref={fileInputRef}
                        onChange={handleFileChange}
                        hidden
                      />
                      <Button
                        size="icon"
                        type="button"
                        onClick={handleFileButtonClick}
                      >
                        <File className="h-4 w-4" />
                      </Button>
                      <Input
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        placeholder="Type your message..."
                        className="flex-1 min:h-4"
                      />
                      <Button type="submit" size="icon">
                        <Send className="h-4 w-4" />
                        <span className="sr-only">Send</span>
                      </Button>
                    </form>
                  </>
                )}
              </div>
            </>
          ) : (
            <div className="flex flex-col items-center justify-center h-full text-center p-6">
              <MessageCircle className="h-16 w-16 text-muted-foreground mb-4 mt-16" />
              <h3 className="text-2xl font-bold mb-2">No chat selected</h3>
              <p className="text-muted-foreground mb-6">
                Select a room from the sidebar or find new rooms to join.
              </p>
              <Button onClick={handleFindNewRooms}>Find Rooms</Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default ChatArea;
