import { Button } from '@mantine/core';
import { useQueryClient } from '@tanstack/react-query';
import { Link, useNavigate } from '@tanstack/react-router';
import { useQueryStates } from 'nuqs';
import { useTranslation } from 'react-i18next';
import {
  getGetApiAdminMetricsQueryKey,
  getMetricColumns,
  MetricFilters,
  type MetricFiltersValue,
  metricFilterParsers,
  metricToFormValues,
  toMetricWrite,
  useGetApiAdminMetrics,
  usePutApiAdminMetricsKey,
} from '@/entities/metric';
import { routes } from '@/shared/lib/routes';
import { sortOrderBetween } from '@/shared/lib/sort-order';
import CatalogPage from '@/shared/ui/CatalogPage';
import DataTable from '@/shared/ui/DataTable';

export default function MetricsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const updateMutation = usePutApiAdminMetricsKey();

  const [query, setQuery] = useQueryStates(metricFilterParsers, { history: 'replace' });

  const filters: MetricFiltersValue = {
    search: query.search,
    enabled: query.enabled ?? undefined,
    isPublic: query.isPublic ?? undefined,
    includeInLlm: query.includeInLlm ?? undefined,
  };

  const { data, isLoading, isError } = useGetApiAdminMetrics({
    search: filters.search || undefined,
    enabled: filters.enabled,
    isPublic: filters.isPublic,
    includeInLlm: filters.includeInLlm,
  });

  function onFiltersChange(next: MetricFiltersValue) {
    setQuery({
      search: next.search,
      enabled: next.enabled ?? null,
      isPublic: next.isPublic ?? null,
      includeInLlm: next.includeInLlm ?? null,
    });
  }

  const items = data?.items ?? [];

  async function onReorder(nextItems: typeof items, movedId: string) {
    const movedIndex = nextItems.findIndex((row) => row.key === movedId);
    const moved = nextItems[movedIndex];
    if (!moved) {
      return;
    }

    const sortOrder = sortOrderBetween(nextItems, movedIndex);
    if (sortOrder === moved.sortOrder) {
      return;
    }

    await updateMutation.mutateAsync({
      key: moved.key,
      data: { ...toMetricWrite(metricToFormValues(moved)), sortOrder },
    });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminMetricsQueryKey() });
  }

  return (
    <CatalogPage
      title={t('metrics.title')}
      actions={
        <Button component={Link} to={routes.metricCreate}>
          {t('metrics.create')}
        </Button>
      }
      filters={<MetricFilters value={filters} onChange={onFiltersChange} />}
      isError={isError}
      errorMessage={t('metrics.loadError')}
    >
      <DataTable
        columns={getMetricColumns(t)}
        items={items}
        emptyMessage={t('metrics.empty')}
        getRowId={(row) => row.key}
        isLoading={isLoading}
        onReorder={(items, movedId) => {
          void onReorder(items, movedId);
        }}
        onRowClick={(row) => {
          navigate({ to: routes.metricByKey, params: { key: row.key } });
        }}
      />
    </CatalogPage>
  );
}
