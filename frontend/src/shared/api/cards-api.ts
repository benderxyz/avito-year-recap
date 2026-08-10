import type { RecapPayload } from '@recap-engine/core';
import { buildRecapUrl, buildShareRecapUrl, httpClient } from './http';

export async function fetchRecapPayload(
  cardsBaseUrl: string,
  year: number,
  userId: string,
): Promise<RecapPayload> {
  const response = await httpClient.get<RecapPayload>(buildRecapUrl(cardsBaseUrl, year, userId));
  return response.data;
}

export async function fetchSharedRecapPayload(
  cardsBaseUrl: string,
  token: string,
): Promise<RecapPayload> {
  const response = await httpClient.get<RecapPayload>(buildShareRecapUrl(cardsBaseUrl, token));
  return response.data;
}
