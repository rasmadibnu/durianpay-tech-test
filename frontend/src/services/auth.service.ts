import api from "@/lib/api";
import type { LoginRequest } from "@/types/auth";
import type { User } from "@/types/user";

// Auth
export async function login(data: LoginRequest): Promise<User> {
  const res = await api.post<User>("/dashboard/v1/auth/login", data);
  return res.data;
}
