import type {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';
import axios from 'axios';
import { API_BASE_URL } from 'shared/config';
import { Logger } from 'shared/lib/utils/logger';
import { useAuthStoreOutside } from 'shared/stores/auth.store';

enum StatusCode {
  Unauthorized = 401,
  Forbidden = 403,
}

// grpc-gateway reads nested fields with dots and arrays as repeated keys:
// { list_params: { sort: ['-price'] } } → list_params.sort=-price
function serializeParams(params: Record<string, unknown>): string {
  const query = new URLSearchParams();
  const walk = (value: unknown, key: string) => {
    if (value === undefined || value === null || value === '')
      return;
    if (Array.isArray(value))
      value.forEach(v => walk(v, key));
    else if (typeof value === 'object')
      Object.entries(value).forEach(([k, v]) => walk(v, key ? `${key}.${k}` : k));
    else
      query.append(key, String(value));
  };
  walk(params, '');
  return query.toString();
}

// The only URL outside repositories: refresh is transport mechanics, not a slice scenario
const REFRESH_URL = '/auth/refresh';

interface RequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

interface ITokenPair {
  access_token: string;
  refresh_token: string;
}

interface IHttpHooks {
  onSessionLost?: () => void;
  onForbidden?: () => void;
}

const hooks: IHttpHooks = {};

// The app layer injects the reaction to session loss and 403: shared knows nothing about the router
export function setHttpHooks(value: IHttpHooks) {
  Object.assign(hooks, value);
}

const _logger = new Logger('Http');

class Http {
  private readonly instance: AxiosInstance;
  private refreshing: Promise<boolean> | null = null;

  constructor(private readonly baseURL: string) {
    this.instance = axios.create({
      baseURL,
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json; charset=utf-8',
      },
      paramsSerializer: { serialize: serializeParams },
    });

    this.instance.interceptors.request.use(this.handleRequest);
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error: AxiosError) => this.handleError(error),
    );
  }

  private handleRequest(this: void, config: InternalAxiosRequestConfig) {
    const { token } = useAuthStoreOutside();
    if (token && !config.headers.Authorization)
      config.headers.Authorization = `Bearer ${token}`;
    return config;
  }

  private async handleError(error: AxiosError) {
    const { response } = error;
    const originalRequest = error.config as RequestConfig | undefined;

    if (response?.status === StatusCode.Unauthorized && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true;

      if (await this.refreshTokens()) {
        originalRequest.headers.Authorization = `Bearer ${useAuthStoreOutside().token}`;
        return this.instance.request(originalRequest);
      }

      _logger.log('Session lost');
      useAuthStoreOutside().logout();
      hooks.onSessionLost?.();
    } else if (response?.status === StatusCode.Forbidden) {
      hooks.onForbidden?.();
    }

    return Promise.reject(error);
  }

  // Concurrent 401s wait for one shared refresh instead of each sending its own
  private refreshTokens() {
    this.refreshing ??= (async () => {
      const authStore = useAuthStoreOutside();
      if (!authStore.refreshToken)
        return false;
      try {
        const { data } = await axios.post<ITokenPair>(`${this.baseURL}${REFRESH_URL}`, {
          refresh_token: authStore.refreshToken,
        });
        authStore.setTokens(data.access_token, data.refresh_token);
        return true;
      } catch {
        return false;
      }
    })().finally(() => {
      this.refreshing = null;
    });

    return this.refreshing;
  }

  request<T = unknown>(config: AxiosRequestConfig) {
    return this.instance.request<T>(config);
  }
}

export const http = new Http(API_BASE_URL);
