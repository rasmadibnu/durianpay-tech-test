import type { PaginationResponse } from "./pagination";

export interface Payment {
  id: string;
  merchant_id: number;
  merchant_name: string;
  amount: string;
  status: "completed" | "processing" | "failed";
  created_at: string;
}

export interface PaymentListResponse extends PaginationResponse<Payment> {
  payments: Payment[];
}
