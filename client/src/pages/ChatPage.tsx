import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getJoinedRoom, getOldMessage } from "@/services/chatServices";
import { useSocket } from "@/hooks/useSocket";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Hash, MessageCircle, Plus, Search, Send, Users } from "lucide-react";
import Header from "@/components/chats/Header";
import { useAuth } from "@/context/AuthContext";
import { Message, ChatRoom } from "@/types/messageTypes";

function ChatPage() {
  const { user } = useAuth();

  const [message, setMessage] = useState("");
  const { name } = useParams();
  const [chatGroups, setChatGroups] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [messages, setMessages] = useState<Message[]>([]);
  const [initialized, setInitialized] = useState<boolean>(false);
  const [isJoined, setIsJoined] = useState<boolean>(true);
  const navigate = useNavigate();
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
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
          setMessages([]); // Fallback to an empty array
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

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getJoinedRoom();
        if (Array.isArray(rooms)) {
          setChatGroups(rooms);

          if (name) {
            const isRoomJoined = rooms.some((room) => room.name === name);
            setIsJoined(isRoomJoined);
          }

          if (!rooms || rooms.length === 0) navigate("/rooms");
        } else {
          console.error("Expected an array but got:", rooms);
          setChatGroups([]); // Fallback to an empty array
        }
      } catch (error) {
        console.error("Failed to load chat rooms", error);
      } finally {
        setLoading(false);
      }
    };
    fetchRooms();
  }, [name, navigate]);

  const handleSendMessage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!message.trim()) return;

    const newMessageData: Message = {
      userId: user?.userId ?? 0,
      username: user?.username ?? "Unknown",
      roomName: name!,
      content: message,
      time: Date.now(),
    };

    sendMessage(JSON.stringify(newMessageData));
    setMessages((prev) => [...prev, newMessageData]);
    setMessage("");
  };

  const handleGroupClick = (groupName: string) => {
    navigate(`/chat/${groupName}`);
  };

  const handleJoinRoom = () => {
    setIsJoined(true);
  };

  const handleFindNewRooms = () => {
    navigate("/rooms");
  };

  return (
    <div className="flex flex-col min-h-screen bg-background">
      {/* Header */}
      <Header />

      {/* Main Content */}
      <div className="flex-1">
        <div className="p-3 fixed w-full grid grid-cols-1 lg:grid-cols-3 gap-3">
          {/* Sidebar with Rooms */}
          <div className="lg:col-span-1">
            <Card className="h-full border-2 border-muted">
              <CardContent className="p-0">
                <div className="p-4 border-b">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="font-semibold text-lg">Chat Rooms</h3>
                    <Button variant="ghost" size="icon">
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                  <div className="relative">
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input placeholder="Search rooms" className="pl-8" />
                  </div>
                </div>
                <ScrollArea className="h-[calc(100vh-17rem)]">
                  <div className="p-2">
                    <div className="space-y-1">
                      {chatGroups?.map((group) => (
                        <Button
                          key={group?.id}
                          variant={name === group?.name ? "secondary" : "ghost"}
                          className="w-full justify-start gap-2"
                          onClick={() => handleGroupClick(group?.name)}
                        >
                          <Hash className="h-4 w-4" />
                          {group?.name}
                        </Button>
                      ))}
                    </div>
                  </div>
                </ScrollArea>
                <div className="p-4 border-t">
                  <Button className="w-full gap-2" onClick={handleFindNewRooms}>
                    <Plus className="h-4 w-4" />
                    Find New Rooms
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Chat Area */}
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
                              messages.map((msg) => (
                                <div
                                  key={msg?.time} // Use a unique key like `time`
                                  className={`flex ${
                                    msg?.userId === user?.userId
                                      ? "justify-end"
                                      : "justify-start"
                                  } my-2`}
                                >
                                  {msg?.userId !== user?.userId && (
                                    <Avatar className="mr-2">
                                      <AvatarFallback>
                                        {msg?.username?.slice(0, 2) || "U"}
                                      </AvatarFallback>
                                    </Avatar>
                                  )}
                                  <div
                                    className={`p-3 max-w-[75%] rounded-lg ${
                                      msg?.userId === user?.userId
                                        ? "bg-primary text-primary-foreground"
                                        : "bg-muted text-foreground"
                                    }`}
                                  >
                                    {msg?.userId !== user?.userId && (
                                      <p className="text-xs font-medium mb-1">
                                        {msg?.username || "User"}
                                      </p>
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
                                  Be the first to start the conversation in this
                                  room!
                                </p>
                              </div>
                            )}
                            <div ref={messagesEndRef} />
                          </div>
                        ) : (
                          <div className="flex flex-col items-center justify-center h-full text-center p-6">
                            <Hash className="h-12 w-12 text-muted-foreground mb-4" />
                            <h3 className="text-xl font-bold mb-2">
                              Join {name}
                            </h3>
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
                            <Input
                              value={message}
                              onChange={(e) => setMessage(e.target.value)}
                              placeholder="Type your message..."
                              className="flex-1"
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
                    <h3 className="text-2xl font-bold mb-2">
                      No chat selected
                    </h3>
                    <p className="text-muted-foreground mb-6">
                      Select a room from the sidebar or find new rooms to join.
                    </p>
                    <Button onClick={handleFindNewRooms}>Find Rooms</Button>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}

export default ChatPage;
