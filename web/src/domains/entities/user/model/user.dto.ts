export type UserRole = 'student' | 'moderator' | 'admin';

export interface IUserDTO {
  id: string;
  created_at: string;
  email: string;
  name: string;
  campus: string;
  role: UserRole;
  telegram: string;
}

export interface IUserUpdate {
  name?: string;
  campus?: string;
  telegram?: string;
}
