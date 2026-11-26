import type { IDictDTO } from './dict.dto';
import { store } from 'app/providers/store';
import { defineStore } from 'pinia';

const STORE_ID = 'dict';

export const useDictStore = defineStore(STORE_ID, () => {
  const dict = shallowRef<IDictDTO | null>(null);

  const categories = computed(() => dict.value?.categories ?? []);
  const campuses = computed(() => dict.value?.campuses ?? []);

  return { dict, categories, campuses };
});

export function useDictStoreOutside() {
  return useDictStore(store);
}
