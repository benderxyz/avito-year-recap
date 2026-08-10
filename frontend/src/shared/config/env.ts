type Env = {
  cardsBaseUrl: string;
  usersBaseUrl: string;
  recapYear: number;
};

export const env: Env = {
  cardsBaseUrl: import.meta.env.VITE_CARDS_BASE_URL ?? 'http://localhost:8081',
  usersBaseUrl: import.meta.env.VITE_USERS_BASE_URL ?? 'http://localhost:8082/users',
  recapYear: Number(import.meta.env.VITE_RECAP_YEAR ?? '2026'),
};
