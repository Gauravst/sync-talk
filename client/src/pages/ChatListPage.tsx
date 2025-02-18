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
import { group } from "console";

function RoomListPage() {
  const [chatGroups, setChatGroups] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getChatRooms();
        setChatGroups(rooms);
      } catch (error) {
        console.error("Failed to load chat rooms");
        console.log(error);
      } finally {
        setLoading(false);
      }
    };

    fetchRooms();
  }, []);

  const handleJoinClick = async () => {
    try {
      await joinChatRoom(group.name);
      alert(`Successfully joined ${group.name}!`);
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
                      onClick={handleJoinClick}
                      asChild
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
