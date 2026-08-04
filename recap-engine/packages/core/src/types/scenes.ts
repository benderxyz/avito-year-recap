import type { SceneAction } from './actions';
import type { PluralForms, RecapContext, ValueOrFn } from './context';
import type { EMotionPreset, MotionConfig } from './motion';

export enum ESceneType {
  Intro = 'intro',
  Stat = 'stat',
  Insight = 'insight',
  Achievement = 'achievement',
  Upsell = 'upsell',
  Blocks = 'blocks',
  Outro = 'outro',
  Custom = 'custom',
}

export enum ESceneBlockType {
  Stat = 'stat',
  Text = 'text',
  Callout = 'callout',
}

export enum EBackgroundType {
  Color = 'color',
  Image = 'image',
  Gradient = 'gradient',
  Slot = 'slot',
}

export enum EImageFit {
  Cover = 'cover',
  Contain = 'contain',
}

export enum EBackgroundSlotName {
  Background = 'background',
}

type ColorSceneBackground = { type: EBackgroundType.Color; value: string };

type ImageSceneBackground = {
  type: EBackgroundType.Image;
  src: string;
  fit?: EImageFit;
};

type GradientSceneBackground = { type: EBackgroundType.Gradient; value: string };

type SlotSceneBackground = {
  type: EBackgroundType.Slot;
  name: EBackgroundSlotName;
};

export type SceneBackground =
  | ColorSceneBackground
  | ImageSceneBackground
  | GradientSceneBackground
  | SlotSceneBackground;

export type ComparisonConfig<D> = {
  template: string;
  percentile: ValueOrFn<number, D>;
};

export type SceneBase<_D = unknown> = {
  id: string;
  key?: string;
  background?: SceneBackground;
  motion?: EMotionPreset | MotionConfig;
  durationMs?: number;
  actions?: SceneAction[];
};

export type IntroScene<D> = SceneBase<D> & {
  type: ESceneType.Intro;
  title: ValueOrFn<string, D>;
  subtitle?: ValueOrFn<string, D>;
};

export type StatScene<D> = SceneBase<D> & {
  type: ESceneType.Stat;
  eyebrow?: ValueOrFn<string, D>;
  title?: ValueOrFn<string, D>;
  value: ValueOrFn<number, D>;
  valueFormat?: Intl.NumberFormatOptions;
  unit?: ValueOrFn<string | PluralForms, D>;
  comparison?: ComparisonConfig<D>;
  blockMotion?: EMotionPreset.CountUp | EMotionPreset.None;
};

export type InsightScene<D> = SceneBase<D> & {
  type: ESceneType.Insight;
  title?: ValueOrFn<string, D>;
  text: ValueOrFn<string, D>;
  linksTo?: string;
  blockMotion?: EMotionPreset.StaggerText | EMotionPreset.None;
};

export type UpsellScene<D> = SceneBase<D> & {
  type: ESceneType.Upsell;
  title?: ValueOrFn<string, D>;
  text: ValueOrFn<string, D>;
  callout?: ValueOrFn<string, D>;
  blockMotion?: EMotionPreset.CalloutIn | EMotionPreset.None;
};

export type AchievementScene<D> = SceneBase<D> & {
  type: ESceneType.Achievement;
  badgeId?: string;
  title?: ValueOrFn<string, D>;
  description?: ValueOrFn<string, D>;
  icon?: string;
  blockMotion?: EMotionPreset.BadgePop | EMotionPreset.None;
};

export type OutroScene<D> = SceneBase<D> & {
  type: ESceneType.Outro;
  title: ValueOrFn<string, D>;
  subtitle?: ValueOrFn<string, D>;
};

export type StatBlock<D> = {
  type: ESceneBlockType.Stat;
  eyebrow?: ValueOrFn<string, D>;
  title?: ValueOrFn<string, D>;
  value: ValueOrFn<number, D>;
  valueFormat?: Intl.NumberFormatOptions;
  unit?: ValueOrFn<string | PluralForms, D>;
  comparison?: ComparisonConfig<D>;
  blockMotion?: EMotionPreset.CountUp | EMotionPreset.None;
};

export type TextBlock<D> = {
  type: ESceneBlockType.Text;
  text: ValueOrFn<string, D>;
  blockMotion?: EMotionPreset.StaggerText | EMotionPreset.None;
};

export type CalloutBlock<D> = {
  type: ESceneBlockType.Callout;
  text: ValueOrFn<string, D>;
  blockMotion?: EMotionPreset.CalloutIn | EMotionPreset.None;
};

export type SceneBlock<D> = StatBlock<D> | TextBlock<D> | CalloutBlock<D>;

export type BlocksScene<D> = SceneBase<D> & {
  type: ESceneType.Blocks;
  blocks: SceneBlock<D>[];
};

export type CustomScene<D> = SceneBase<D> & {
  type: ESceneType.Custom;
  sceneType?: string;
  props?: Record<string, unknown>;
  render?: (ctx: RecapContext<D> & { goNext: () => void; goPrev: () => void }) => unknown;
};

export type SceneDefinition<D = unknown> =
  | IntroScene<D>
  | StatScene<D>
  | InsightScene<D>
  | UpsellScene<D>
  | AchievementScene<D>
  | OutroScene<D>
  | BlocksScene<D>
  | CustomScene<D>;
