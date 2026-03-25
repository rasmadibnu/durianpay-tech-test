import { useState, useMemo } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  getMerchants,
  createMerchant,
  updateMerchant,
  deleteMerchant,
} from "@/services/merchant.service";
import type { Merchant } from "@/types/merchant";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
import { Plus, Pencil, Trash2, Clock } from "lucide-react";
import { formatDate } from "@/lib/utils";
import TablePagination from "@/components/TablePagination";

const PAGE_SIZE = 10;

export default function MerchantsPage() {
  const qc = useQueryClient();
  const [page, setPage] = useState(1);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Merchant | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Merchant | null>(null);
  const [name, setName] = useState("");
  const [search, setSearch] = useState("");

  const {
    data,
    isLoading,
    error: queryError,
  } = useQuery({
    queryKey: ["merchants", search, page],
    queryFn: () =>
      getMerchants({ search: search || undefined, page, limit: PAGE_SIZE }),
  });
  const merchants: Merchant[] = data?.merchants ?? [];
  const totalCount = data?.total ?? 0;
  const totalPages = data?.total_pages ?? 1;
  const currentPage = page;

  const createMut = useMutation({
    mutationFn: createMerchant,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["merchants"] });
      closeModal();
      toast.success("Merchant created successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to create merchant"),
  });

  const updateMut = useMutation({
    mutationFn: ({ id, name }: { id: number; name: string }) =>
      updateMerchant(id, { name }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["merchants"] });
      closeModal();
      toast.success("Merchant updated successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to update merchant"),
  });

  const deleteMut = useMutation({
    mutationFn: deleteMerchant,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["merchants"] });
      toast.success("Merchant deleted");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to delete merchant"),
  });

  function openCreate() {
    setEditing(null);
    setName("");
    setModalOpen(true);
  }
  function openEdit(m: Merchant) {
    setEditing(m);
    setName(m.name);
    setModalOpen(true);
  }
  function closeModal() {
    setModalOpen(false);
    setEditing(null);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (editing) {
      updateMut.mutate({ id: editing.id, name });
    } else {
      createMut.mutate({ name });
    }
  }

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-xl font-bold text-foreground">Merchants</h1>
        <div className="flex gap-3 items-center">
          <Button onClick={openCreate}>
            <Plus className="size-4" /> Add Merchant
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
            <Clock className="size-4 animate-spin" /> Loading merchants...
          </div>
        ) : queryError ? (
          <div className="p-12 text-center text-destructive text-sm">
            Failed to load merchants.
          </div>
        ) : merchants.length === 0 ? (
          <div className="p-12 text-center text-muted-foreground text-sm">
            No merchants found.
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Updated</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {merchants.map((m) => (
                  <TableRow key={m.id}>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      {m.id}
                    </TableCell>
                    <TableCell>{m.name}</TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {formatDate(m.created_at)}
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {formatDate(m.updated_at)}
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => openEdit(m)}
                          title="Edit"
                        >
                          <Pencil className="size-3.5" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => setDeleteTarget(m)}
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
              limit={PAGE_SIZE}
              onPageChange={setPage}
              page={currentPage}
              totalPages={totalPages}
              total={totalCount}
            />
          </>
        )}
      </div>

      <Dialog
        open={modalOpen}
        onOpenChange={(open) => {
          if (!open) closeModal();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editing ? "Edit Merchant" : "Create Merchant"}
            </DialogTitle>
          </DialogHeader>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label>Name</Label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Merchant name"
                required
                autoFocus
              />
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

      <AlertDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Merchant</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete{" "}
              <span className="font-semibold">{deleteTarget?.name}</span>? This
              action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (deleteTarget) {
                  deleteMut.mutate(deleteTarget.id);
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
