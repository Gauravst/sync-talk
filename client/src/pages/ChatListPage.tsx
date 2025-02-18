import React from "react";
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
import { Link } from "react-router-dom"; // Assuming you're using React Router for navigation

const chatGroups = [
  {
    id: 1,
    name: "Family",
    username: "@family",
    members: 5,
    image: "/family-group.jpg",
  },
  {
    id: 2,
    name: "Work Team",
    username: "@workteam",
    members: 8,
    image: "/work-team.jpg",
  },
  {
    id: 3,
    name: "Friends",
    username: "@friends",
    members: 12,
    image: "/friends-group.jpg",
  },
  // Add more chat groups as needed
];

function RoomListPage() {
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Your Chat Groups</h1>
      <ScrollArea className="h-[calc(100vh-100px)]">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {chatGroups.map((group) => (
            <Card key={group.id}>
              <CardHeader>
                <div className="flex items-center space-x-4">
                  <Avatar>
                    <AvatarImage src={group.image} alt={group.name} />
                    <AvatarFallback>{group.name.slice(0, 2)}</AvatarFallback>
                  </Avatar>
                  <div>
                    <CardTitle>{group.name}</CardTitle>
                    <CardDescription>{group.username}</CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <p>{group.members} members</p>
              </CardContent>
              <CardFooter>
                <Button asChild className="w-full">
                  <Link to={`/chat/${group.id}`}>Join Chat</Link>
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}

export default RoomListPage;
