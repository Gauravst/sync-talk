import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login, loginWithoutAuth } from "@/services/authService";
import { useAuth } from "@/context/AuthContext";

function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const { user, loading } = useAuth();
  const [loginLoading, setLoginLoading] = useState(false);
  const [login2Loading, setLogin2Loading] = useState(false);

  const handleLogin = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoginLoading(true);
    try {
      const loginStatus = await login(username, password);
      if (loginStatus) {
        navigate(`/chat`, { replace: true });
        window.location.reload();
      }
    } catch (error) {
      console.log("Invalid login credentials");
      console.log(error);
      setLoginLoading(false);
    }
  };

  const handleLogin2 = async () => {
    setLogin2Loading(true);
    try {
      const loginStatus = await loginWithoutAuth();
      if (loginStatus) {
        navigate(`/rooms`, { replace: true });
        window.location.reload();
      }
    } catch (error) {
      setLogin2Loading(false);
      console.error("Invalid login credentials", error);
    }
  };

  useEffect(() => {
    if (user) {
      navigate("/chat", { replace: true });
    }
  }, [user, navigate]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>Welcome To Sync Talk</CardTitle>
          <CardDescription>Create or Login Your account</CardDescription>
        </CardHeader>

        {/* Form submits with handleLogin */}
        <form onSubmit={handleLogin}>
          <CardContent>
            <div className="grid w-full items-center gap-4">
              {/* Username Input */}
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  placeholder="Enter your username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>
              {/* Password Input */}
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Enter your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </div>
            </div>
          </CardContent>
          <CardFooter className="flex flex-col gap-2 pb-2">
            {/* Submit button - triggers handleLogin */}
            <Button disabled={loginLoading} type="submit" className="w-full">
              {loginLoading && <Loader2 className="animate-spin mr-2" />}
              {loginLoading ? "Please wait" : "Continue"}
            </Button>
          </CardFooter>
        </form>

        <CardFooter className="flex flex-col">
          <Button
            disabled={login2Loading}
            onClick={handleLogin2}
            variant="outline"
            className="w-full"
          >
            {login2Loading && <Loader2 className="animate-spin mr-2" />}
            {login2Loading ? "Please wait" : "Continue Without Auth"}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}

export default LoginPage;
