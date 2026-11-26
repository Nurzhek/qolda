import type { AxiosError } from 'axios';
import type { IListingDTO, IListingEdit } from 'domains/entities/listing';
import type { IErrorData } from 'shared/lib/types/app/app';
import type { ListingFormValues } from './listing-form.schema';
import { ListingRepository } from 'domains/entities/listing';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';

// Send only the fields of the chosen type: the backend rejects e.g. a price on an exchange
export function toListingPayload(values: ListingFormValues): IListingEdit {
  const base = {
    title: values.title.trim(),
    description: values.description?.trim() ?? '',
    category: values.category,
    campus: values.campus ?? '',
    photos: values.photos ?? [],
  };

  switch (values.type) {
    case 'rent':
      return {
        ...base,
        price: values.price ?? 0,
        deposit: values.deposit ?? 0,
        available_from: values.available_from ?? '',
        available_to: values.available_to ?? '',
      };
    case 'sale':
      return { ...base, price: values.price ?? 0 };
    default:
      return { ...base, wanted: values.wanted?.trim() ?? '' };
  }
}

export function toFormValues(listing: IListingDTO): Partial<ListingFormValues> {
  return {
    type: listing.type,
    title: listing.title,
    description: listing.description,
    category: listing.category,
    campus: listing.campus,
    photos: [...listing.photos],
    price: listing.type === 'exchange' ? undefined : listing.price,
    deposit: listing.type === 'rent' ? listing.deposit : undefined,
    wanted: listing.wanted || undefined,
    available_from: listing.available_from || undefined,
    available_to: listing.available_to || undefined,
  };
}

type SaveResult = { ok: true; id: string } | { ok: false; error: AxiosError<IErrorData> };

export function useListingSave() {
  const loadingState = useLoading();

  // Returns the listing id or the request error (for form fields)
  async function save(values: ListingFormValues, id?: string): Promise<SaveResult> {
    const payload = toListingPayload(values);

    if (id) {
      const result = await Executor.run({
        request: ListingRepository.updateListing(id, payload),
        loadingState,
      });
      return result.ok ? { ok: true, id } : { ok: false, error: result.error };
    }

    const result = await Executor.run({
      request: ListingRepository.createListing({ ...payload, type: values.type }),
      loadingState,
    });
    return result.ok ? { ok: true, id: result.data.id } : { ok: false, error: result.error };
  }

  return { isSaving: loadingState.isLoading, save };
}
