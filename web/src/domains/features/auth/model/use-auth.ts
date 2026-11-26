import { useUserService, useUserStore } from 'domains/entities/user';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';
import { PageName } from 'shared/lib/types/app/pages';
import { useAuthStore } from 'shared/stores/auth.store';
import { AuthRepository } from '../api/auth.repository';

// No open redirects: only return to our own path
export function safeRedirect(value: unknown) {
  return typeof value === 'string' && value.startsWith('/') && !value.startsWith('//') ? value : '/';
}

export function useAuth() {
  const router = useRouter();
  const authStore = useAuthStore();
  const userStore = useUserStore();
  const { fetchMe } = useUserService();
  const loadingState = useLoading();

  function signInWithGoogle(redirect: string) {
    window.location.assign(AuthRepository.googleSignInUrl(safeRedirect(redirect)));
  }

  async function finishSignIn(access: string, refresh: string, redirect: string) {
    authStore.setTokens(access, refresh);
    userStore.reset();
    if (!(await fetchMe()))
      return false;
    await router.replace(safeRedirect(redirect));
    return true;
  }

  async function signInDev(values: { email: string; name: string }, redirect: string) {
    const result = await Executor.run({
      request: AuthRepository.createDevSession(values),
      loadingState,
    });
    if (!result.ok)
      return result.error;
    await finishSignIn(result.data.access_token, result.data.refresh_token, redirect);
    return null;
  }

  function signOut() {
    authStore.logout();
    userStore.reset();
    router.replace({ name: PageName.SIGN_IN });
  }

  return {
    isLoading: loadingState.isLoading,
    signInWithGoogle,
    signInDev,
    finishSignIn,
    signOut,
  };
}
