import type { IListParams, IPaginationInfo } from 'shared/model/common.dto';

export type ListingType = 'rent' | 'exchange' | 'sale';
export type ListingStatus = 'active' | 'paused' | 'archived';

export const LISTING_TYPES: ListingType[] = ['rent', 'exchange', 'sale'];

export interface IListingOwnerDTO {
  id: string;
  name: string;
  campus: string;
}

export interface IListingDTO {
  id: string;
  created_at: string;
  modified_at: string;
  owner_id: string;
  owner: IListingOwnerDTO | null;
  type: ListingType;
  status: ListingStatus;
  title: string;
  description: string;
  category: string;
  campus: string;
  photos: string[];
  price: number; // KZT: per day (rent), fixed price (sale), 0 (exchange)
  deposit: number;
  wanted: string;
  available_from: string; // YYYY-MM-DD or ''
  available_to: string;
}

export interface IListingListParams {
  query?: string;
  type?: ListingType;
  category?: string;
  campus?: string;
  price_min?: number;
  price_max?: number;
  free_from?: string;
  free_to?: string;
  owner_id?: string;
  status?: ListingStatus;
  mine?: boolean;
  list_params: IListParams;
}

export interface IListingListRep {
  pagination_info: IPaginationInfo;
  results: IListingDTO[];
}

// create/update body: fields that don't belong to the type are not sent
export interface IListingEdit {
  type?: ListingType;
  title: string;
  description: string;
  category: string;
  campus: string;
  photos: string[];
  price?: number;
  deposit?: number;
  wanted?: string;
  available_from?: string;
  available_to?: string;
}

export interface IListingCreateRep {
  id: string;
}
