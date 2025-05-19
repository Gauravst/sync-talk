import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Hash, Lock, Users, Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";

import { UserProps } from "@/types";
import { ChatRoomProps } from "@/types/messageTypes";

type ChatAreaHeaderProps = {
  user: UserProps;
  roomData: ChatRoomProps;
  name: string;
  onlineUsers: number;
};

export const ChatAreaHeader = ({
  user,
  roomData,
  name,
  onlineUsers,
}: ChatAreaHeaderProps) => {
  const [copied, setCopied] = useState(false);
  const handleCopy = () => {
    navigator.clipboard.writeText(roomData.code!);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <div className="flex items-center justify-between border-b px-4 py-3 h-[60px]">
      <div className="flex items-center gap-2">
        <Hash className="h-5 w-5 text-muted-foreground" />
        <h3 className="font-semibold">{name}</h3>
        {roomData && roomData.userId == user?.id && (
          <div className="flex gap-x-2 mx-2">
            <Badge variant="outline" className="gap-1">
              <Users className="h-3 w-3" />
              <span>You are the owner</span>
            </Badge>
            <Badge variant="outline" className="gap-1">
              <Lock className="h-3 w-3" />
              Private
            </Badge>
            <Button
              onClick={handleCopy}
              title="copy code"
              className="flex items-center gap-1 text-sm w-6 h-6"
            >
              {copied ? (
                <Check size={10} className="text-green-500" />
              ) : (
                <Copy size={10} />
              )}
            </Button>
          </div>
        )}
      </div>
      <div className="flex items-center gap-2">
        <Badge variant="outline" className="gap-1">
          <Users className="h-3 w-3" />
          <span>{onlineUsers} online</span>
        </Badge>
      </div>
    </div>
  );
};
