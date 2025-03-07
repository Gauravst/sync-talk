import { Card, CardContent } from "@/components/ui/card";
import { Zap, Shield, Globe } from "lucide-react";

const FeaturesSection = () => {
  return (
    <section className="container py-12 space-y-8">
      <div className="flex flex-col items-center text-center space-y-4">
        <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
          Powerful Features
        </h2>
        <p className="max-w-[700px] text-muted-foreground">
          Discover what makes Sync Talk the best choice for real-time
          communication.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        <Card className="bg-background/50">
          <CardContent className="pt-6">
            <div className="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
              <Zap className="h-5 w-5 text-primary" />
            </div>
            <h3 className="text-xl font-bold">Lightning Fast</h3>
            <p className="text-muted-foreground mt-2">
              Powered by Go and WebSockets for instant message delivery with
              minimal latency.
            </p>
          </CardContent>
        </Card>

        <Card className="bg-background/50">
          <CardContent className="pt-6">
            <div className="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
              <Shield className="h-5 w-5 text-primary" />
            </div>
            <h3 className="text-xl font-bold">Secure</h3>
            <p className="text-muted-foreground mt-2">
              End-to-end encryption and robust security measures to protect your
              conversations.
            </p>
          </CardContent>
        </Card>

        <Card className="bg-background/50">
          <CardContent className="pt-6">
            <div className="mb-4 flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
              <Globe className="h-5 w-5 text-primary" />
            </div>
            <h3 className="text-xl font-bold">Scalable</h3>
            <p className="text-muted-foreground mt-2">
              Built to handle millions of concurrent connections without
              compromising performance.
            </p>
          </CardContent>
        </Card>
      </div>
    </section>
  );
};

export default FeaturesSection;
