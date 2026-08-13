import type { TFunction } from 'i18next';
import type { MetricDefinition } from '@/shared/api/generated/model/metricDefinition';
import type { DataTableColumn } from '@/shared/ui/DataTable/DataTable';
import { getMetricSourceFieldLabel, getMetricValueTypeLabel } from './metric-enum-labels';

export function getMetricColumns(t: TFunction): DataTableColumn<MetricDefinition>[] {
  return [
    { key: 'key', header: t('metrics.columns.key'), render: (row) => row.key },
    {
      key: 'valueType',
      header: t('metrics.columns.valueType'),
      render: (row) => getMetricValueTypeLabel(t, row.valueType),
    },
    {
      key: 'sourceField',
      header: t('metrics.columns.sourceField'),
      render: (row) => getMetricSourceFieldLabel(t, row.sourceField),
    },
    { key: 'sourceKey', header: t('metrics.columns.sourceKey'), render: (row) => row.sourceKey },
    {
      key: 'enabled',
      header: t('metrics.columns.enabled'),
      render: (row) => (row.enabled ? t('metrics.yes') : t('metrics.no')),
    },
    {
      key: 'isPublic',
      header: t('metrics.columns.isPublic'),
      render: (row) => (row.isPublic ? t('metrics.yes') : t('metrics.no')),
    },
    {
      key: 'includeInLlm',
      header: t('metrics.columns.includeInLlm'),
      render: (row) => (row.includeInLlm ? t('metrics.yes') : t('metrics.no')),
    },
  ];
}
