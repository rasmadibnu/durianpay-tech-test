import type { PaginationResponse } from "./pagination";

export interface Merchant {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface MerchantListResponse extends PaginationResponse<Merchant> {
  merchants: Merchant[];
}
