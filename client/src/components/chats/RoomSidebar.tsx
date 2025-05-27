import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChatRoomProps } from "@/types/messageTypes";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { Hash, Plus, Search } from "lucide-react";

import { useNavigate } from "react-router-dom";
import { useSocket } from "@/hooks/useSocket";
import { useAuth } from "@/context/AuthContext";

type RoomSidebarProps = {
  chatGroups: ChatRoomProps[];
  setNewRoomPopup: (value: boolean) => void;
  name: string;
  handleGroupClick: (groupName: string, data: ChatRoomProps) => void;
};

const RoomSidebar = ({
  chatGroups,
  setNewRoomPopup,
  name,
  handleGroupClick,
}: RoomSidebarProps) => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { closeSocket } = useSocket(null);

  const handleFindNewRooms = async () => {
    await closeSocket();
    navigate("/rooms");
  };
  return (
    <div className="lg:col-span-1">
      <Card className="h-full border-2 border-muted">
        <CardContent className="p-0">
          <div className="p-4 border-b">
            <div className="flex items-center justify-between mb-4">
              <h3 className="font-semibold text-lg">Chat Rooms</h3>
              <Button
                onClick={() => setNewRoomPopup(true)}
                variant="ghost"
                size="icon"
              >
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
                    onClick={() => handleGroupClick(group?.name, group)}
                  >
                    <Hash className="h-4 w-4" />
                    {group?.name}

                    {group.userId == user?.id && (
                      <div className="flex gap-x-2 mx-2">
                        <Badge
                          variant="outline"
                          className="gap-1 bg-background"
                        >
                          <span>You are the owner</span>
                        </Badge>
                        <Badge
                          variant="outline"
                          className="gap-1 bg-background"
                        >
                          Private
                        </Badge>
                      </div>
                    )}
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
  );
};

export default RoomSidebar;
