import { useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ChatRoom, getChatRooms, joinChatRoom } from "@/services/chatServices";
import { useSocket } from "@/hooks/useSocket"; // Updated WebSocket hook

function RoomListPage() {
  const [chatGroups, setChatGroups] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [activeRoom, setActiveRoom] = useState<string | null>(null);
  const [messages, setMessages] = useState<string[]>([]);
  const { socket, sendMessage } = useSocket(activeRoom);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getChatRooms();
        setChatGroups(rooms);
      } catch (error) {
        console.error("Failed to load chat rooms");
      } finally {
        setLoading(false);
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

  const handleJoinClick = async (roomName: string) => {
    try {
      await joinChatRoom(roomName);
      setActiveRoom(roomName);
      alert(`Successfully joined ${roomName}!`);
    } catch {
      alert("Failed to join the room.");
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Your Chat Groups</h1>
      {loading ? (
        <p>Loading chat rooms...</p>
      ) : (
        <ScrollArea className="h-[calc(100vh-100px)]">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {chatGroups.length > 0 ? (
              chatGroups.map((group) => (
                <Card key={group.id}>
                  <CardHeader>
                    <div className="flex items-center space-x-4">
                      <Avatar>
                        <AvatarImage src={group.profilePic} alt={group.name} />
                        <AvatarFallback>
                          {group.name.slice(0, 2)}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <CardTitle>{group.name}</CardTitle>
                        <CardDescription>{group.name}</CardDescription>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p>{45} members</p>
                  </CardContent>
                  <CardFooter>
                    <Button
                      onClick={() => handleJoinClick(group.name)}
                      className="w-full"
                    >
                      Join Chat
                    </Button>
                  </CardFooter>
                </Card>
              ))
            ) : (
              <p>No chat rooms available.</p>
            )}
          </div>
        </ScrollArea>
      )}

      {/* Chat Window */}
      {activeRoom && (
        <div className="mt-6 p-4 border rounded-lg">
          <h2 className="text-xl font-bold mb-2">Chat in {activeRoom}</h2>
          <div className="h-40 overflow-y-auto border p-2 bg-gray-100">
            {messages.map((msg, index) => (
              <div key={index} className="p-1 bg-white my-1 rounded">
                {msg}
              </div>
            ))}
          </div>
          <input
            type="text"
            placeholder="Type a message..."
            className="border p-2 w-full mt-2"
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                sendMessage(e.currentTarget.value);
                e.currentTarget.value = "";
              }
            }}
          />
        </div>
      )}
    </div>
  );
}

export default RoomListPage;
