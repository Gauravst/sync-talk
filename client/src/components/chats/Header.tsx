import { Settings } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "@/context/AuthContext";
import { Link, useNavigate } from "react-router-dom";
import { useSocket } from "@/hooks/useSocket";

type HeaderProps = {
  handleProfileClick: () => void;
};

const Header = ({ handleProfileClick }: HeaderProps) => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const { closeSocket } = useSocket(null);

  const handleLogout = async () => {
    closeSocket();
    await logout();
    navigate("/login");
  };

  return (
    <>
      <header className="px-5 sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="w-full flex h-14 items-center">
          <div className="w-full flex items-center font-bold text-xl">
            <Link to="/" className="flex gap-2 justify-center items-center">
              <img className="w-7 h-7" src="/icon.png" alt="logo" />
              <span>Sync Talk</span>
            </Link>
          </div>
          <div className="ml-auto flex items-center gap-4">
            <DropdownMenu>
              <DropdownMenuTrigger>
                <div className="p-2 bg-gray-900 rounded-md">
                  <Settings className="h-5 w-5" />
                </div>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="mr-3">
                <DropdownMenuItem
                  className="cursor-pointer"
                  onClick={handleProfileClick}
                >
                  Profile
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  className="cursor-pointer"
                  onClick={handleLogout}
                >
                  Logout
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>
    </>
  );
};

export default Header;
