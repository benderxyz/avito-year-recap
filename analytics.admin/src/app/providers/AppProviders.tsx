import { QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode, useState } from 'react';
import { createQueryClient } from '@/shared/api/query-client';
import I18nProvider from '@/shared/i18n/I18nProvider';

import { MantineProvider, mantineTheme } from '@/shared/ui';

type AppProvidersProps = {
  children?: ReactNode;
};

export default function AppProviders({ children }: AppProvidersProps) {
  const [queryClient] = useState(createQueryClient);

  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider forceColorScheme="dark" theme={mantineTheme} defaultColorScheme="light">
        <I18nProvider>{children}</I18nProvider>
      </MantineProvider>
    </QueryClientProvider>
  );
}
