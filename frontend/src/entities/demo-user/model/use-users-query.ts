import { useQuery } from '@tanstack/react-query';
import type { UserDto } from '@/shared/api/users-api';
import { fetchUsers } from '@/shared/api/users-api';
import { env } from '@/shared/config';
import type { DemoUser } from './demo-users';

export function usersQueryKey() {
  return ['users'] as const;
}

export function mapUserDto(dto: UserDto): DemoUser {
  return {
    id: String(dto.user_id),
    name: dto.username,
    externalId: dto.external_id,
  };
}

export function useUsersQuery() {
  return useQuery({
    queryKey: usersQueryKey(),
    queryFn: async () => {
      const users = await fetchUsers(env.usersBaseUrl);
      return users.map(mapUserDto);
    },
    staleTime: 60_000,
  });
}

export function useDemoUser(userId: string | undefined) {
  const { data: users } = useUsersQuery();
  if (!userId || !users) {
    return undefined;
  }

  return users.find((user) => user.id === userId);
}
