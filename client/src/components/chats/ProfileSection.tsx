import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useRef, useState, useEffect } from "react";
import { UserProps } from "@/types";
import { Pencil, Loader2 } from "lucide-react";
import { updateUser } from "@/services/authService";

type ProfileProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
  userData: UserProps;
};

const ProfileDialog = ({ open, setOpen, userData }: ProfileProps) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [username, setUsername] = useState(userData.username);
  const [password, setPassword] = useState("");
  const [image, setImage] = useState<string | null>(
    userData?.profilePic || null,
  );
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (open) {
      setUsername(userData.username);
      setPassword("");
      setImage(userData.profilePic || null);
      setImageFile(null);
      setLoading(false);
    }
  }, [open, userData]);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      setImage(imageUrl);
      setImageFile(file);
    }
  };

  const handleCloseDialog = (isOpen: boolean) => {
    if (!isOpen) {
      setUsername(userData.username);
      setPassword("");
      setImage(userData.profilePic || null);
      setImageFile(null);
      setLoading(false);
    }
    setOpen(isOpen);
  };

  const handleSaveChanges = async () => {
    const formData = new FormData();
    const trimmedUsername = username.trim();

    if (
      trimmedUsername === userData.username &&
      password === "" &&
      !imageFile
    ) {
      return;
    }

    setLoading(true);

    if (trimmedUsername !== userData.username) {
      formData.append("username", trimmedUsername);
    }

    if (password) {
      formData.append("password", password);
    }

    if (imageFile) {
      formData.append("profilePic", imageFile);
    }

    try {
      await updateUser(formData);
      setOpen(false);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleCloseDialog}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Profile</DialogTitle>
          <DialogDescription>Your Profile Username and DP</DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="relative w-fit gap-4 ml-4">
            <Avatar className="w-20 h-20">
              <AvatarImage src={image!} />
              <AvatarFallback>ST</AvatarFallback>
            </Avatar>

            <input
              type="file"
              accept="image/*"
              ref={fileInputRef}
              onChange={handleImageChange}
              className="hidden"
            />

            <Button
              variant="ghost"
              size="icon"
              className="absolute bottom-0 right-0 bg-background shadow-md"
              onClick={() => fileInputRef.current?.click()}
            >
              <Pencil className="w-4 h-4" />
            </Button>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="username" className="text-right">
              Username
            </Label>
            <Input
              id="username"
              onChange={(e) => setUsername(e.target.value)}
              value={username}
              className="col-span-3"
              placeholder="Your Username"
            />
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="password" className="text-right">
              Password
            </Label>
            <Input
              id="password"
              onChange={(e) => setPassword(e.target.value)}
              value={password}
              className="col-span-3"
              placeholder="Your New Password"
              type="password"
            />
          </div>
        </div>
        <DialogFooter>
          <Button onClick={handleSaveChanges} disabled={loading}>
            {loading && <Loader2 className="animate-spin mr-2" />}
            {loading ? "Saving..." : "Save changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default ProfileDialog;
