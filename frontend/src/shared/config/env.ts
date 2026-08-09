type Env = {
  cardsBaseUrl: string;
  recapYear: number;
};

export const env: Env = {
  cardsBaseUrl: import.meta.env.VITE_CARDS_BASE_URL ?? 'http://localhost:8081',
  recapYear: Number(import.meta.env.VITE_RECAP_YEAR ?? '2026'),
};
