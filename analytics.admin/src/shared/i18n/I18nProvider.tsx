import { Center, Loader } from '@mantine/core';
import { type ReactNode, Suspense } from 'react';
import { I18nextProvider } from 'react-i18next';
import { i18n } from './config';

type I18nProviderProps = {
  children?: ReactNode;
};

export default function I18nProvider({ children }: I18nProviderProps) {
  return (
    <I18nextProvider i18n={i18n}>
      <Suspense
        fallback={
          <Center h="100vh">
            <Loader />
          </Center>
        }
      >
        {children}
      </Suspense>
    </I18nextProvider>
  );
}
