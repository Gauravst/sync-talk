import { useEffect, useRef, useState } from "react";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useNavigate } from "react-router-dom";
import { Hash, MessageCircle } from "lucide-react";

import { ChatRoomProps, MessageProps } from "@/types/messageTypes";
import { useSocket } from "@/hooks/useSocket";
import { useAuth } from "@/context/AuthContext";
import { getOldMessage } from "@/services/chatServices";

import { ImagePreview } from "./ImagePreview";
import { SelectImagePreview } from "./SelectImagePreview";
import { ChatAreaHeader } from "./ChatAreaHeader";
import { ChatAreaFooter } from "./ChatAreaFooter";

type ChatAreaProps = {
  name: string;
  roomData: ChatRoomProps;
  isJoined: boolean;
  setIsJoined: (value: boolean) => void;
};

const ChatArea = ({ name, roomData, isJoined, setIsJoined }: ChatAreaProps) => {
  const { user } = useAuth();
  const navigate = useNavigate();

  const [loading, setLoading] = useState<boolean>(true);
  const [messages, setMessages] = useState<MessageProps[]>([]);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [isUploading, setIsUploading] = useState(false);
  const [initialized, setInitialized] = useState<boolean>(false);
  const [previewPopup, setPreviewPopup] = useState(false);

  console.log(previewUrl);
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

  const handleJoinRoom = () => {
    setIsJoined(true);
  };

  const handleFindNewRooms = () => {
    navigate("/rooms");
  };

  const handleSelectImagePreviewClose = () => {
    setPreviewUrl(null);
  };

  return (
    <div className="col-span-1 md:col-span-2">
      <Card className="h-full flex flex-col border-2 border-muted">
        <CardContent className="p-0 flex-1 relative">
          <SelectImagePreview
            url={previewUrl!}
            open={previewPopup}
            close={handleSelectImagePreviewClose}
          />
          {name ? (
            <>
              <ChatAreaHeader
                user={user!}
                roomData={roomData}
                name={name}
                onlineUsers={onlineUsers}
              />
              <div className="flex flex-col h-[calc(100%-70px)]">
                <ScrollArea className="p-4">
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
                              <Avatar className="mr-2 cursor-pointer">
                                <AvatarFallback>
                                  {msg?.username?.slice(0, 2) || "U"}
                                </AvatarFallback>
                              </Avatar>
                            )}
                            <div
                              className={`p-1 max-w-[75%] rounded-lg ${
                                msg?.userId === user?.id
                                  ? "bg-primary text-primary-foreground rounded-l-lg rounded-br-2xl rounded-tr-none"
                                  : "bg-muted text-foreground rounded-r-lg rounded-tl-none rounded-bl-2xl"
                              }`}
                            >
                              {msg?.userId !== user?.id && (
                                <p className="text-xs font-medium mx-1 mt-1 cursor-pointer hover:underline">
                                  {msg?.username || "User"}
                                </p>
                              )}
                              {msg?.file && (
                                <ImagePreview
                                  file={msg.file}
                                  isUploading={isUploading}
                                  progress={uploadProgress}
                                />
                              )}
                              {msg?.content && (
                                <p className="text-sm p-2">{msg?.content}</p>
                              )}
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="flex flex-col items-center justify-center h-full mt-20 text-center p-6">
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
                      <Hash className="h-12 w-12 text-muted-foreground" />
                      <h3 className="text-xl font-bold mb-2">Join {name}</h3>
                      <p className="text-muted-foreground mb-6">
                        You need to join this room to see messages and
                        participate in the conversation.
                      </p>
                      <Button onClick={handleJoinRoom}>Join Room</Button>
                    </div>
                  )}
                </ScrollArea>
              </div>
            </>
          ) : (
            <div className="flex flex-col items-center justify-center h-full text-center p-6">
              <MessageCircle className="h-16 w-16 text-muted-foreground mb-4" />
              <h3 className="text-2xl font-bold mb-2">No chat selected</h3>
              <p className="text-muted-foreground mb-6">
                Select a room from the sidebar or find new rooms to join.
              </p>
              <Button onClick={handleFindNewRooms}>Find Rooms</Button>
            </div>
          )}
        </CardContent>
        <CardFooter className="p-0">
          {isJoined && (
            <ChatAreaFooter
              user={user!}
              sendMessage={sendMessage}
              setMessages={setMessages}
              previewUrl={previewUrl!}
              setPreviewUrl={setPreviewUrl}
              name={name!}
              setUploadProgress={setUploadProgress}
              setIsUploading={setIsUploading}
              setPreviewPopup={setPreviewPopup}
            />
          )}
        </CardFooter>
      </Card>
    </div>
  );
};

export default ChatArea;
