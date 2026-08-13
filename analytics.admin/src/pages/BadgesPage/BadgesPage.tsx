import { Button } from '@mantine/core';
import { Link, useNavigate } from '@tanstack/react-router';
import { useQueryStates } from 'nuqs';
import { useTranslation } from 'react-i18next';
import {
  BadgeFilters,
  type BadgeFiltersValue,
  badgeFilterParsers,
  getBadgeColumns,
  useGetApiAdminBadges,
} from '@/entities/badge';
import { routes } from '@/shared/lib/routes';
import CatalogPage from '@/shared/ui/CatalogPage';
import DataTable from '@/shared/ui/DataTable';

export default function BadgesPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [query, setQuery] = useQueryStates(badgeFilterParsers, { history: 'replace' });

  const filters: BadgeFiltersValue = {
    search: query.search,
    enabled: query.enabled ?? undefined,
    visibility: query.visibility ?? undefined,
    metric: query.metric,
  };

  const { data, isLoading, isError } = useGetApiAdminBadges({
    search: filters.search || undefined,
    enabled: filters.enabled,
    visibility: filters.visibility,
    metric: filters.metric || undefined,
  });

  function onFiltersChange(next: BadgeFiltersValue) {
    setQuery({
      search: next.search,
      enabled: next.enabled ?? null,
      visibility: next.visibility ?? null,
      metric: next.metric,
    });
  }

  return (
    <CatalogPage
      title={t('badges.title')}
      actions={
        <Button component={Link} to={routes.badgeCreate}>
          {t('badges.create')}
        </Button>
      }
      filters={<BadgeFilters value={filters} onChange={onFiltersChange} />}
      isError={isError}
      errorMessage={t('badges.loadError')}
    >
      <DataTable
        columns={getBadgeColumns(t)}
        items={data?.items ?? []}
        emptyMessage={t('badges.empty')}
        getRowId={(row) => row.id}
        isLoading={isLoading}
        onRowClick={(row) => {
          navigate({ to: routes.badgeById, params: { id: row.id } });
        }}
      />
    </CatalogPage>
  );
}
