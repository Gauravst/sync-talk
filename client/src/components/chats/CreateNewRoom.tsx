import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { useState, useEffect } from "react";
import { Loader2 } from "lucide-react";
import { createNewRoom } from "@/services/chatServices";

type ProfileProps = {
  open: boolean;
  setOpen: (open: boolean) => void;
};

const CreateNewRoom = ({ open, setOpen }: ProfileProps) => {
  const [username, setUsername] = useState("");
  const [description, setDescription] = useState("");
  const [loading, setLoading] = useState(false);
  const [disable, setDisable] = useState(true);

  useEffect(() => {
    if (username.trim() && description.trim()) {
      setDisable(false);
    } else {
      setDisable(true);
    }
  }, [username, description]);

  const handleDialogToggle = (isOpen: boolean) => {
    if (!isOpen) {
      setUsername("");
      setDescription("");
      setLoading(false);
    }
    setOpen(isOpen);
  };

  const handleCreateRoomClick = async () => {
    setLoading(true);
    await createNewRoom(username, description);
    setLoading(false);
    setOpen(false);
    setUsername("");
    setDescription("");
  };

  if (!open) return null;

  return (
    <Dialog open={open} onOpenChange={handleDialogToggle}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Room</DialogTitle>
          <DialogDescription>
            Create New Personal Private Room
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="username" className="text-right">
              Username
            </Label>
            <Input
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="col-span-3"
              placeholder="Room's Username"
            />
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="description" className="text-right">
              Description
            </Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="col-span-3 resize-none"
              placeholder="Room's Description"
            />
          </div>
        </div>
        <DialogFooter>
          <Button onClick={handleCreateRoomClick} disabled={loading || disable}>
            {loading && <Loader2 className="animate-spin mr-2" />}
            {loading ? "Please wait" : "Create Private Room"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default CreateNewRoom;
