export {};

declare global {
  interface ImportMetaEnv {
    readonly VITE_BASE_PATH?: string;
    readonly VITE_LOCALES_URL?: string;
    readonly VITE_API_BASE_URL?: string;
    readonly VITE_ADMIN_API_TOKEN?: string;
  }
}
