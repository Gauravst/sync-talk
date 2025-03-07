import { Button } from "@/components/ui/button";
import { Settings } from "lucide-react";

const Header = () => {
  return (
    <header className="px-5 sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="w-full flex h-14 items-center">
        <div className="w-full flex items-center gap-2 font-bold text-xl">
          <img className="w-7 h-7" src="/icon.png" alt="logo" />
          <span>Sync Talk</span>
        </div>
        <div className="ml-auto flex items-center gap-4">
          <Button variant="ghost" size="icon">
            <Settings className="h-5 w-5" />
          </Button>
        </div>
      </div>
    </header>
  );
};

export default Header;
