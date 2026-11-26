import type { IListingListParams } from './listing.dto';
import { storeToRefs } from 'pinia';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';
import { ListingRepository } from '../api/listing.repository';
import { useListingStore } from './listing.store';

// The entity service handles reads. Changes live in the listing-form and listing-status features.
export function useListingService() {
  const loadingState = useLoading();
  const { items, total, current, isEmpty } = storeToRefs(useListingStore());
  const hasError = shallowRef(false);

  async function fetchListings(params: IListingListParams) {
    hasError.value = false;
    return Executor.run({
      request: ListingRepository.fetchListings(params),
      loadingState,
      ignoreError: true,
      onResult(data) {
        items.value = data.results;
        total.value = Number(data.pagination_info?.total_count ?? data.results.length);
      },
      onError() {
        hasError.value = true;
      },
    });
  }

  async function fetchListing(id: string) {
    hasError.value = false;
    current.value = null;
    return Executor.run({
      request: ListingRepository.readListing(id),
      loadingState,
      ignoreError: true,
      onResult(data) {
        current.value = data;
      },
      onError() {
        hasError.value = true;
      },
    });
  }

  return {
    items,
    total,
    current,
    isEmpty,
    hasError,
    isLoading: loadingState.isLoading,
    fetchListings,
    fetchListing,
  };
}
