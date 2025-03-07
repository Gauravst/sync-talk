import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Send, Users, Hash, Plus, Search } from "lucide-react";
import { useState } from "react";

const ChatPreviewSection = () => {
  const [message, setMessage] = useState("");
  return (
    <section className="container py-12 space-y-8">
      <div className="flex flex-col items-center text-center space-y-4">
        <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
          Join Multiple Chat Rooms
        </h2>
        <p className="max-w-[700px] text-muted-foreground">
          Connect with different communities in topic-specific rooms. Find your
          people, join the conversation.
        </p>
      </div>

      <div className="mx-auto max-w-6xl w-full p-4">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Sidebar with Rooms */}
          <div className="lg:col-span-1">
            <Card className="h-full border-2 border-muted">
              <CardContent className="p-0">
                <div className="p-4 border-b">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="font-semibold text-lg">Chat Rooms</h3>
                    <Button variant="ghost" size="icon">
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                  <div className="relative">
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input placeholder="Search rooms" className="pl-8" />
                  </div>
                </div>
                <ScrollArea className="[400px]">
                  <div className="p-2">
                    <div className="space-y-1">
                      <Button
                        variant="secondary"
                        className="w-full justify-start gap-2 font-medium"
                      >
                        <Hash className="h-4 w-4" />
                        general
                        <Badge className="ml-auto">24</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        tech-talk
                        <Badge className="ml-auto">8</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        golang
                        <Badge className="ml-auto">12</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        react
                        <Badge className="ml-auto">15</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        websockets
                        <Badge className="ml-auto">6</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        gaming
                        <Badge className="ml-auto">32</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        music
                        <Badge className="ml-auto">18</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        movies
                        <Badge className="ml-auto">9</Badge>
                      </Button>
                      <Button
                        variant="ghost"
                        className="w-full justify-start gap-2"
                      >
                        <Hash className="h-4 w-4" />
                        books
                        <Badge className="ml-auto">5</Badge>
                      </Button>
                    </div>
                  </div>
                </ScrollArea>
                <div className="p-4 border-t">
                  <Button className="w-full gap-2">
                    <Plus className="h-4 w-4" />
                    Create New Room
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Chat Area */}
          <div className="lg:col-span-2">
            <Card className="h-full border-2 border-muted">
              <CardContent className="p-0">
                <div className="flex items-center justify-between border-b px-4 py-3">
                  <div className="flex items-center gap-2">
                    <Hash className="h-5 w-5 text-muted-foreground" />
                    <h3 className="font-semibold">general</h3>
                  </div>
                  <div className="flex items-center gap-2">
                    <Badge variant="outline" className="gap-1">
                      <Users className="h-3 w-3" />
                      <span>24 online</span>
                    </Badge>
                  </div>
                </div>

                <div className="flex h-[400px] justify-between flex-col">
                  <ScrollArea className="flex-1 h-[200px] p-4">
                    <div className="space-y-4 h-full">
                      {/* Message 1 */}
                      <div className="flex items-start gap-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=96&h=96&fit=crop&auto=format" />
                          <AvatarFallback>JD</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="flex items-center gap-2">
                            <div className="font-semibold">John Doe</div>
                            <div className="text-xs text-muted-foreground">
                              10:23 AM
                            </div>
                          </div>
                          <div className="mt-1 rounded-lg bg-muted p-3">
                            <p>
                              Hey everyone! Has anyone tried the new WebSocket
                              implementation?
                            </p>
                          </div>
                        </div>
                      </div>

                      {/* Message 2 */}
                      <div className="flex items-start gap-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=96&h=96&fit=crop&auto=format" />
                          <AvatarFallback>AS</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="flex items-center gap-2">
                            <div className="font-semibold">Alice Smith</div>
                            <div className="text-xs text-muted-foreground">
                              10:25 AM
                            </div>
                          </div>
                          <div className="mt-1 rounded-lg bg-muted p-3">
                            <p>
                              Yes! The performance is incredible. Messages sync
                              almost instantly.
                            </p>
                          </div>
                        </div>
                      </div>

                      {/* Message 3 */}
                      <div className="flex items-start gap-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1599566150163-29194dcaad36?w=96&h=96&fit=crop&auto=format" />
                          <AvatarFallback>RJ</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="flex items-center gap-2">
                            <div className="font-semibold">Robert Johnson</div>
                            <div className="text-xs text-muted-foreground">
                              10:27 AM
                            </div>
                          </div>
                          <div className="mt-1 rounded-lg bg-muted p-3">
                            <p>
                              I'm impressed with how well it handles multiple
                              users. No lag at all!
                            </p>
                          </div>
                        </div>
                      </div>

                      {/* Message 1 */}
                      <div className="flex items-start gap-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=96&h=96&fit=crop&auto=format" />
                          <AvatarFallback>JD</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="flex items-center gap-2">
                            <div className="font-semibold">John Doe</div>
                            <div className="text-xs text-muted-foreground">
                              10:23 AM
                            </div>
                          </div>
                          <div className="mt-1 rounded-lg bg-muted p-3">
                            <p>Yaah, Go is cool</p>
                          </div>
                        </div>
                      </div>

                      {/* Message 4 */}
                      <div className="flex items-start gap-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1607746882042-944635dfe10e?w=96&h=96&fit=crop&auto=format" />
                          <AvatarFallback>EM</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="flex items-center gap-2">
                            <div className="font-semibold">Emily Martinez</div>
                            <div className="text-xs text-muted-foreground">
                              10:30 AM
                            </div>
                          </div>
                          <div className="mt-1 rounded-lg bg-muted p-3">
                            <p>The UI is so clean too. Love the dark mode!</p>
                          </div>
                        </div>
                      </div>

                      {/* Message 5 */}
                      <div className="flex items-start gap-3">
                        <Avatar>
                          <AvatarImage src="https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=96&h=96&fit=crop&auto=format" />
                          <AvatarFallback>JD</AvatarFallback>
                        </Avatar>
                        <div>
                          <div className="flex items-center gap-2">
                            <div className="font-semibold">John Doe</div>
                            <div className="text-xs text-muted-foreground">
                              10:32 AM
                            </div>
                          </div>
                          <div className="mt-1 rounded-lg bg-muted p-3">
                            <p>
                              Agreed! The Golang backend makes everything super
                              fast.
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  </ScrollArea>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </section>
  );
};

export default ChatPreviewSection;
