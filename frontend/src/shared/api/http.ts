import axios from 'axios';

export const httpClient = axios.create({
  timeout: 15_000,
});

export function buildRecapUrl(cardsBaseUrl: string, year: number, userId: string): string {
  return `${cardsBaseUrl.replace(/\/$/, '')}/api/recap/${year}/${userId}`;
}

export function buildShareRecapUrl(cardsBaseUrl: string, token: string): string {
  const normalizedToken = token.trim().split(/\s+/, 1)[0];
  return `${cardsBaseUrl.replace(/\/$/, '')}/api/share/${encodeURIComponent(normalizedToken)}`;
}

export const DUMMYJSON_BASE_URL = 'https://dummyjson.com';
