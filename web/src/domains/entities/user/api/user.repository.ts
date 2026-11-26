import type { IUserDTO, IUserUpdate } from '../model/user.dto';
import { http } from 'shared/api/http';

export class UserRepository {
  static readMe() {
    return http.request<IUserDTO>({
      method: 'GET',
      url: '/me',
    });
  }

  static updateMe(data: IUserUpdate) {
    return http.request<unknown>({
      method: 'PUT',
      url: '/me',
      data,
    });
  }

  static readUser(id: string) {
    return http.request<IUserDTO>({
      method: 'GET',
      url: `/usr/${id}`,
    });
  }
}
