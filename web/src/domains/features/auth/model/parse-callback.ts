export type CallbackResult
  = | { ok: true; access: string; refresh: string; redirect: string }
    | { ok: false; error: string };

// The backend returns tokens in the fragment (#...) so they stay out of logs and Referer
export function parseCallbackHash(hash: string): CallbackResult {
  const params = new URLSearchParams(hash.replace(/^#/, ''));
  const access = params.get('access_token');
  const refresh = params.get('refresh_token');
  if (access && refresh)
    return { ok: true, access, refresh, redirect: params.get('redirect') ?? '/' };
  return { ok: false, error: params.get('error') || 'oauth_failed' };
}
