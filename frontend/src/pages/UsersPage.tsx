import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  getUsers,
  createUser,
  updateUser,
  updateUserPassword,
  deleteUser,
} from "@/services/user.service";
import type { UserRecord } from "@/types/user";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import TablePagination from "@/components/TablePagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Plus, Pencil, Trash2, Lock, Clock } from "lucide-react";

const PAGE_SIZE = 10;

export default function UsersPage() {
  const qc = useQueryClient();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [modalOpen, setModalOpen] = useState(false);
  const [pwModalOpen, setPwModalOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<UserRecord | null>(null);
  const [editing, setEditing] = useState<UserRecord | null>(null);
  const [pwTarget, setPwTarget] = useState<UserRecord | null>(null);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("cs");
  const [newPassword, setNewPassword] = useState("");

  const {
    data,
    isLoading,
    error: queryError,
  } = useQuery({
    queryKey: ["users", search, page],
    queryFn: () =>
      getUsers({ search: search || undefined, page, limit: PAGE_SIZE }),
  });
  const users: UserRecord[] = data?.users ?? [];
  const totalCount = data?.total ?? 0;
  const totalPages = data?.total_pages ?? 1;
  const currentPage = page;

  const createMut = useMutation({
    mutationFn: createUser,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
      closeModal();
      toast.success("User created successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to create user"),
  });

  const updateMut = useMutation({
    mutationFn: ({
      id,
      email,
      role,
    }: {
      id: number;
      email: string;
      role: string;
    }) => updateUser(id, { email, role }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
      closeModal();
      toast.success("User updated successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to update user"),
  });

  const pwMut = useMutation({
    mutationFn: ({ id, password }: { id: number; password: string }) =>
      updateUserPassword(id, password),
    onSuccess: () => {
      closePwModal();
      toast.success("Password updated successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to update password"),
  });

  const deleteMut = useMutation({
    mutationFn: (id: number) => deleteUser(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
      toast.success("User deleted");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to delete user"),
  });

  function openCreate() {
    setEditing(null);
    setEmail("");
    setPassword("");
    setRole("cs");
    setModalOpen(true);
  }
  function openEdit(u: UserRecord) {
    setEditing(u);
    setEmail(u.email);
    setRole(u.role);
    setModalOpen(true);
  }
  function openPwModal(u: UserRecord) {
    setPwTarget(u);
    setNewPassword("");
    setPwModalOpen(true);
  }
  function closeModal() {
    setModalOpen(false);
    setEditing(null);
  }
  function closePwModal() {
    setPwModalOpen(false);
    setPwTarget(null);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (editing) {
      updateMut.mutate({ id: Number(editing.id), email, role });
    } else {
      createMut.mutate({ email, password, role });
    }
  }

  function handlePwSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (pwTarget)
      pwMut.mutate({ id: Number(pwTarget.id), password: newPassword });
  }

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-bold text-foreground">Users</h1>
        <div className="flex gap-3 items-center">
          <Button size="sm" onClick={openCreate}>
            <Plus className="size-4" /> Add User
          </Button>
          <Input
            placeholder="Search..."
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1);
            }}
            className="w-48"
          />
        </div>
      </div>

      <div className="rounded-xl border border-border bg-card overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center gap-3 p-12 text-muted-foreground text-sm">
            <Clock className="size-4 animate-spin" /> Loading users...
          </div>
        ) : queryError ? (
          <div className="p-12 text-center text-destructive text-sm">
            Failed to load users.
          </div>
        ) : totalCount === 0 ? (
          <div className="p-12 text-center text-muted-foreground text-sm">
            No users found.
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.map((u) => (
                  <TableRow key={u.id}>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      {u.id}
                    </TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          u.role === "operation" ? "default" : "secondary"
                        }
                      >
                        {u.role}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => openEdit(u)}
                          title="Edit"
                        >
                          <Pencil className="size-3.5" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => openPwModal(u)}
                          title="Change Password"
                        >
                          <Lock className="size-3.5" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => setDeleteTarget(u)}
                          title="Delete"
                          className="hover:text-destructive"
                        >
                          <Trash2 className="size-3.5" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            <TablePagination
              page={currentPage}
              totalPages={totalPages}
              total={totalCount}
              limit={PAGE_SIZE}
              onPageChange={setPage}
            />
          </>
        )}
      </div>

      {/* Create / Edit User */}
      <Dialog
        open={modalOpen}
        onOpenChange={(open) => {
          if (!open) closeModal();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? "Edit User" : "Create User"}</DialogTitle>
          </DialogHeader>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label>Email</Label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="user@example.com"
                required
                autoFocus
              />
            </div>
            {!editing && (
              <div className="space-y-2">
                <Label>Password</Label>
                <Input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter password"
                  required
                />
              </div>
            )}
            <div className="space-y-2">
              <Label>Role</Label>
              <Select value={role} onValueChange={setRole}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="cs">CS</SelectItem>
                  <SelectItem value="operation">Operation</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={closeModal}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={createMut.isPending || updateMut.isPending}
              >
                {editing ? "Update" : "Create"}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Change Password */}
      <Dialog
        open={pwModalOpen}
        onOpenChange={(open) => {
          if (!open) closePwModal();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change Password</DialogTitle>
          </DialogHeader>
          <form className="space-y-4" onSubmit={handlePwSubmit}>
            <div className="space-y-2">
              <Label>User</Label>
              <Input value={pwTarget?.email ?? ""} disabled />
            </div>
            <div className="space-y-2">
              <Label>New Password</Label>
              <Input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="Enter new password"
                required
                autoFocus
              />
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={closePwModal}>
                Cancel
              </Button>
              <Button type="submit" disabled={pwMut.isPending}>
                Update Password
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <AlertDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete User</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete{" "}
              <span className="font-semibold">{deleteTarget?.email}</span>? This
              action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (deleteTarget) {
                  deleteMut.mutate(Number(deleteTarget.id));
                  setDeleteTarget(null);
                }
              }}
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
