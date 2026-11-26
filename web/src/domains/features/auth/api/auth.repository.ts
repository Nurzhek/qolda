import { http } from 'shared/api/http';
import { API_BASE_URL } from 'shared/config';

export interface ITokenPairDTO {
  access_token: string;
  refresh_token: string;
}

export class AuthRepository {
  // Google OAuth is a full browser navigation: the backend redirects to Google and back
  static googleSignInUrl(redirect: string) {
    return `${API_BASE_URL}/auth/google?${new URLSearchParams({ redirect })}`;
  }

  static createDevSession(data: { email: string; name: string }) {
    return http.request<ITokenPairDTO>({
      method: 'POST',
      url: '/auth/dev',
      data,
    });
  }
}
