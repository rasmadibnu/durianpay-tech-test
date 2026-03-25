import api from "@/lib/api";
import type { UserListResponse, UserRecord } from "@/types/user";

export async function getUsers(params: {
  search?: string;
  page?: number;
  limit?: number;
}): Promise<UserListResponse> {
  const res = await api.get<UserListResponse>("/dashboard/v1/users", {
    params,
  });
  return res.data;
}

export async function getUser(id: number): Promise<UserRecord> {
  const res = await api.get<UserRecord>(`/dashboard/v1/users/${id}`);
  return res.data;
}

export async function createUser(data: {
  email: string;
  password: string;
  role: string;
}): Promise<UserRecord> {
  const res = await api.post<UserRecord>("/dashboard/v1/users", data);
  return res.data;
}

export async function updateUser(
  id: number,
  data: { email: string; role: string },
): Promise<UserRecord> {
  const res = await api.put<UserRecord>(`/dashboard/v1/users/${id}`, data);
  return res.data;
}

export async function updateUserPassword(
  id: number,
  password: string,
): Promise<void> {
  await api.patch(`/dashboard/v1/users/${id}/password`, { password });
}

export async function deleteUser(id: number): Promise<void> {
  await api.delete(`/dashboard/v1/users/${id}`);
}
