import type { RecapPayload } from '@recap-engine/core';
import { useQuery } from '@tanstack/react-query';
import { fetchSharedRecapPayload } from '@/shared/api';
import { env } from '@/shared/config';

export function sharedRecapQueryKey(token: string) {
  return ['shared-recap', token] as const;
}

export function useSharedRecapQuery(token: string, enabled = true) {
  return useQuery({
    queryKey: sharedRecapQueryKey(token),
    queryFn: () => fetchSharedRecapPayload(env.cardsBaseUrl, token),
    enabled: enabled && Boolean(token),
    staleTime: 60_000,
  });
}

export type { RecapPayload };
