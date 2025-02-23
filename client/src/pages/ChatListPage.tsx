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
import {
  ChatRoom,
  getChatRooms,
  getJoinedRoom,
  joinChatRoom,
  leaveRoom,
} from "../services/chatServices";
import { useNavigate } from "react-router-dom";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";

function RoomListPage() {
  const [chatGroups, setChatGroups] = useState<ChatRoom[]>([]);
  const [joinedRooms, setJoinedRooms] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [leaveRoomName, setLeaveRoomName] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const rooms = await getChatRooms();
        const joinedRoomsData = await getJoinedRoom();
        if (rooms) setChatGroups(rooms);
        if (joinedRoomsData) setJoinedRooms(joinedRoomsData);
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
      navigate(`/chat/${roomName}`);
    } else {
      console.log("Failed to join the room.");
    }
  };

  const handleLeaveClick = (roomName: string) => {
    setLeaveRoomName(roomName); // Open confirmation popup
  };

  const confirmLeaveRoom = async () => {
    if (leaveRoomName) {
      console.log(`Leaving room: ${leaveRoomName}`);
      const leaved = await leaveRoom(leaveRoomName);
      setLeaveRoomName(null);
      if (leaved) {
        navigate(`/rooms`);
      } else {
        console.log("error in leaveing room");
      }
    }
  };

  if (loading) {
    return (
      <div className="container flex flex-col items-center mx-auto p-4">
        <h1 className="text-2xl font-bold mt-6 mb-4">List of Chat Rooms</h1>
        <div className="text-lg text-gray-500">Loading...</div>
      </div>
    );
  }

  if (chatGroups.length === 0) {
    return (
      <div className="container flex flex-col items-center mx-auto p-4">
        <h1 className="text-2xl font-bold mt-6 mb-4">List of Chat Rooms</h1>
        <div className="text-lg text-gray-500 mt-10">No rooms available.</div>
      </div>
    );
  }

  return (
    <div className="container flex w-full items-center flex-col mx-auto p-4">
      <h1 className="text-2xl font-bold mt-6 mb-4">List of Chat Rooms</h1>
      <ScrollArea className="h-[calc(100vh-100px)] w-full">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {chatGroups.map((group) => {
            const isJoined = joinedRooms.some(
              (room) => room.name === group.name,
            );

            return (
              <Card key={group.id}>
                <CardHeader>
                  <div className="flex items-center space-x-4">
                    <Avatar>
                      <AvatarImage src={group.profilePic} alt={group.name} />
                      <AvatarFallback>{group.name?.slice(0, 2)}</AvatarFallback>
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
                <CardFooter className="flex gap-2">
                  {isJoined ? (
                    <>
                      <Button
                        onClick={() => navigate(`/chat/${group.name}`)}
                        className="w-full"
                      >
                        Go to Chat
                      </Button>
                      <Button
                        variant="destructive"
                        onClick={() => handleLeaveClick(group.name)}
                        className="w-full"
                      >
                        Leave
                      </Button>
                    </>
                  ) : (
                    <Button
                      onClick={() => handleJoinClick(group.name)}
                      className="w-full"
                    >
                      Join Chat
                    </Button>
                  )}
                </CardFooter>
              </Card>
            );
          })}
        </div>
      </ScrollArea>

      {/* Leave Confirmation Popup */}
      {leaveRoomName && (
        <Dialog
          open={!!leaveRoomName}
          onOpenChange={() => setLeaveRoomName(null)}
        >
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Confirm Leave</DialogTitle>
            </DialogHeader>
            <p>
              Are you sure you want to leave <strong>{leaveRoomName}</strong>?
            </p>
            <DialogFooter>
              <Button variant="outline" onClick={() => setLeaveRoomName(null)}>
                Cancel
              </Button>
              <Button variant="destructive" onClick={confirmLeaveRoom}>
                Confirm Leave
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}

export default RoomListPage;
