import { store } from 'app/providers/store';
import { defineStore } from 'pinia';

interface IAuthState {
  token: string | null;
  refreshToken: string | null;
}

const STORE_ID = 'auth';

export const useAuthStore = defineStore(
  STORE_ID,
  () => {
    const state = reactive<IAuthState>({
      token: null,
      refreshToken: null,
    });

    function setTokens(access: string | null, refresh: string | null) {
      state.token = access;
      state.refreshToken = refresh;
    }

    const isLoggedIn = computed(() => !!state.token);

    function logout() {
      state.token = null;
      state.refreshToken = null;
    }

    return { ...toRefs(state), setTokens, logout, isLoggedIn };
  },
  {
    persist: {
      key: 'qolda-auth',
    },
  },
);

// Need to be used outside the setup
export function useAuthStoreOutside() {
  return useAuthStore(store);
}
