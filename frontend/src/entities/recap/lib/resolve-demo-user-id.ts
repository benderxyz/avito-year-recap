import type { DemoUser } from '@/entities/demo-user';

export function resolveDemoUserId(metaUserId: string, users?: DemoUser[]): string {
  if (/^\d+$/.test(metaUserId)) {
    return metaUserId;
  }

  const match = users?.find((user) => user.externalId === metaUserId);
  return match?.id ?? metaUserId;
}
