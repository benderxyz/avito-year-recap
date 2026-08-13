import { Button } from '@mantine/core';
import { Link, useNavigate } from '@tanstack/react-router';
import { useQueryStates } from 'nuqs';
import { useTranslation } from 'react-i18next';
import { getMetricColumns } from '@/entities/metric/get-metric-columns';
import MetricFilters, { type MetricFiltersValue } from '@/entities/metric/MetricFilters';
import { metricFilterParsers } from '@/entities/metric/metric-filter-parsers';
import { useGetApiAdminMetrics } from '@/shared/api/generated/metrics/metrics';
import { routes } from '@/shared/lib/routes';
import CatalogPage from '@/shared/ui/CatalogPage';
import DataTable from '@/shared/ui/DataTable';

export default function MetricsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

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
        items={data?.items ?? []}
        emptyMessage={t('metrics.empty')}
        getRowId={(row) => row.key}
        isLoading={isLoading}
        onRowClick={(row) => {
          navigate({ to: routes.metricByKey, params: { key: row.key } });
        }}
      />
    </CatalogPage>
  );
}
