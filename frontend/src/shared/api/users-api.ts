import { buildUsersListUrl, httpClient } from './http';

export type UserDto = {
  user_id: number;
  external_id: string;
  username: string;
  timezone: string;
  created_at: string;
  updated_at: string;
};

export async function fetchUsers(usersBaseUrl: string): Promise<UserDto[]> {
  const response = await httpClient.get<UserDto[]>(buildUsersListUrl(usersBaseUrl));
  return response.data;
}
