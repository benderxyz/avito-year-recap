import { metricNumber } from '../payload/metrics';
import { defineScenes } from '../scenes/defineScenes';
import { ESceneActionType } from '../types/actions';
import { EMotionPreset } from '../types/motion';
import type { RecapPayload } from '../types/payload';
import {
  ESceneBlockType,
  ESceneType,
  type SceneBlock,
  type SceneDefinition,
  type StatBlock,
} from '../types/scenes';
import type { StoryBlock, StoryItem, StoryStatBlock } from '../types/story-items';

function narrative(payload: RecapPayload, id: string) {
  return payload.narrative?.scenes?.[id];
}

function fillMoneyTemplate(
  template: string,
  ctx: { data: RecapPayload; format: { currency: (n: number) => string } },
  metricKey: string | undefined,
): string {
  if (!metricKey || !template.includes('{{value}}')) return template;
  const amount = metricNumber(ctx.data.metrics, metricKey);
  return template.replace(/\{\{value\}\}/g, ctx.format.currency(amount));
}

function mapStatBlock(block: StoryStatBlock): StatBlock<RecapPayload> {
  return {
    type: ESceneBlockType.Stat,
    title: block.title,
    eyebrow: block.eyebrow,
    valueFormat: block.valueFormat,
    blockMotion: block.blockMotion,
    unit: block.unit,
    value: (ctx) => metricNumber(ctx.data.metrics, block.value),
    comparison: block.percentile
      ? {
          template: block.comparisonTemplate ?? 'больше, чем у {{percentile}}% пользователей',
          percentile: (ctx) => metricNumber(ctx.data.metrics, block.percentile as string),
        }
      : undefined,
  };
}

function mapBlock(block: StoryBlock): SceneBlock<RecapPayload> {
  if (block.type === ESceneBlockType.Stat) return mapStatBlock(block);
  if (block.type === ESceneBlockType.Text) {
    return { type: ESceneBlockType.Text, text: block.text, blockMotion: block.blockMotion };
  }
  return {
    type: ESceneBlockType.Callout,
    text: block.text,
    blockMotion: block.blockMotion,
  };
}

function mapStoryItem(item: StoryItem, payload: RecapPayload): SceneDefinition<RecapPayload> {
  const base = {
    id: item.id,
    motion: item.motion,
    durationMs: item.durationMs,
    actions: item.actions,
  };

  switch (item.type) {
    case ESceneType.Intro:
      return {
        ...base,
        type: ESceneType.Intro,
        title: (ctx) =>
          narrative(ctx.data, item.id)?.title ??
          item.title ??
          `${ctx.data.meta.user.displayName}, ваш ${ctx.data.meta.year}`,
        subtitle: (ctx) =>
          narrative(ctx.data, item.id)?.subtitle ?? item.subtitle ?? 'Год на Авито в цифрах',
        actions: item.actions ?? [{ type: ESceneActionType.Next, label: 'Начать' }],
      };

    case ESceneType.Stat: {
      const percentileKey = item.percentile;
      const comparisonTemplate =
        narrative(payload, item.id)?.comparison ??
        item.comparisonTemplate ??
        'это больше, чем у {{percentile}}% пользователей';

      return {
        ...base,
        type: ESceneType.Stat,
        blockMotion: item.blockMotion ?? EMotionPreset.CountUp,
        valueFormat: item.valueFormat,
        unit: item.unit,
        eyebrow:
          item.eyebrow !== undefined ? item.eyebrow : (ctx) => `За ${ctx.data.meta.year} год`,
        title: (ctx) => narrative(ctx.data, item.id)?.title ?? item.title ?? '',
        value: (ctx) => metricNumber(ctx.data.metrics, item.value),
        comparison: percentileKey
          ? {
              template: comparisonTemplate,
              percentile: (ctx) => metricNumber(ctx.data.metrics, percentileKey),
            }
          : undefined,
      };
    }

    case ESceneType.Insight:
      return {
        ...base,
        type: ESceneType.Insight,
        linksTo: item.linksTo,
        blockMotion: item.blockMotion ?? EMotionPreset.StaggerText,
        title: (ctx) => narrative(ctx.data, item.id)?.title ?? item.title ?? '',
        text: (ctx) => narrative(ctx.data, item.id)?.body ?? item.text ?? '',
      };

    case ESceneType.Achievement:
      return {
        ...base,
        type: ESceneType.Achievement,
        badgeId: item.badgeId,
        title: item.title,
        description: item.description,
        icon: item.icon,
        blockMotion: item.blockMotion ?? EMotionPreset.BadgePop,
      };

    case ESceneType.Upsell:
      return {
        ...base,
        type: ESceneType.Upsell,
        blockMotion: item.blockMotion ?? EMotionPreset.CalloutIn,
        title: (ctx) => narrative(ctx.data, item.id)?.title ?? item.title ?? '',
        text: (ctx) => {
          const raw = narrative(ctx.data, item.id)?.body ?? item.text ?? '';
          return fillMoneyTemplate(raw, ctx, item.value);
        },
        callout: (ctx) => {
          const raw = item.callout ?? narrative(ctx.data, item.id)?.comparison ?? '';
          return fillMoneyTemplate(raw, ctx, item.value);
        },
      };

    case ESceneType.Outro:
      return {
        ...base,
        type: ESceneType.Outro,
        title: (ctx) =>
          narrative(ctx.data, item.id)?.title ?? item.title ?? 'Это был ваш год на Авито',
        subtitle: (ctx) =>
          narrative(ctx.data, item.id)?.subtitle ??
          item.subtitle ??
          'Поделитесь итогами или загляните на главную',
      };

    case ESceneType.Blocks:
      return {
        ...base,
        type: ESceneType.Blocks,
        blocks: item.blocks.map(mapBlock),
      };

    case ESceneType.Custom:
      return {
        ...base,
        type: ESceneType.Custom,
        sceneType: item.sceneType,
        props: item.props,
      };

    default: {
      const _exhaustive: never = item;
      return _exhaustive;
    }
  }
}

/** Turn backend `story` playlist into SceneDefinition[] bound to RecapPayload as `data`. */
export function buildScenesFromStory(payload: RecapPayload): SceneDefinition<RecapPayload>[] {
  return defineScenes(payload.story.map((item) => mapStoryItem(item, payload)));
}
