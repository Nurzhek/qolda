import type { ListingStatus } from 'domains/entities/listing';
import { ListingRepository, useListingStore } from 'domains/entities/listing';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';
import { useI18n } from 'vue-i18n';

// What the owner can do from each status
export const STATUS_ACTIONS: Record<ListingStatus, ListingStatus[]> = {
  active: ['paused', 'archived'],
  paused: ['active', 'archived'],
  archived: ['active'],
};

export function useListingStatus() {
  const loadingState = useLoading();
  const listingStore = useListingStore();
  const toast = useToast();
  const { t } = useI18n();

  async function setStatus(id: string, status: ListingStatus) {
    const result = await Executor.run({
      request: ListingRepository.updateListingStatus(id, status),
      loadingState,
    });
    if (result.ok) {
      listingStore.patchStatus(id, status);
      toast.add({ title: t(`listing.my.changed.${status}`), color: 'success' });
    }
    return result.ok;
  }

  return { isLoading: loadingState.isLoading, setStatus };
}
