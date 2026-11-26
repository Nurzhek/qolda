export type CallbackResult
  = | { ok: true; access: string; refresh: string; redirect: string }
    | { ok: false; error: string };

// Set right before this browser leaves for Google. Tokens in a callback link this browser didn't start
// are rejected, so nobody can sign a victim into the attacker's account with a crafted link (login CSRF).
const PENDING_KEY = 'qolda-oauth-pending';

export function markGoogleSignInStarted() {
  sessionStorage.setItem(PENDING_KEY, '1');
}

// The backend returns tokens in the fragment (#...) so they stay out of logs and Referer
export function parseCallbackHash(hash: string): CallbackResult {
  const started = sessionStorage.getItem(PENDING_KEY) === '1';
  sessionStorage.removeItem(PENDING_KEY);

  const params = new URLSearchParams(hash.replace(/^#/, ''));
  const access = params.get('access_token');
  const refresh = params.get('refresh_token');
  if (access && refresh)
    return started ? { ok: true, access, refresh, redirect: params.get('redirect') ?? '/' } : { ok: false, error: 'oauth_failed' };
  return { ok: false, error: params.get('error') || 'oauth_failed' };
}
