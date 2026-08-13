import type { TFunction } from 'i18next';
import type { BadgeRule } from '@/shared/api/generated/model/badgeRule';
import type { DataTableColumn } from '@/shared/ui/DataTable/DataTable';
import { getBadgeVisibilityLabel } from './enum-labels';

export function getBadgeColumns(t: TFunction): DataTableColumn<BadgeRule>[] {
  return [
    { key: 'id', header: t('badges.columns.id'), render: (row) => row.id },
    { key: 'title', header: t('badges.columns.title'), render: (row) => row.title },
    {
      key: 'visibility',
      header: t('badges.columns.visibility'),
      render: (row) => getBadgeVisibilityLabel(t, row.visibility),
    },
    {
      key: 'enabled',
      header: t('badges.columns.enabled'),
      render: (row) => (row.enabled ? t('badges.yes') : t('badges.no')),
    },
    { key: 'metric', header: t('badges.columns.metric'), render: (row) => row.when.metric },
    { key: 'sortOrder', header: t('badges.columns.sortOrder'), render: (row) => row.sortOrder },
  ];
}
