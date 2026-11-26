import type { IUserUpdate } from 'domains/entities/user';
import { UserRepository, useUserStore } from 'domains/entities/user';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';

export function useProfileEdit() {
  const loadingState = useLoading();
  const userStore = useUserStore();

  // Returns the request error (for form fields) or null on success
  async function saveProfile(data: IUserUpdate) {
    const result = await Executor.run({
      request: UserRepository.updateMe(data),
      loadingState,
    });
    if (!result.ok)
      return result.error;
    if (userStore.me)
      userStore.me = { ...userStore.me, ...data };
    return null;
  }

  return {
    isLoading: loadingState.isLoading,
    saveProfile,
  };
}
