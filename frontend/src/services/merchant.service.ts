import api from "@/lib/api";
import type { Merchant, MerchantListResponse } from "@/types/merchant";

export async function getMerchants(params: {
  search?: string;
  page?: number;
  limit?: number;
}): Promise<MerchantListResponse> {
  const res = await api.get<MerchantListResponse>("/dashboard/v1/merchants", {
    params,
  });
  return res.data;
}

export async function getMerchant(id: number): Promise<Merchant> {
  const res = await api.get<Merchant>(`/dashboard/v1/merchants/${id}`);
  return res.data;
}

export async function createMerchant(data: {
  name: string;
}): Promise<Merchant> {
  const res = await api.post<Merchant>("/dashboard/v1/merchants", data);
  return res.data;
}

export async function updateMerchant(
  id: number,
  data: { name: string },
): Promise<Merchant> {
  const res = await api.put<Merchant>(`/dashboard/v1/merchants/${id}`, data);
  return res.data;
}

export async function deleteMerchant(id: number): Promise<void> {
  await api.delete(`/dashboard/v1/merchants/${id}`);
}
