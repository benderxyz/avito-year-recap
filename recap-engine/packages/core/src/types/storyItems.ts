import type { SceneAction } from './actions';
import type { PluralForms } from './context';
import type { EMotionPreset, MotionConfig } from './motion';
import type { ESceneBlockType, ESceneType } from './scenes';

export type StoryItemBase = {
  id: string;
  type: ESceneType;
  motion?: EMotionPreset | MotionConfig;
  durationMs?: number;
  actions?: SceneAction[];
};

export type StoryIntroItem = StoryItemBase & {
  type: ESceneType.Intro;
  title?: string;
  subtitle?: string;
};

export type StoryStatItem = StoryItemBase & {
  type: ESceneType.Stat;
  value: string;
  percentile?: string;
  unit?: string | PluralForms;
  title?: string;
  eyebrow?: string;
  comparisonTemplate?: string;
  valueFormat?: Intl.NumberFormatOptions;
  blockMotion?: EMotionPreset.CountUp | EMotionPreset.None;
};

export type StoryInsightItem = StoryItemBase & {
  type: ESceneType.Insight;
  text?: string;
  title?: string;
  value?: string;
  dateFormat?: Intl.DateTimeFormatOptions;
  linksTo?: string;
  blockMotion?: EMotionPreset.StaggerText | EMotionPreset.None;
};

export type StoryAchievementItem = StoryItemBase & {
  type: ESceneType.Achievement;
  badgeId?: string;
  title?: string;
  description?: string;
  icon?: string;
  blockMotion?: EMotionPreset.BadgePop | EMotionPreset.None;
};

export type StoryUpsellItem = StoryItemBase & {
  type: ESceneType.Upsell;
  title?: string;
  text?: string;
  callout?: string;
  value?: string;
  dateFormat?: Intl.DateTimeFormatOptions;
  blockMotion?: EMotionPreset.CalloutIn | EMotionPreset.None;
};

export type StoryOutroItem = StoryItemBase & {
  type: ESceneType.Outro;
  title?: string;
  subtitle?: string;
};

export type StoryStatBlock = {
  type: ESceneBlockType.Stat;
  value: string;
  percentile?: string;
  unit?: string | PluralForms;
  title?: string;
  eyebrow?: string;
  comparisonTemplate?: string;
  valueFormat?: Intl.NumberFormatOptions;
  blockMotion?: EMotionPreset.CountUp | EMotionPreset.None;
};

export type StoryTextBlock = {
  type: ESceneBlockType.Text;
  text: string;
  blockMotion?: EMotionPreset.StaggerText | EMotionPreset.None;
};

export type StoryCalloutBlock = {
  type: ESceneBlockType.Callout;
  text: string;
  blockMotion?: EMotionPreset.CalloutIn | EMotionPreset.None;
};

export type StoryBlock = StoryStatBlock | StoryTextBlock | StoryCalloutBlock;

export type StoryBlocksItem = StoryItemBase & {
  type: ESceneType.Blocks;
  blocks: StoryBlock[];
};

export type StoryCustomItem = StoryItemBase & {
  type: ESceneType.Custom;
  sceneType?: string;
  props?: Record<string, unknown>;
};

export type StoryItem =
  | StoryIntroItem
  | StoryStatItem
  | StoryInsightItem
  | StoryAchievementItem
  | StoryUpsellItem
  | StoryOutroItem
  | StoryBlocksItem
  | StoryCustomItem;
