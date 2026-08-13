import type { TFunction } from 'i18next';
import type { StoryRule } from '@/shared/api/generated/model/storyRule';
import type { DataTableColumn } from '@/shared/ui/DataTable/DataTable';
import { getStorySceneTypeLabel, getStoryVisibilityLabel } from './enum-labels';

export function getStoryColumns(t: TFunction): DataTableColumn<StoryRule>[] {
  return [
    { key: 'id', header: t('stories.columns.id'), render: (row) => row.id },
    {
      key: 'sceneType',
      header: t('stories.columns.sceneType'),
      render: (row) => getStorySceneTypeLabel(t, row.sceneType),
    },
    {
      key: 'visibility',
      header: t('stories.columns.visibility'),
      render: (row) => getStoryVisibilityLabel(t, row.visibility),
    },
    {
      key: 'enabled',
      header: t('stories.columns.enabled'),
      render: (row) => (row.enabled ? t('stories.yes') : t('stories.no')),
    },
    {
      key: 'metric',
      header: t('stories.columns.metric'),
      render: (row) => row.when?.metric ?? '',
    },
    { key: 'sortOrder', header: t('stories.columns.sortOrder'), render: (row) => row.sortOrder },
  ];
}
