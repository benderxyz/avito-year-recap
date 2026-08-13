import { parseAsBoolean, parseAsString, parseAsStringEnum } from 'nuqs';
import { GetApiAdminStoriesSceneType } from '@/shared/api/generated/model/getApiAdminStoriesSceneType';
import { GetApiAdminStoriesVisibility } from '@/shared/api/generated/model/getApiAdminStoriesVisibility';

export const storyFilterParsers = {
  search: parseAsString.withDefault(''),
  enabled: parseAsBoolean,
  visibility: parseAsStringEnum([
    GetApiAdminStoriesVisibility.private,
    GetApiAdminStoriesVisibility.public,
    GetApiAdminStoriesVisibility.both,
  ]),
  sceneType: parseAsStringEnum([
    GetApiAdminStoriesSceneType.intro,
    GetApiAdminStoriesSceneType.stat,
    GetApiAdminStoriesSceneType.insight,
    GetApiAdminStoriesSceneType.achievement,
    GetApiAdminStoriesSceneType.upsell,
    GetApiAdminStoriesSceneType.blocks,
    GetApiAdminStoriesSceneType.outro,
    GetApiAdminStoriesSceneType.custom,
  ]),
  metric: parseAsString.withDefault(''),
};
