import { Button } from '@mantine/core';
import { Link, useNavigate } from '@tanstack/react-router';
import { useQueryStates } from 'nuqs';
import { useTranslation } from 'react-i18next';
import {
  getStoryColumns,
  StoryFilters,
  type StoryFiltersValue,
  storyFilterParsers,
  useGetApiAdminStories,
} from '@/entities/story';
import { routes } from '@/shared/lib/routes';
import CatalogPage from '@/shared/ui/CatalogPage';
import DataTable from '@/shared/ui/DataTable';

export default function StoriesPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const [query, setQuery] = useQueryStates(storyFilterParsers, { history: 'replace' });

  const filters: StoryFiltersValue = {
    search: query.search,
    enabled: query.enabled ?? undefined,
    visibility: query.visibility ?? undefined,
    sceneType: query.sceneType ?? undefined,
    metric: query.metric,
  };

  const { data, isLoading, isError } = useGetApiAdminStories({
    search: filters.search || undefined,
    enabled: filters.enabled,
    visibility: filters.visibility,
    sceneType: filters.sceneType,
    metric: filters.metric || undefined,
  });

  function onFiltersChange(next: StoryFiltersValue) {
    setQuery({
      search: next.search,
      enabled: next.enabled ?? null,
      visibility: next.visibility ?? null,
      sceneType: next.sceneType ?? null,
      metric: next.metric,
    });
  }

  return (
    <CatalogPage
      title={t('stories.title')}
      actions={
        <Button component={Link} to={routes.storyCreate}>
          {t('stories.create')}
        </Button>
      }
      filters={<StoryFilters value={filters} onChange={onFiltersChange} />}
      isError={isError}
      errorMessage={t('stories.loadError')}
    >
      <DataTable
        columns={getStoryColumns(t)}
        items={data?.items ?? []}
        emptyMessage={t('stories.empty')}
        getRowId={(row) => row.id}
        isLoading={isLoading}
        onRowClick={(row) => {
          navigate({ to: routes.storyById, params: { id: row.id } });
        }}
      />
    </CatalogPage>
  );
}
