import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { MessageSquare, Users, Zap } from "lucide-react";

const HeroSection = () => {
  return (
    <section className="container py-24 space-y-8 md:py-32">
      <div className="flex flex-col items-center text-center space-y-4">
        <Badge variant="outline" className="px-3 py-1 text-sm">
          <Zap className="mr-1 h-3 w-3" />
          <span>Now in Beta</span>
        </Badge>
        <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl lg:text-7xl">
          Real-time conversations,
          <br />
          <span className="bg-gradient-to-r from-blue-500 to-purple-500 text-transparent bg-clip-text">
            seamlessly synced.
          </span>
        </h1>
        <p className="max-w-[700px] text-muted-foreground md:text-xl">
          Experience lightning-fast messaging powered by WebSockets and Go.
          Connect with anyone, anywhere, instantly.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 mt-8">
          <Button size="lg" className="gap-2">
            <MessageSquare className="h-5 w-5" />
            Start Chatting
          </Button>
          <Button size="lg" variant="outline" className="gap-2">
            <Users className="h-5 w-5" />
            Join Community
          </Button>
        </div>
      </div>
    </section>
  );
};

export default HeroSection;
