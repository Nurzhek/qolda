import { http } from 'shared/api/http';

export class FileRepository {
  static createFile(file: File) {
    const data = new FormData();
    data.append('file', file);
    return http.request<{ url: string }>({
      method: 'POST',
      url: '/file',
      data,
      // Explicit multipart: otherwise axios serializes FormData as JSON because of the default Content-Type
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  }
}
