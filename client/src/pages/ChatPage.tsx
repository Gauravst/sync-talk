import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getJoinedRoom } from "@/services/chatServices";
import Header from "@/components/chats/Header";
import { useAuth } from "@/context/AuthContext";
import { ChatRoomProps } from "@/types/messageTypes";
import ProfileDialog from "@/components/chats/ProfileSection";
import CreateNewRoom from "@/components/chats/CreateNewRoom";
import RoomSidebar from "@/components/chats/RoomSidebar";
import ChatArea from "@/components/chats/ChatArea";

function ChatPage() {
  const { user } = useAuth();
  const { name } = useParams();
  const [chatGroups, setChatGroups] = useState<ChatRoomProps[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isJoined, setIsJoined] = useState<boolean>(false);
  const navigate = useNavigate();
  const [profilePopup, setProfilePopup] = useState(false);
  const [newRoomPopup, setNewRoomPopup] = useState(false);
  const [currentOpenRoom, setCurrentOpenRoom] = useState<ChatRoomProps>();
  console.log(loading);

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
          setChatGroups([]);
        }
      } catch (error) {
        console.error("Failed to load chat rooms", error);
      } finally {
        setLoading(false);
      }
    };
    fetchRooms();
  }, [name, navigate]);

  const handleGroupClick = (groupName: string, groupData: ChatRoomProps) => {
    setCurrentOpenRoom(groupData);
    navigate(`/chat/${groupName}`, { replace: true });
  };

  const updateRoomsList = (room: ChatRoomProps) => {
    setChatGroups((prev) => [room, ...prev]);
  };

  return (
    <div className="flex relative flex-col min-h-screen bg-background">
      {/* Header */}
      <Header handleProfileClick={() => setProfilePopup(true)} />
      <ProfileDialog
        open={profilePopup}
        setOpen={setProfilePopup}
        userData={user!}
      />

      <CreateNewRoom
        open={newRoomPopup}
        setOpen={setNewRoomPopup}
        updateRoomsList={updateRoomsList}
      />

      {/* Main Content */}
      <div className="flex-1">
        <div className="p-3 fixed w-full grid grid-cols-1 md:grid-cols-3 gap-3">
          {/* Sidebar with Rooms */}
          <RoomSidebar
            chatGroups={chatGroups}
            setNewRoomPopup={setNewRoomPopup}
            name={name!}
            handleGroupClick={handleGroupClick}
          />

          {/* Chat Area */}
          <ChatArea
            name={name!}
            roomData={currentOpenRoom!}
            isJoined={isJoined}
            setIsJoined={setIsJoined}
          />
        </div>
      </div>
    </div>
  );
}

export default ChatPage;
