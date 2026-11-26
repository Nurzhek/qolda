import type { AxiosError, AxiosResponse } from 'axios';
import type { IErrorData, LoadingService } from 'shared/lib/types/app/app';
import { i18n } from 'i18n/index';
import { useGlobalLoading } from '../lib/composables/useLoading';
import { getErrorMessage } from './error-handler';

export interface ExecutorOptions<T> {
  request: Promise<AxiosResponse<T>>;
  loadingState?: LoadingService;
  ignoreError?: boolean;
  onResult?: (data: T, response: AxiosResponse<T>) => void;
  onError?: (error: AxiosError<IErrorData>) => void;
  onComplete?: () => void;
}

export type ExecutorResult<T>
  = | { ok: true; data: T }
    | { ok: false; error: AxiosError<IErrorData> };

const SILENT_STATUSES = [401, 403];
const TOAST_DURATION = 4000;

const globalLoadingState = useGlobalLoading();

// Field errors go to the form; a toast for them is redundant
function hasFieldErrors(error: AxiosError<IErrorData>) {
  const { data } = error.response ?? {};
  return data?.code === 'invalid_request' && !!data.fields && Object.keys(data.fields).length > 0;
}

export abstract class Executor {
  static async run<T = unknown>(options: ExecutorOptions<T>): Promise<ExecutorResult<T>> {
    const { request, onResult, onError, onComplete, ignoreError } = options;
    const loadingState = options.loadingState ?? globalLoadingState;

    try {
      loadingState.startLoading();

      const result = await request;
      onResult?.(result.data, result);

      return { ok: true, data: result.data };
    } catch (raw) {
      const error = raw as AxiosError<IErrorData>;
      const status = error.response?.status ?? 0;

      if (!ignoreError && !SILENT_STATUSES.includes(status) && !hasFieldErrors(error)) {
        useToast().add({
          title: i18n.t(getErrorMessage(error.response?.data?.code ?? 'service_not_available')),
          color: 'error',
          duration: TOAST_DURATION,
        });
      }

      onError?.(error);

      return { ok: false, error };
    } finally {
      loadingState.stopLoading();
      onComplete?.();
    }
  }
}
