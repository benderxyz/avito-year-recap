import { NavLink, Stack } from '@mantine/core';
import { Link, useRouterState } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';
import { routes } from '@/shared/lib/routes';

const NAV_ITEMS = [
  { to: routes.metrics, labelKey: 'nav.metrics' },
  { to: routes.badges, labelKey: 'nav.badges' },
  { to: routes.stories, labelKey: 'nav.stories' },
  { to: routes.recommendations, labelKey: 'nav.recommendations' },
  { to: routes.preview, labelKey: 'nav.preview' },
] as const;

export default function AppNav() {
  const { t } = useTranslation();
  const pathname = useRouterState({ select: (state) => state.location.pathname });

  return (
    <Stack gap="xs">
      {NAV_ITEMS.map((item) => (
        <NavLink
          key={item.to}
          component={Link}
          to={item.to}
          label={t(item.labelKey)}
          active={pathname === item.to || pathname.startsWith(`${item.to}/`)}
        />
      ))}
    </Stack>
  );
}
