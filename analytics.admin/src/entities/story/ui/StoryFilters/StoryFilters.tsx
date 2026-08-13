import { Group } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { GetApiAdminStoriesSceneType } from '@/shared/api/generated/model/getApiAdminStoriesSceneType';
import type { GetApiAdminStoriesVisibility } from '@/shared/api/generated/model/getApiAdminStoriesVisibility';
import SearchFilter from '@/shared/ui/filters/SearchFilter';
import SelectFilter from '@/shared/ui/filters/SelectFilter';
import TriStateFilter from '@/shared/ui/filters/TriStateFilter';
import { getStorySceneTypeOptions, getStoryVisibilityOptions } from '../../lib/enum-labels';

export type StoryFiltersValue = {
  search: string;
  enabled: boolean | undefined;
  visibility: GetApiAdminStoriesVisibility | undefined;
  sceneType: GetApiAdminStoriesSceneType | undefined;
  metric: string;
};

type StoryFiltersProps = {
  value: StoryFiltersValue;
  onChange: (value: StoryFiltersValue) => void;
};

export default function StoryFilters({ value, onChange }: StoryFiltersProps) {
  const { t } = useTranslation();

  return (
    <Group align="flex-end">
      <SearchFilter
        label={t('stories.search')}
        value={value.search}
        placeholder={t('stories.search')}
        onChange={(search) => onChange({ ...value, search })}
      />
      <TriStateFilter
        label={t('stories.enabled')}
        value={value.enabled}
        onChange={(enabled) => onChange({ ...value, enabled })}
        trueLabel={t('stories.yes')}
        falseLabel={t('stories.no')}
        allLabel={t('stories.all')}
      />
      <SelectFilter
        label={t('stories.visibility')}
        value={value.visibility ?? null}
        onChange={(visibility) =>
          onChange({
            ...value,
            visibility: (visibility as GetApiAdminStoriesVisibility | null) ?? undefined,
          })
        }
        data={getStoryVisibilityOptions(t)}
        placeholder={t('stories.all')}
      />
      <SelectFilter
        label={t('stories.sceneType')}
        value={value.sceneType ?? null}
        onChange={(sceneType) =>
          onChange({
            ...value,
            sceneType: (sceneType as GetApiAdminStoriesSceneType | null) ?? undefined,
          })
        }
        data={getStorySceneTypeOptions(t)}
        placeholder={t('stories.all')}
      />
      <SearchFilter
        label={t('stories.metric')}
        value={value.metric}
        placeholder={t('stories.metric')}
        onChange={(metric) => onChange({ ...value, metric })}
      />
    </Group>
  );
}
