import { Button } from "@/components/ui/button";
import { MessageSquare } from "lucide-react";

const CTASection = () => {
  return (
    <section className="container py-12 md:py-24">
      <div className="relative overflow-hidden rounded-lg border bg-background p-8">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-500/10 to-purple-500/10" />
        <div className="relative flex flex-col items-center text-center space-y-4">
          <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">
            Ready to start chatting?
          </h2>
          <p className="max-w-[700px] text-muted-foreground md:text-xl">
            Join thousands of users already enjoying Sync Talk's seamless
            experience.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 mt-4">
            <Button size="lg" className="gap-2">
              <MessageSquare className="h-5 w-5" />
              Get Started Free
            </Button>
            <Button size="lg" variant="outline" className="gap-2">
              Learn More
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
};

export default CTASection;
