import { defineConfig } from 'orval';

export default defineConfig({
  admin: {
    input: {
      target: 'https://recaps.hakolr.dev/api/admin/openapi.json',
    },
    output: {
      mode: 'tags-split',
      target: './src/shared/api/generated',
      schemas: './src/shared/api/generated/model',
      client: 'react-query',
      httpClient: 'fetch',
      clean: true,
      override: {
        mutator: {
          path: './src/shared/api/http.ts',
          name: 'customFetch',
        },
        fetch: {
          includeHttpResponseReturnType: false,
        },
      },
    },
  },
});
