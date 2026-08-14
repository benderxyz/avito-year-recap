import { Button } from '@mantine/core';
import { useQueryClient } from '@tanstack/react-query';
import { Link, useNavigate } from '@tanstack/react-router';
import { useQueryStates } from 'nuqs';
import { useTranslation } from 'react-i18next';
import {
  getGetApiAdminRecommendationsQueryKey,
  getRecommendationColumns,
  RecommendationFilters,
  type RecommendationFiltersValue,
  recommendationFilterParsers,
  recommendationToFormValues,
  toRecommendationWrite,
  useGetApiAdminRecommendations,
  usePutApiAdminRecommendationsId,
} from '@/entities/recommendation';
import { routes } from '@/shared/lib/routes';
import { priorityBetween } from '@/shared/lib/sort-order';
import CatalogPage from '@/shared/ui/CatalogPage';
import DataTable from '@/shared/ui/DataTable';

export default function RecommendationsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const updateMutation = usePutApiAdminRecommendationsId();

  const [query, setQuery] = useQueryStates(recommendationFilterParsers, { history: 'replace' });

  const filters: RecommendationFiltersValue = {
    search: query.search,
    enabled: query.enabled ?? undefined,
    metric: query.metric,
    minPriority: query.minPriority ?? undefined,
  };

  const { data, isLoading, isError } = useGetApiAdminRecommendations({
    search: filters.search || undefined,
    enabled: filters.enabled,
    metric: filters.metric || undefined,
    minPriority: filters.minPriority,
  });

  function onFiltersChange(next: RecommendationFiltersValue) {
    setQuery({
      search: next.search,
      enabled: next.enabled ?? null,
      metric: next.metric,
      minPriority: next.minPriority ?? null,
    });
  }

  const items = data?.items ?? [];

  async function onReorder(nextItems: typeof items, movedId: string) {
    const movedIndex = nextItems.findIndex((row) => row.id === movedId);
    const moved = nextItems[movedIndex];
    if (!moved) {
      return;
    }

    const priority = priorityBetween(nextItems, movedIndex);
    if (priority === moved.priority) {
      return;
    }

    await updateMutation.mutateAsync({
      id: moved.id,
      data: { ...toRecommendationWrite(recommendationToFormValues(moved)), priority },
    });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminRecommendationsQueryKey() });
  }

  return (
    <CatalogPage
      title={t('recommendations.title')}
      actions={
        <Button component={Link} to={routes.recommendationCreate}>
          {t('recommendations.create')}
        </Button>
      }
      filters={<RecommendationFilters value={filters} onChange={onFiltersChange} />}
      isError={isError}
      errorMessage={t('recommendations.loadError')}
    >
      <DataTable
        columns={getRecommendationColumns(t)}
        items={items}
        emptyMessage={t('recommendations.empty')}
        getRowId={(row) => row.id}
        isLoading={isLoading}
        onReorder={(nextItems, movedId) => {
          void onReorder(nextItems, movedId);
        }}
        onRowClick={(row) => {
          navigate({ to: routes.recommendationById, params: { id: row.id } });
        }}
      />
    </CatalogPage>
  );
}
