/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_CARDS_BASE_URL: string;
  readonly VITE_RECAP_YEAR: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
