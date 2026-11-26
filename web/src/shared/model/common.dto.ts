export interface IListParams {
  page?: number;
  page_size?: number;
  with_total_count?: boolean;
  sort?: string[];
}

export interface IPaginationInfo {
  page: string;
  page_size: string;
  total_count: string;
}
