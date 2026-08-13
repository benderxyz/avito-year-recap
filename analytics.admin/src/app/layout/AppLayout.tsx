import { AppShell } from '@mantine/core';
import { Outlet } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';

export default function AppLayout() {
  const { t } = useTranslation();

  return (
    <AppShell
      padding="md"
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: 'sm',
      }}
    >
      <AppShell.Header>{t('layout.logo')}</AppShell.Header>
      <AppShell.Navbar>{t('layout.navbar')}</AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
