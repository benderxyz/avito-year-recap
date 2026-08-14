import { Group } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import NumberFilter from '@/shared/ui/filters/NumberFilter';
import SearchFilter from '@/shared/ui/filters/SearchFilter';
import TriStateFilter from '@/shared/ui/filters/TriStateFilter';

export type RecommendationFiltersValue = {
  search: string;
  enabled: boolean | undefined;
  metric: string;
  minPriority: number | undefined;
};

type RecommendationFiltersProps = {
  value: RecommendationFiltersValue;
  onChange: (value: RecommendationFiltersValue) => void;
};

export default function RecommendationFilters({ value, onChange }: RecommendationFiltersProps) {
  const { t } = useTranslation();

  return (
    <Group align="flex-end">
      <SearchFilter
        label={t('recommendations.search')}
        value={value.search}
        placeholder={t('recommendations.search')}
        onChange={(search) => onChange({ ...value, search })}
      />
      <TriStateFilter
        label={t('recommendations.enabled')}
        value={value.enabled}
        onChange={(enabled) => onChange({ ...value, enabled })}
        trueLabel={t('recommendations.yes')}
        falseLabel={t('recommendations.no')}
        allLabel={t('recommendations.all')}
      />
      <SearchFilter
        label={t('recommendations.metric')}
        value={value.metric}
        placeholder={t('recommendations.metric')}
        onChange={(metric) => onChange({ ...value, metric })}
      />
      <NumberFilter
        label={t('recommendations.minPriority')}
        value={value.minPriority ?? null}
        onChange={(minPriority) => onChange({ ...value, minPriority: minPriority ?? undefined })}
      />
    </Group>
  );
}
