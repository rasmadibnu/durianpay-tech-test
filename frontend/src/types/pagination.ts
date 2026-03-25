export interface PaginationResponse<T> {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  [key: string]: T[] | number;
}
