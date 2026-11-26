import type {
  IListingCreateRep,
  IListingDTO,
  IListingEdit,
  IListingListParams,
  IListingListRep,
  ListingStatus,
} from '../model/listing.dto';
import { http } from 'shared/api/http';

// The only place where listing URLs live.
export class ListingRepository {
  static fetchListings(params: IListingListParams) {
    return http.request<IListingListRep>({
      method: 'GET',
      url: '/listing',
      params,
    });
  }

  static readListing(id: string) {
    return http.request<IListingDTO>({
      method: 'GET',
      url: `/listing/${id}`,
    });
  }

  static createListing(data: IListingEdit) {
    return http.request<IListingCreateRep>({
      method: 'POST',
      url: '/listing',
      data,
    });
  }

  static updateListing(id: string, data: IListingEdit) {
    return http.request<unknown>({
      method: 'PUT',
      url: `/listing/${id}`,
      data,
    });
  }

  static updateListingStatus(id: string, status: ListingStatus) {
    return http.request<unknown>({
      method: 'PUT',
      url: `/listing/${id}/status`,
      data: { status },
    });
  }
}
