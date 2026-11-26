import { storeToRefs } from 'pinia';
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';
import { DictRepository } from '../api/dict.repository';
import { useDictStoreOutside } from './dict.store';

export function useDictService() {
  const loadingState = useLoading();
  const { dict } = storeToRefs(useDictStoreOutside());
  const hasError = shallowRef(false);

  async function fetchDict() {
    hasError.value = false;
    return Executor.run({
      request: DictRepository.readDict(),
      loadingState,
      ignoreError: true,
      onResult(data) {
        dict.value = data;
      },
      onError() {
        hasError.value = true;
      },
    });
  }

  // The dictionary doesn't change within a session: load it once
  async function ensureDict() {
    if (!dict.value)
      await fetchDict();
  }

  return {
    dict,
    hasError,
    isLoading: loadingState.isLoading,
    fetchDict,
    ensureDict,
  };
}
