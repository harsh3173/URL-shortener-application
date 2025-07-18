/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_NODE_ENV: string;
  readonly VITE_GOOGLE_ANALYTICS_ID: string;
  readonly VITE_ENABLE_ANALYTICS: string;
  readonly VITE_ENABLE_CUSTOM_DOMAINS: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}