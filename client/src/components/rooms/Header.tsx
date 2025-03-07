import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";

const Header = () => {
  const navigate = useNavigate();
  const handleBackToChat = () => {
    navigate("/chat");
  };

  return (
    <header className="px-10 sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="w-full flex h-16 items-center">
        <div className="flex items-center gap-2 font-bold text-xl">
          <img className="w-7 h-7" src="/icon.png" alt="logo" />
          <span>Sync Talk</span>
        </div>
        <div className="ml-auto flex items-center gap-4">
          <Button
            variant="outline"
            onClick={handleBackToChat}
            className="gap-2"
          >
            <ArrowLeft className="h-4 w-4" />
            Go to Chat
          </Button>
        </div>
      </div>
    </header>
  );
};

export default Header;
