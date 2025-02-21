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
import { ChatRoom, getChatRooms, joinChatRoom } from "../services/chatServices";
import { useNavigate } from "react-router-dom";

function RoomListPage() {
  const [chatGroups, setChatGroups] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getChatRooms();
        setChatGroups(rooms);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    fetchRooms();
  }, []);

  const handleJoinClick = async (roomName: string) => {
    const isJoined = await joinChatRoom(roomName);

    if (isJoined) {
      console.log("Navigating to chat room...");
      navigate(`/chat/${roomName}`);
    } else {
      console.log("Failed to join the room.");
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">List of Chat Rooms</h1>
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
    </div>
  );
}

export default RoomListPage;
