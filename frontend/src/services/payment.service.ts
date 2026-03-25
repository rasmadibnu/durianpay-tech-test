import api from "@/lib/api";
import type { Payment, PaymentListResponse } from "@/types/payment";

export async function getPayments(params: {
  status?: string;
  search?: string;
  page?: number;
  limit?: number;
}): Promise<PaymentListResponse> {
  const res = await api.get<PaymentListResponse>("/dashboard/v1/payments", {
    params,
  });
  return res.data;
}

export async function createPayment(data: {
  merchant_id: number;
  amount: string;
  status: string;
}): Promise<Payment> {
  const res = await api.post<Payment>("/dashboard/v1/payments", data);
  return res.data;
}

export async function updatePayment(
  id: string,
  data: { merchant_id: number; amount: string; status: string },
): Promise<Payment> {
  const res = await api.put<Payment>(`/dashboard/v1/payments/${id}`, data);
  return res.data;
}

export async function reviewPayment(id: string, status: string): Promise<void> {
  await api.put(`/dashboard/v1/payments/${id}/review`, { status });
}

export async function deletePayment(id: string): Promise<void> {
  await api.delete(`/dashboard/v1/payments/${id}`);
}
