import { createConnectTransport } from '@connectrpc/connect-web';

// Connect transport for the Go backend.
// In dev, Vite proxies /portal.v1.* to localhost:8080.
// In prod, this should point to the Go service URL (set via env var).
export const transport = createConnectTransport({
  baseUrl: import.meta.env.PUBLIC_API_URL || '/',
});
