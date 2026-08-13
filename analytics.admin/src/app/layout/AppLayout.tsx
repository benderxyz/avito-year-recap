import { Outlet } from '@tanstack/react-router';
import { AppShell } from '@mantine/core';

export default function AppLayout() {
  return (
    <AppShell
      padding="md"
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: 'sm',
      }}
    >
      <AppShell.Header>Logotype</AppShell.Header>
      <AppShell.Navbar>Navbar</AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
