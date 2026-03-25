import type { PaginationResponse } from "./pagination";

export interface User {
  email: string;
  role: string;
  token: string;
}

export interface UserRecord {
  id: string;
  email: string;
  role: string;
}
export interface UserListResponse extends PaginationResponse<UserRecord> {
  users: UserRecord[];
}
