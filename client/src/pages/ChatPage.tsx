import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { ChatRoom, getJoinedRoom } from "@/services/chatServices";
import { useSocket } from "@/hooks/useSocket";

function ChatPage() {
  const [message, setMessage] = useState("");
  const { name } = useParams();
  const [chatGroups, setChatGroups] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [messages, setMessages] = useState<string[]>([]);
  const { socket, sendMessage } = useSocket(name!);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getJoinedRoom();
        setChatGroups(rooms);
      } catch (error) {
        console.error("Failed to load chat rooms");
        console.log(error);
      } finally {
        setLoading(false);
        console.log(loading);
      }
    };

    fetchRooms();
  }, []);

  useEffect(() => {
    if (!socket) return;

    socket.onmessage = (event) => {
      setMessages((prev) => [...prev, event.data]);
    };

    return () => {
      socket.onmessage = null;
    };
  }, [socket]);

  const handleSendMessage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (message.trim()) {
      sendMessage(message);
      setMessages((prev) => [...prev, message]);
      setMessage("");
    }
  };

  return (
    <div className="flex h-screen">
      {/* Sidebar */}
      <Card className="w-64 h-full rounded-none">
        <CardHeader>
          <CardTitle>Chats</CardTitle>
        </CardHeader>
        <ScrollArea className="h-[calc(100vh-60px)]">
          {chatGroups.map((group) => (
            <div
              key={group.id}
              className="p-4 hover:bg-gray-100 cursor-pointer"
            >
              <div className="flex items-center space-x-4">
                <Avatar>
                  <AvatarImage src={group.profilePic} alt={group.name} />
                  <AvatarFallback>{group.name.slice(0, 2)}</AvatarFallback>
                </Avatar>
                <div>
                  <p className="font-medium">{group.name}</p>
                  <p className="text-sm text-gray-500">@{group.name}</p>
                </div>
              </div>
            </div>
          ))}
        </ScrollArea>
      </Card>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col">
        <Card className="flex-1">
          <CardHeader>
            <CardTitle>Chat Group Name</CardTitle>
          </CardHeader>
          <CardContent>
            <ScrollArea className="h-[calc(100vh-200px)]">
              {messages.map((msg, index) => (
                <div key={index} className="p-2 bg-gray-200 my-1 rounded">
                  {msg}
                </div>
              ))}
            </ScrollArea>
          </CardContent>
        </Card>
        <Separator />
        <form onSubmit={handleSendMessage} className="p-4 bg-white">
          <div className="flex space-x-2">
            <Input
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Type your message..."
              className="flex-1"
            />
            <Button type="submit">Send</Button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default ChatPage;
