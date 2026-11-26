import type { IUserDTO } from './user.dto';
import { store } from 'app/providers/store';
import { defineStore } from 'pinia';

const STORE_ID = 'user';

export const useUserStore = defineStore(STORE_ID, () => {
  const me = shallowRef<IUserDTO | null>(null);
  // The public profile being viewed
  const viewed = shallowRef<IUserDTO | null>(null);

  function reset() {
    me.value = null;
    viewed.value = null;
  }

  return { me, viewed, reset };
});

export function useUserStoreOutside() {
  return useUserStore(store);
}
