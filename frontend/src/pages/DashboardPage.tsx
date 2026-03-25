import { useState } from "react";
import {
  useQuery,
  useQueries,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";
import {
  getPayments,
  createPayment,
  updatePayment,
  reviewPayment,
  deletePayment,
} from "@/services/payment.service";
import { getMerchants } from "@/services/merchant.service";
import type { Payment } from "@/types/payment";
import type { Merchant } from "@/types/merchant";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import TablePagination from "@/components/TablePagination";
import { cn, formatDate } from "@/lib/utils";
import {
  Plus,
  Eye,
  Pencil,
  Trash2,
  CreditCard,
  CheckCircle2,
  Clock,
  XCircle,
} from "lucide-react";

const STATUS_OPTIONS = [
  { value: "", label: "All" },
  { value: "completed", label: "Completed" },
  { value: "processing", label: "Processing" },
  { value: "failed", label: "Failed" },
] as const;

const PAGE_SIZE = 10;

const badgeVariant: Record<
  string,
  "default" | "secondary" | "destructive" | "outline"
> = {
  completed: "default",
  processing: "secondary",
  failed: "destructive",
};

function formatCurrency(amount: string) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(parseFloat(amount));
}

export default function DashboardPage() {
  const { role } = useAuth();
  const qc = useQueryClient();
  const isOperation = role === "operation";

  const [statusFilter, setStatusFilter] = useState("");
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [modalOpen, setModalOpen] = useState(false);
  const [reviewOpen, setReviewOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<Payment | null>(null);
  const [editing, setEditing] = useState<Payment | null>(null);
  const [reviewTarget, setReviewTarget] = useState<Payment | null>(null);

  const [merchantId, setMerchantId] = useState("");
  const [amount, setAmount] = useState("");
  const [status, setStatus] = useState("processing");
  const [reviewStatus, setReviewStatus] = useState("completed");

  const {
    data,
    isLoading,
    error: queryError,
  } = useQuery({
    queryKey: ["payments", statusFilter, search, page],
    queryFn: () =>
      getPayments({
        status: statusFilter || undefined,
        search: search || undefined,
        page,
        limit: PAGE_SIZE,
      }),
  });

  const { data: merchantsData } = useQuery({
    queryKey: ["merchants-all"],
    queryFn: () => getMerchants({ page: 1, limit: 100 }),
    enabled: isOperation,
  });

  const merchants: Merchant[] = merchantsData?.merchants ?? [];
  const payments: Payment[] = data?.payments ?? [];
  const totalCount = data?.total ?? 0;
  const totalPages = data?.total_pages ?? 1;
  const currentPage = page;

  const paymentSummaryQueries = useQueries({
    queries: [
      {
        queryKey: ["payment-summary", "all"],
        queryFn: () => getPayments({ page: 1, limit: 1 }),
      },
      {
        queryKey: ["payment-summary", "completed"],
        queryFn: () => getPayments({ status: "completed", page: 1, limit: 1 }),
      },
      {
        queryKey: ["payment-summary", "processing"],
        queryFn: () => getPayments({ status: "processing", page: 1, limit: 1 }),
      },
      {
        queryKey: ["payment-summary", "failed"],
        queryFn: () => getPayments({ status: "failed", page: 1, limit: 1 }),
      },
    ],
  });

  const [allSummary, completedSummary, processingSummary, failedSummary] =
    paymentSummaryQueries;

  const createMut = useMutation({
    mutationFn: (d: { merchant_id: number; amount: string; status: string }) =>
      createPayment(d),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["payments"] });
      closeModal();
      toast.success("Payment created successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to create payment"),
  });

  const updateMut = useMutation({
    mutationFn: ({
      id,
      ...d
    }: {
      id: string;
      merchant_id: number;
      amount: string;
      status: string;
    }) => updatePayment(id, d),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["payments"] });
      qc.invalidateQueries({ queryKey: ["payment-summary"] });
      closeModal();
      toast.success("Payment updated successfully");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to update payment"),
  });

  const reviewMut = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      reviewPayment(id, status),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["payments"] });
      qc.invalidateQueries({ queryKey: ["payment-summary"] });
      closeReview();
      toast.success("Payment status updated");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to update status"),
  });

  const deleteMut = useMutation({
    mutationFn: deletePayment,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["payments"] });
      qc.invalidateQueries({ queryKey: ["payment-summary"] });
      toast.success("Payment deleted");
    },
    onError: (e: any) =>
      toast.error(e?.response?.data?.message || "Failed to delete payment"),
  });

  function openCreate() {
    setEditing(null);
    setMerchantId(merchants[0]?.id?.toString() ?? "");
    setAmount("");
    setStatus("processing");
    setModalOpen(true);
  }

  function openEdit(p: Payment) {
    setEditing(p);
    setMerchantId(p.merchant_id.toString());
    setAmount(p.amount);
    setStatus(p.status);
    setModalOpen(true);
  }

  function openReview(p: Payment) {
    setReviewTarget(p);
    setReviewStatus(p.status);
    setReviewOpen(true);
  }

  function closeModal() {
    setModalOpen(false);
    setEditing(null);
  }
  function closeReview() {
    setReviewOpen(false);
    setReviewTarget(null);
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const payload = { merchant_id: Number(merchantId), amount, status };
    if (editing) {
      updateMut.mutate({ id: editing.id, ...payload });
    } else {
      createMut.mutate(payload);
    }
  }

  function handleReview(e: React.FormEvent) {
    e.preventDefault();
    if (reviewTarget)
      reviewMut.mutate({ id: reviewTarget.id, status: reviewStatus });
  }

  function handleFilterChange(value: string) {
    setStatusFilter(value);
    setPage(1);
  }

  function handleSearchChange(value: string) {
    setSearch(value);
    setPage(1);
  }

  const widgets = [
    {
      label: "Total Payments",
      value: allSummary.data?.total ?? 0,
      icon: CreditCard,
      color: "text-indigo-400 bg-indigo-500/10",
    },
    {
      label: "Success",
      value: completedSummary.data?.total ?? 0,
      icon: CheckCircle2,
      color: "text-emerald-400 bg-emerald-500/10",
    },
    {
      label: "Processing",
      value: processingSummary.data?.total ?? 0,
      icon: Clock,
      color: "text-amber-400 bg-amber-500/10",
    },
    {
      label: "Failed",
      value: failedSummary.data?.total ?? 0,
      icon: XCircle,
      color: "text-red-400 bg-red-500/10",
    },
  ];

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      {/* Widgets */}
      <div className="grid grid-cols-4 gap-4 max-lg:grid-cols-2 max-sm:grid-cols-1">
        {widgets.map((w) => (
          <div
            key={w.label}
            className="flex items-center gap-4 p-5 rounded-xl border border-border bg-card"
          >
            <div
              className={cn(
                "flex items-center justify-center size-11 rounded-xl",
                w.color,
              )}
            >
              <w.icon className="size-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground uppercase tracking-wider font-medium">
                {w.label}
              </p>
              <p className="text-2xl font-bold text-foreground">{w.value}</p>
            </div>
          </div>
        ))}
      </div>

      {/* Toolbar */}
      <div className="flex justify-between items-center flex-wrap gap-3">
        <h2 className="text-lg font-semibold text-foreground">Transactions</h2>
        <div className="flex gap-3 items-center">
          <div className="flex gap-1 bg-muted/50 border border-border rounded-lg p-1">
            {STATUS_OPTIONS.map((opt) => (
              <button
                key={opt.value}
                className={cn(
                  "px-3 py-1.5 rounded-md text-sm font-medium transition-colors flex items-center gap-1.5",
                  statusFilter === opt.value
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground",
                )}
                onClick={() => handleFilterChange(opt.value)}
              >
                {opt.value === "completed" && (
                  <span className="size-1.5 rounded-full bg-emerald-400" />
                )}
                {opt.value === "processing" && (
                  <span className="size-1.5 rounded-full bg-amber-400" />
                )}
                {opt.value === "failed" && (
                  <span className="size-1.5 rounded-full bg-red-400" />
                )}
                {opt.label}
              </button>
            ))}
          </div>
          {isOperation && (
            <Button onClick={openCreate}>
              <Plus className="size-4" />
              Add Payment
            </Button>
          )}
          <Input
            placeholder="Search..."
            value={search}
            onChange={(e) => handleSearchChange(e.target.value)}
            className="w-48"
          />
        </div>
      </div>

      {/* Table */}
      <div className="rounded-xl border border-border bg-card overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center gap-3 p-12 text-muted-foreground text-sm">
            <Clock className="size-4 animate-spin" /> Loading payments...
          </div>
        ) : queryError ? (
          <div className="p-12 text-center text-destructive text-sm">
            Failed to load payments.
          </div>
        ) : totalCount === 0 ? (
          <div className="p-12 text-center text-muted-foreground text-sm">
            No payments found.
          </div>
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Payment ID</TableHead>
                  <TableHead>Merchant</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {payments.map((p) => (
                  <TableRow key={p.id}>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      {p.id}
                    </TableCell>
                    <TableCell>{p.merchant_name}</TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {formatDate(p.created_at)}
                    </TableCell>
                    <TableCell className="font-semibold tabular-nums">
                      {formatCurrency(p.amount)}
                    </TableCell>
                    <TableCell>
                      <Badge variant={badgeVariant[p.status] ?? "outline"}>
                        {p.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        <Button
                          variant="ghost"
                          size="icon-xs"
                          onClick={() => openReview(p)}
                          title="Review"
                        >
                          <Eye className="size-3.5" />
                        </Button>
                        {isOperation && (
                          <>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              onClick={() => openEdit(p)}
                              title="Edit"
                            >
                              <Pencil className="size-3.5" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon-xs"
                              onClick={() => setDeleteTarget(p)}
                              title="Delete"
                              className="hover:text-destructive"
                            >
                              <Trash2 className="size-3.5" />
                            </Button>
                          </>
                        )}
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

      {/* Create/Edit Dialog */}
      <Dialog
        open={modalOpen}
        onOpenChange={(open) => {
          if (!open) closeModal();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editing ? "Edit Payment" : "Create Payment"}
            </DialogTitle>
          </DialogHeader>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label>Merchant</Label>
              <Select value={merchantId} onValueChange={setMerchantId} required>
                <SelectTrigger>
                  <SelectValue placeholder="Select merchant" />
                </SelectTrigger>
                <SelectContent>
                  {merchants.map((m) => (
                    <SelectItem key={m.id} value={m.id.toString()}>
                      {m.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Amount</Label>
              <Input
                type="number"
                step="0.01"
                min="0"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="0.00"
                required
              />
            </div>
            <div className="space-y-2">
              <Label>Status</Label>
              <Select value={status} onValueChange={setStatus}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="processing">Processing</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="failed">Failed</SelectItem>
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

      {/* Review Dialog */}
      <Dialog
        open={reviewOpen}
        onOpenChange={(open) => {
          if (!open) closeReview();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Review Payment</DialogTitle>
          </DialogHeader>
          <form className="space-y-4" onSubmit={handleReview}>
            <div className="space-y-2">
              <Label>Payment ID</Label>
              <Input value={reviewTarget?.id ?? ""} disabled />
            </div>
            <div className="space-y-2">
              <Label>Merchant</Label>
              <Input value={reviewTarget?.merchant_name ?? ""} disabled />
            </div>
            <div className="space-y-2">
              <Label>Amount</Label>
              <Input
                value={reviewTarget ? formatCurrency(reviewTarget.amount) : ""}
                disabled
              />
            </div>
            <div className="space-y-2">
              <Label>Update Status</Label>
              <Select value={reviewStatus} onValueChange={setReviewStatus}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="processing">Processing</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="failed">Failed</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={closeReview}>
                Cancel
              </Button>
              <Button type="submit" disabled={reviewMut.isPending}>
                Update Status
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation AlertDialog */}
      <AlertDialog
        open={!!deleteTarget}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Payment</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete payment{" "}
              <span className="font-mono font-semibold">
                {deleteTarget?.id}
              </span>
              ? This action cannot be undone.
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
