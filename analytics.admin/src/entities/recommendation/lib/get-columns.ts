import type { TFunction } from 'i18next';
import type { RecommendationRule } from '@/shared/api/generated/model/recommendationRule';
import type { DataTableColumn } from '@/shared/ui/DataTable/DataTable';

export function getRecommendationColumns(t: TFunction): DataTableColumn<RecommendationRule>[] {
  return [
    { key: 'id', header: t('recommendations.columns.id'), render: (row) => row.id },
    { key: 'title', header: t('recommendations.columns.title'), render: (row) => row.title },
    { key: 'path', header: t('recommendations.columns.path'), render: (row) => row.path },
    {
      key: 'enabled',
      header: t('recommendations.columns.enabled'),
      render: (row) => (row.enabled ? t('recommendations.yes') : t('recommendations.no')),
    },
    {
      key: 'metric',
      header: t('recommendations.columns.metric'),
      render: (row) => row.when.predicates[0]?.metric ?? '',
    },
  ];
}
