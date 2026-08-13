import { Group } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import SearchFilter from '@/shared/ui/filters/SearchFilter';
import TriStateFilter from '@/shared/ui/filters/TriStateFilter';

export type MetricFiltersValue = {
  search: string;
  enabled: boolean | undefined;
  isPublic: boolean | undefined;
  includeInLlm: boolean | undefined;
};

type MetricFiltersProps = {
  value: MetricFiltersValue;
  onChange: (value: MetricFiltersValue) => void;
};

export default function MetricFilters({ value, onChange }: MetricFiltersProps) {
  const { t } = useTranslation();

  return (
    <Group align="flex-end">
      <SearchFilter
        label={t('metrics.search')}
        value={value.search}
        placeholder={t('metrics.search')}
        onChange={(search) => onChange({ ...value, search })}
      />
      <TriStateFilter
        label={t('metrics.enabled')}
        value={value.enabled}
        onChange={(enabled) => onChange({ ...value, enabled })}
        trueLabel={t('metrics.yes')}
        falseLabel={t('metrics.no')}
        allLabel={t('metrics.all')}
      />
      <TriStateFilter
        label={t('metrics.isPublic')}
        value={value.isPublic}
        onChange={(isPublic) => onChange({ ...value, isPublic })}
        trueLabel={t('metrics.yes')}
        falseLabel={t('metrics.no')}
        allLabel={t('metrics.all')}
      />
      <TriStateFilter
        label={t('metrics.includeInLlm')}
        value={value.includeInLlm}
        onChange={(includeInLlm) => onChange({ ...value, includeInLlm })}
        trueLabel={t('metrics.yes')}
        falseLabel={t('metrics.no')}
        allLabel={t('metrics.all')}
      />
    </Group>
  );
}
