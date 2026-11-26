import type { IListingListParams, ListingType } from 'domains/entities/listing';
import type { LocationQuery, LocationQueryRaw } from 'vue-router';
import { LISTING_TYPES } from 'domains/entities/listing';

export const PAGE_SIZE = 24;

export const SORTS = {
  newest: '-created_at',
  priceAsc: 'price',
  priceDesc: '-price',
} as const;

export type SortKey = keyof typeof SORTS;

export interface IListingFilterState {
  q: string;
  type?: ListingType;
  category?: string;
  campus?: string;
  priceMin?: number;
  priceMax?: number;
  from?: string;
  to?: string;
  sort: SortKey;
  page: number; // zero-based, as in the API
}

function str(value: LocationQuery[string]) {
  const first = Array.isArray(value) ? value[0] : value;
  return first ? String(first) : undefined;
}

function int(value: LocationQuery[string]) {
  const parsed = Number.parseInt(str(value) ?? '', 10);
  return Number.isFinite(parsed) && parsed >= 0 ? parsed : undefined;
}

function isListingType(value?: string): value is ListingType {
  return !!value && (LISTING_TYPES as string[]).includes(value);
}

function isSortKey(value?: string): value is SortKey {
  return !!value && value in SORTS;
}

export function parseFilterQuery(query: LocationQuery): IListingFilterState {
  const type = str(query.type);
  const sort = str(query.sort);
  return {
    q: str(query.q) ?? '',
    type: isListingType(type) ? type : undefined,
    category: str(query.category),
    campus: str(query.campus),
    priceMin: int(query.price_min),
    priceMax: int(query.price_max),
    from: str(query.from),
    to: str(query.to),
    sort: isSortKey(sort) ? sort : 'newest',
    page: int(query.page) ?? 0,
  };
}

function toQuery(state: IListingFilterState): LocationQueryRaw {
  const query: LocationQueryRaw = {
    q: state.q.trim() || undefined,
    type: state.type,
    category: state.category,
    campus: state.campus,
    price_min: state.priceMin,
    price_max: state.priceMax,
    from: state.from,
    to: state.to,
    sort: state.sort === 'newest' ? undefined : state.sort,
    page: state.page > 0 ? state.page : undefined,
  };
  return Object.fromEntries(Object.entries(query).filter(([, v]) => v !== undefined && v !== ''));
}

export function toListParams(state: IListingFilterState): IListingListParams {
  return {
    query: state.q.trim() || undefined,
    type: state.type,
    category: state.category,
    campus: state.campus,
    price_min: state.priceMin,
    price_max: state.priceMax,
    free_from: state.from,
    free_to: state.to,
    list_params: {
      page: state.page,
      page_size: PAGE_SIZE,
      with_total_count: true,
      sort: [SORTS[state.sort]],
    },
  };
}

// Filter state lives in the URL, so a search can be shared in the dorm chat
export function useListingFilter() {
  const route = useRoute();
  const router = useRouter();

  const state = computed(() => parseFilterQuery(route.query));
  const params = computed(() => toListParams(state.value));

  const activeCount = computed(() => {
    const s = state.value;
    return [s.category, s.campus, s.priceMin ?? s.priceMax, s.from ?? s.to].filter(v => v !== undefined).length;
  });

  const hasAnyFilter = computed(() => activeCount.value > 0 || !!state.value.q || !!state.value.type);

  // Any filter change goes back to the first page
  function update(patch: Partial<IListingFilterState>) {
    router.replace({ query: toQuery({ ...state.value, page: 0, ...patch }) });
  }

  function reset() {
    router.replace({ query: toQuery({ q: '', sort: state.value.sort, page: 0 }) });
  }

  return { state, params, activeCount, hasAnyFilter, update, reset };
}
