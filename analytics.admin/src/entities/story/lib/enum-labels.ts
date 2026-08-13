import type { TFunction } from 'i18next';
import type { PredicateOp } from '@/shared/api/generated/model/predicateOp';
import { PredicateOp as PredicateOpValues } from '@/shared/api/generated/model/predicateOp';
import type { StoryCreateSceneType } from '@/shared/api/generated/model/storyCreateSceneType';
import { StoryCreateSceneType as StoryCreateSceneTypeValues } from '@/shared/api/generated/model/storyCreateSceneType';
import type { StoryCreateVisibility } from '@/shared/api/generated/model/storyCreateVisibility';
import { StoryCreateVisibility as StoryCreateVisibilityValues } from '@/shared/api/generated/model/storyCreateVisibility';

const VISIBILITY_KEYS = {
  private: 'stories.enums.visibility.private',
  public: 'stories.enums.visibility.public',
  both: 'stories.enums.visibility.both',
} as const satisfies Record<StoryCreateVisibility, string>;

const SCENE_TYPE_KEYS = {
  intro: 'stories.enums.sceneType.intro',
  stat: 'stories.enums.sceneType.stat',
  insight: 'stories.enums.sceneType.insight',
  achievement: 'stories.enums.sceneType.achievement',
  upsell: 'stories.enums.sceneType.upsell',
  blocks: 'stories.enums.sceneType.blocks',
  outro: 'stories.enums.sceneType.outro',
  custom: 'stories.enums.sceneType.custom',
} as const satisfies Record<StoryCreateSceneType, string>;

const PREDICATE_OP_KEYS = {
  gt: 'stories.enums.predicateOp.gt',
  gte: 'stories.enums.predicateOp.gte',
  eq: 'stories.enums.predicateOp.eq',
  exists: 'stories.enums.predicateOp.exists',
} as const satisfies Record<PredicateOp, string>;

export function getStoryVisibilityLabel(t: TFunction, value: StoryCreateVisibility) {
  return t(VISIBILITY_KEYS[value]);
}

export function getStorySceneTypeLabel(t: TFunction, value: StoryCreateSceneType) {
  return t(SCENE_TYPE_KEYS[value]);
}

export function getPredicateOpLabel(t: TFunction, value: PredicateOp) {
  return t(PREDICATE_OP_KEYS[value]);
}

export function getStoryVisibilityOptions(t: TFunction) {
  return Object.values(StoryCreateVisibilityValues).map((value) => ({
    value,
    label: getStoryVisibilityLabel(t, value),
  }));
}

export function getStorySceneTypeOptions(t: TFunction) {
  return Object.values(StoryCreateSceneTypeValues).map((value) => ({
    value,
    label: getStorySceneTypeLabel(t, value),
  }));
}

export function getPredicateOpOptions(t: TFunction) {
  return Object.values(PredicateOpValues).map((value) => ({
    value,
    label: getPredicateOpLabel(t, value),
  }));
}
