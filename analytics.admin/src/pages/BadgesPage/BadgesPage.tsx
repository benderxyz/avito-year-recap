import { Button } from '@mantine/core';
import { useQueryClient } from '@tanstack/react-query';
import { Link, useNavigate } from '@tanstack/react-router';
import { useQueryStates } from 'nuqs';
import { useTranslation } from 'react-i18next';
import {
  BadgeFilters,
  type BadgeFiltersValue,
  badgeFilterParsers,
  badgeToFormValues,
  getBadgeColumns,
  getGetApiAdminBadgesQueryKey,
  toBadgeWrite,
  useGetApiAdminBadges,
  usePutApiAdminBadgesId,
} from '@/entities/badge';
import { routes } from '@/shared/lib/routes';
import { sortOrderBetween } from '@/shared/lib/sort-order';
import CatalogPage from '@/shared/ui/CatalogPage';
import DataTable from '@/shared/ui/DataTable';

export default function BadgesPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const updateMutation = usePutApiAdminBadgesId();

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

  const items = data?.items ?? [];

  async function onReorder(nextItems: typeof items, movedId: string) {
    const movedIndex = nextItems.findIndex((row) => row.id === movedId);
    const moved = nextItems[movedIndex];
    if (!moved) {
      return;
    }

    const sortOrder = sortOrderBetween(nextItems, movedIndex);
    if (sortOrder === moved.sortOrder) {
      return;
    }

    await updateMutation.mutateAsync({
      id: moved.id,
      data: { ...toBadgeWrite(badgeToFormValues(moved)), sortOrder },
    });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminBadgesQueryKey() });
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
        items={items}
        emptyMessage={t('badges.empty')}
        getRowId={(row) => row.id}
        isLoading={isLoading}
        onReorder={(items, movedId) => {
          void onReorder(items, movedId);
        }}
        onRowClick={(row) => {
          navigate({ to: routes.badgeById, params: { id: row.id } });
        }}
      />
    </CatalogPage>
  );
}
