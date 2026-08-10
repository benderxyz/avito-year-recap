import type { RecapPayload } from '@recap-engine/core';
import { useQuery } from '@tanstack/react-query';
import { fetchRecapPayload } from '@/shared/api';
import { env } from '@/shared/config';

export function recapQueryKey(userId: string) {
  return ['recap', env.recapYear, userId] as const;
}

export function useRecapQuery(userId: string, enabled = true) {
  return useQuery({
    queryKey: recapQueryKey(userId),
    queryFn: () => fetchRecapPayload(env.cardsBaseUrl, env.recapYear, userId),
    enabled: enabled && Boolean(userId),
    staleTime: 60_000,
  });
}

export type { RecapPayload };
