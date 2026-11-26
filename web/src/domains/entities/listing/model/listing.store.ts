import type { IListingDTO, ListingStatus } from './listing.dto';
import { defineStore } from 'pinia';

const STORE_ID = 'listing';

// The store holds state and derived values. The service makes requests, not the store.
export const useListingStore = defineStore(STORE_ID, () => {
  const items = shallowRef<IListingDTO[]>([]);
  const total = shallowRef(0);
  const current = shallowRef<IListingDTO | null>(null);

  const isEmpty = computed(() => items.value.length === 0);

  // A status change from a feature is reflected in already loaded data without a refetch
  function patchStatus(id: string, status: ListingStatus) {
    items.value = items.value.map(item => (item.id === id ? { ...item, status } : item));
    if (current.value?.id === id)
      current.value = { ...current.value, status };
  }

  return { items, total, current, isEmpty, patchStatus };
});
