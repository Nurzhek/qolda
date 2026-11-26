import type { IDictDTO } from '../model/dict.dto';
import { http } from 'shared/api/http';

export class DictRepository {
  static readDict() {
    return http.request<IDictDTO>({
      method: 'GET',
      url: '/dict',
    });
  }
}
