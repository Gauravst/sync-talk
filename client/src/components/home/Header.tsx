import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

const Header = () => {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-16 mx-8 items-center">
        <div className="flex items-center gap-2 font-bold text-xl">
          <img className="w-7 h-7" src="/icon.png" alt="logo" />
          <span>Sync Talk</span>
        </div>
        <div className="ml-auto flex items-center gap-4">
          <Link to="/features">
            <Button variant="ghost">Features</Button>
          </Link>
          <Link to="/pricing">
            <Button variant="ghost">Pricing</Button>
          </Link>
          <Link to="/about">
            <Button variant="ghost">About</Button>
          </Link>
          <Link to="/login">
            <Button variant="default">Get Started</Button>
          </Link>
        </div>
      </div>
    </header>
  );
};

export default Header;
