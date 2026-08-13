import { AppShell, Burger, Group } from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { Outlet } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import AppNav from '@/app/layout/AppNav';

export default function AppLayout() {
  const { t } = useTranslation();
  const [opened, { toggle }] = useDisclosure();

  return (
    <AppShell
      padding="md"
      header={{ height: 60 }}
      navbar={{
        width: 260,
        breakpoint: 'sm',
        collapsed: { mobile: !opened },
      }}
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          {t('layout.logo')}
        </Group>
      </AppShell.Header>
      <AppShell.Navbar p="md">
        <AppNav />
      </AppShell.Navbar>
      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
