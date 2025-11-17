import { Button } from "@/src/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/src/components/ui/dialog";
import { Input } from "../../ui/input";
import Link from "next/link";

export type NeedLoginButtonViewProps = {
  isLoggedIn: boolean;
  content: string;
  loggedInLink?: string;
  className?: string;
  isDialogOpen: boolean;
  userName: string;
  isLoading: boolean;
  error: string | null;
  onButtonClick: () => void;
  onDialogOpenChange: (open: boolean) => void;
  onUserNameChange: (userName: string) => void;
  onSubmit: () => void;
};

export function NeedLoginButtonView({
  isLoggedIn,
  content,
  loggedInLink,
  className,
  isDialogOpen,
  userName,
  isLoading,
  error,
  onButtonClick,
  onDialogOpenChange,
  onUserNameChange,
  onSubmit,
}: NeedLoginButtonViewProps) {
  return (
    <>
      <Button onClick={onButtonClick} className={className}>
        {loggedInLink && isLoggedIn ? (
          <Link href={loggedInLink}>{content}</Link>
        ) : (
          content
        )}
      </Button>

      <Dialog open={isDialogOpen} onOpenChange={onDialogOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>ユーザー作成</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Input
              type="text"
              placeholder="ユーザー名"
              value={userName}
              onChange={(e) => onUserNameChange(e.target.value)}
              disabled={isLoading}
              className="bg-white"
            />
            {error && <p className="text-sm text-red-500">{error}</p>}
            <Button
              type="submit"
              onClick={onSubmit}
              disabled={isLoading}
              className="bg-button-accent hover:bg-button-accent/80 w-full"
            >
              {isLoading ? "作成中..." : "作成"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
