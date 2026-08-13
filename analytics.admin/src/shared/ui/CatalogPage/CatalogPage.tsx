import { Alert, Group, Stack, Title } from '@mantine/core';
import type { ReactNode } from 'react';

type CatalogPageProps = {
  title: ReactNode;
  actions?: ReactNode;
  filters?: ReactNode;
  isError?: boolean;
  errorMessage?: string;
  children: ReactNode;
};

export default function CatalogPage({
  title,
  actions,
  filters,
  isError = false,
  errorMessage = 'Failed to load',
  children,
}: CatalogPageProps) {
  return (
    <Stack gap="md">
      <Group justify="space-between" align="center">
        <Title order={2}>{title}</Title>
        {actions}
      </Group>
      {filters}
      {isError ? <Alert color="red">{errorMessage}</Alert> : null}
      {children}
    </Stack>
  );
}
