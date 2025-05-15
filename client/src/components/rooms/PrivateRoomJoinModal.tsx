import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState, useEffect } from "react";
import { Hash, Users, Lock, Loader2, Search } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import {
  getPrivateChatRoom,
  joinPrivateRoom,
  leaveRoom,
} from "@/services/chatServices";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PrivateChatRoomProps } from "@/types/messageTypes";
import { useAuth } from "@/context/AuthContext";

type ModalProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
};

const PrivateRoomJoinModal = ({ open, setOpen }: ModalProps) => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [roomCode, setRoomCode] = useState("");
  const [loading, setLoading] = useState(false);
  const [disable, setDisable] = useState(true);
  const [room, setRoom] = useState<PrivateChatRoomProps | null>(null);
  const [leaveRoomName, setLeaveRoomName] = useState<string | null>(null);
  const [searched, setSearched] = useState(false);

  useEffect(() => {
    if (roomCode.trim()) {
      setDisable(false);
    } else {
      setDisable(true);
    }
  }, [roomCode]);

  const handleDialogToggle = (isOpen: boolean) => {
    if (!isOpen) {
      setRoomCode("");
      setRoom(null);
      setSearched(false);
      setLoading(false);
    }
    setOpen(isOpen);
  };

  const handleRoomJoinClick = async () => {
    if (room) {
      setLoading(true);
      await joinPrivateRoom(roomCode);
      setLoading(false);
      setOpen(false);
      navigate(`/chat/${room.name}`, { replace: true });
    }
  };

  const handleRoomSearchClick = async () => {
    setLoading(true);
    const data = await getPrivateChatRoom(roomCode);
    setRoom(data || null);
    setSearched(true);
    setLoading(false);
  };

  const handleLeaveClick = (roomName: string) => {
    setLeaveRoomName(roomName);
  };

  const confirmLeaveRoom = async () => {
    if (leaveRoomName) {
      const leaved = await leaveRoom(leaveRoomName);
      setLeaveRoomName(null);
      if (leaved) {
        setRoom((prev) =>
          prev ? { ...prev, isMember: false, members: prev.members - 1 } : null,
        );
      }
    }
  };

  if (!open) return null;

  return (
    <>
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
      <Dialog open={open} onOpenChange={handleDialogToggle}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Join Private Room</DialogTitle>
            <DialogDescription>Join Private Room Using Code</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-[auto_1fr_auto] items-center gap-3">
              <Label htmlFor="roomCode" className="text-right">
                Room Code
              </Label>
              <Input
                id="roomCode"
                value={roomCode}
                onChange={(e) => setRoomCode(e.target.value)}
                className="w-full"
                placeholder="Enter Code"
              />

              <Button
                size="icon"
                onClick={handleRoomSearchClick}
                disabled={loading || disable}
                className="flex justify-center items-center"
              >
                {loading && <Loader2 className="animate-spin" />}
                {!loading && <Search />}
              </Button>
            </div>
            {room && user ? (
              <Card className="border-2 border-muted">
                <CardHeader className="pb-2">
                  <div className="flex justify-between items-center">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <Hash className="h-5 w-5 text-primary" />
                        {room.name}
                      </CardTitle>
                      <div className="flex gap-2 mt-2">
                        {/* Owner Tag */}
                        {room.userId === user.id && (
                          <Badge
                            variant="secondary"
                            className="text-xs cursor-pointer"
                          >
                            You are the owner
                          </Badge>
                        )}

                        {/* Privacy Tag */}
                        <Badge
                          variant="outline"
                          className="text-xs flex items-center gap-1 cursor-pointer"
                        >
                          <Lock className="h-3 w-3" />
                          Private
                        </Badge>
                      </div>
                    </div>

                    <Badge variant="outline" className="gap-1">
                      <Users className="h-3 w-3" />
                      <span>{room.members} members</span>
                    </Badge>
                  </div>
                </CardHeader>

                <CardContent>
                  <p className="text-muted-foreground">{room.description}</p>
                </CardContent>

                <CardFooter className="flex gap-2">
                  {room.isMember ? (
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
                    <Button className="w-full" onClick={handleRoomJoinClick}>
                      Join Room
                    </Button>
                  )}
                </CardFooter>
              </Card>
            ) : (
              searched && (
                <div className="w-full py-10">
                  <p className="text-sm text-center text-gray-200">
                    No Chat Room Found
                  </p>
                </div>
              )
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
};

export default PrivateRoomJoinModal;
