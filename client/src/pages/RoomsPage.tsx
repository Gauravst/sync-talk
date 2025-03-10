import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { Search, Hash, Users } from "lucide-react";
import Header from "@/components/rooms/Header";
import {
  getChatRooms,
  getJoinedRoom,
  joinChatRoom,
  leaveRoom,
} from "@/services/chatServices";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { ChatRoom } from "@/types/messageTypes";

function RoomsPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const [availableRooms, setAvailableRooms] = useState<ChatRoom[]>([]);
  const [joinedRooms, setJoinedRooms] = useState<ChatRoom[]>([]);
  const [filteredRooms, setFilteredRooms] = useState<ChatRoom[]>([]);
  const [leaveRoomName, setLeaveRoomName] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getChatRooms();
        const joinedRoomsData = await getJoinedRoom();
        if (rooms) setAvailableRooms(rooms);
        if (joinedRoomsData) setJoinedRooms(joinedRoomsData);
      } catch (error) {
        console.error(error);
      }
    };
    fetchRooms();
  }, []);

  useEffect(() => {
    if (availableRooms.length > 0) {
      const results = availableRooms.filter(
        (room) =>
          room.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          room.description.toLowerCase().includes(searchTerm.toLowerCase()),
      );
      setFilteredRooms(results);
    } else {
      setFilteredRooms([]);
    }
  }, [searchTerm, availableRooms]);

  const handleLeaveClick = (roomName: string) => {
    setLeaveRoomName(roomName);
  };

  const confirmLeaveRoom = async () => {
    if (leaveRoomName) {
      const leaved = await leaveRoom(leaveRoomName);
      setLeaveRoomName(null);
      if (leaved) {
        setJoinedRooms((prev) =>
          prev.filter((room) => room.name !== leaveRoomName),
        );
      }
    }
  };

  const handleJoinRoom = async (roomName: string) => {
    const isJoined = await joinChatRoom(roomName);
    if (isJoined) {
      navigate(`/chat/${roomName}`, { replace: true });
    }
  };

  return (
    <div className="w-full items-center flex flex-col min-h-screen bg-background text-white">
      <Header />
      <div className="w-full fixed top-16 container py-8">
        <div className="w-full max-w-4xl mx-auto">
          <div className="mb-8 text-center">
            <h1 className="text-3xl font-bold mb-4">Discover Chat Rooms</h1>
            <p className="text-muted-foreground mb-6">
              Find and join rooms based on your interests and start chatting
              with the community.
            </p>
            <div className="relative max-w-md mx-auto">
              <Search className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search for rooms..."
                className="pl-10"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>
          <ScrollArea className="h-[calc(100vh-16rem)]">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {filteredRooms.map((room) => {
                const isJoined = joinedRooms.some(
                  (joined) => joined.name === room.name,
                );
                return (
                  <Card key={room.id} className="border-2 border-muted">
                    <CardHeader className="pb-2">
                      <div className="flex justify-between items-center">
                        <CardTitle className="flex items-center gap-2">
                          <Hash className="h-5 w-5 text-primary" />
                          {room.name}
                        </CardTitle>
                        <Badge variant="outline" className="gap-1">
                          <Users className="h-3 w-3" />
                          <span>{room.members} members</span>
                        </Badge>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <p className="text-muted-foreground">
                        {room.description}
                      </p>
                    </CardContent>
                    <CardFooter className="flex gap-2">
                      {isJoined ? (
                        <>
                          <Button
                            className="w-full"
                            variant="default"
                            onClick={() => navigate(`/chat/${room.name}`)}
                          >
                            Go to Chat
                          </Button>
                          <Button
                            className="w-full bg-red-600 hover:bg-red-500"
                            variant="destructive"
                            onClick={() => handleLeaveClick(room.name)}
                          >
                            Leave Room
                          </Button>
                        </>
                      ) : (
                        <Button
                          className="w-full"
                          onClick={() => handleJoinRoom(room.name)}
                        >
                          Join Room
                        </Button>
                      )}
                    </CardFooter>
                  </Card>
                );
              })}
            </div>
          </ScrollArea>
        </div>
      </div>
      {leaveRoomName && (
        <Dialog open={true} onOpenChange={() => setLeaveRoomName(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Leave Room</DialogTitle>
            </DialogHeader>
            <p>Are you sure you want to leave {leaveRoomName}?</p>
            <DialogFooter>
              <Button
                variant="secondary"
                onClick={() => setLeaveRoomName(null)}
              >
                Cancel
              </Button>
              <Button
                className="bg-red-600 hover:bg-red-500"
                variant="destructive"
                onClick={confirmLeaveRoom}
              >
                Leave
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

export default RoomsPage;
