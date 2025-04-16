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
import { useRef, useState } from "react";
import { UserProps } from "@/types";
import { Pencil } from "lucide-react";

type ProfileProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
  userData: UserProps;
};

const ProfileDialog = ({ open, setOpen, userData }: ProfileProps) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [image, setImage] = useState<string | null>(null);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      setImage(imageUrl);
    }
  };
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Profile</DialogTitle>
          <DialogDescription>Your Profile Username and DP</DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="relative w-fit gap-4">
            <Avatar className="w-20 h-20">
              <AvatarImage src={image || "https://github.com/gauravst.png"} />
              <AvatarFallback>GS</AvatarFallback>
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
            <Input id="username" value="@peduarte" className="col-span-3" />
          </div>
        </div>
        <DialogFooter>
          <Button type="submit">Save changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default ProfileDialog;
