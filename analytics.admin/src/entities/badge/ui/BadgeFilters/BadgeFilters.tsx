import { Group } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { GetApiAdminBadgesVisibility } from '@/shared/api/generated/model/getApiAdminBadgesVisibility';
import SearchFilter from '@/shared/ui/filters/SearchFilter';
import SelectFilter from '@/shared/ui/filters/SelectFilter';
import TriStateFilter from '@/shared/ui/filters/TriStateFilter';
import { getBadgeVisibilityOptions } from '../../lib/enum-labels';

export type BadgeFiltersValue = {
  search: string;
  enabled: boolean | undefined;
  visibility: GetApiAdminBadgesVisibility | undefined;
  metric: string;
};

type BadgeFiltersProps = {
  value: BadgeFiltersValue;
  onChange: (value: BadgeFiltersValue) => void;
};

export default function BadgeFilters({ value, onChange }: BadgeFiltersProps) {
  const { t } = useTranslation();

  return (
    <Group align="flex-end">
      <SearchFilter
        label={t('badges.search')}
        value={value.search}
        placeholder={t('badges.search')}
        onChange={(search) => onChange({ ...value, search })}
      />
      <TriStateFilter
        label={t('badges.enabled')}
        value={value.enabled}
        onChange={(enabled) => onChange({ ...value, enabled })}
        trueLabel={t('badges.yes')}
        falseLabel={t('badges.no')}
        allLabel={t('badges.all')}
      />
      <SelectFilter
        label={t('badges.visibility')}
        value={value.visibility ?? null}
        onChange={(visibility) =>
          onChange({
            ...value,
            visibility: (visibility as GetApiAdminBadgesVisibility | null) ?? undefined,
          })
        }
        data={getBadgeVisibilityOptions(t)}
        placeholder={t('badges.all')}
      />
      <SearchFilter
        label={t('badges.metric')}
        value={value.metric}
        placeholder={t('badges.metric')}
        onChange={(metric) => onChange({ ...value, metric })}
      />
    </Group>
  );
}
