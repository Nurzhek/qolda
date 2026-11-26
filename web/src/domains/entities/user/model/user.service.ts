import { storeToRefs } from 'pinia';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';
import { UserRepository } from '../api/user.repository';
import { useUserStoreOutside } from './user.store';

// Reading users. Profile changes live in the profile-edit feature
export function useUserService() {
  const loadingState = useLoading();
  const { me, viewed } = storeToRefs(useUserStoreOutside());
  const hasError = shallowRef(false);

  async function fetchMe() {
    const result = await Executor.run({
      request: UserRepository.readMe(),
      loadingState,
      ignoreError: true,
      onResult(data) {
        me.value = data;
      },
    });
    return result.ok;
  }

  async function fetchUser(id: string) {
    hasError.value = false;
    viewed.value = null;
    return Executor.run({
      request: UserRepository.readUser(id),
      loadingState,
      ignoreError: true,
      onResult(data) {
        viewed.value = data;
      },
      onError() {
        hasError.value = true;
      },
    });
  }

  return {
    me,
    viewed,
    hasError,
    isLoading: loadingState.isLoading,
    fetchMe,
    fetchUser,
  };
}
