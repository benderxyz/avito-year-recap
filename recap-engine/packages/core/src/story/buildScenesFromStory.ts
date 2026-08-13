import { metricDate, metricNumber } from '../payload/metrics';
import { defineScenes } from '../scenes/defineScenes';
import { ESceneActionType } from '../types/actions';
import type { RecapContext } from '../types/context';
import { EMotionPreset } from '../types/motion';
import { EMetricType, type RecapPayload } from '../types/payload';
import {
  ESceneBlockType,
  ESceneType,
  type SceneBlock,
  type SceneDefinition,
  type StatBlock,
} from '../types/scenes';
import type { StoryBlock, StoryItem, StoryStatBlock } from '../types/storyItems';

function narrative(payload: RecapPayload, id: string) {
  return payload.narrative?.scenes?.[id];
}

function formatMetricValue(
  ctx: RecapContext<RecapPayload>,
  metricKey: string,
  dateFormat: Intl.DateTimeFormatOptions | undefined,
): string | null {
  const metric = ctx.data.metrics[metricKey];

  if (!metric) return null;

  switch (metric.type) {
    case EMetricType.Money:
      return ctx.format.currency(metric.value, metric.currency);
    case EMetricType.Number:
    case EMetricType.Percentile:
    case EMetricType.Ratio:
      return ctx.format.number(metric.value);
    case EMetricType.Date: {
      const parsed = metricDate(ctx.data.metrics, metricKey);
      return parsed ? ctx.format.date(parsed, dateFormat) : null;
    }
    case EMetricType.String:
      return metric.value;
    default:
      return null;
  }
}

function fillValueTemplate(
  template: string,
  ctx: RecapContext<RecapPayload>,
  metricKey: string | undefined,
  dateFormat?: Intl.DateTimeFormatOptions,
): string {
  if (!metricKey || !template.includes('{{value}}')) return template;

  const formatted = formatMetricValue(ctx, metricKey, dateFormat);

  if (formatted === null) return template;

  return template.replace(/\{\{value\}\}/g, formatted);
}

function mapStatBlock(block: StoryStatBlock): StatBlock<RecapPayload> {
  const percentileKey = block.percentile;

  return {
    type: ESceneBlockType.Stat,
    title: block.title,
    eyebrow: block.eyebrow,
    valueFormat: block.valueFormat,
    blockMotion: block.blockMotion,
    unit: block.unit,
    value: (ctx) => metricNumber(ctx.data.metrics, block.value),
    comparison: percentileKey
      ? {
          template: block.comparisonTemplate ?? 'больше, чем у {{percentile}}% пользователей',
          percentile: (ctx) => metricNumber(ctx.data.metrics, percentileKey),
        }
      : undefined,
  };
}

function mapBlock(block: StoryBlock): SceneBlock<RecapPayload> {
  switch (block.type) {
    case ESceneBlockType.Stat:
      return mapStatBlock(block);
    case ESceneBlockType.Text:
      return { type: ESceneBlockType.Text, text: block.text, blockMotion: block.blockMotion };
    case ESceneBlockType.Callout:
      return {
        type: ESceneBlockType.Callout,
        text: block.text,
        blockMotion: block.blockMotion,
      };
    default: {
      const _exhaustive: never = block;
      return _exhaustive;
    }
  }
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
        title: (ctx) => {
          const raw = narrative(ctx.data, item.id)?.title ?? item.title ?? '';
          return fillValueTemplate(raw, ctx, item.value, item.dateFormat);
        },
        text: (ctx) => {
          const raw = narrative(ctx.data, item.id)?.body ?? item.text ?? '';
          return fillValueTemplate(raw, ctx, item.value, item.dateFormat);
        },
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
          return fillValueTemplate(raw, ctx, item.value, item.dateFormat);
        },
        callout: (ctx) => {
          const raw = item.callout ?? narrative(ctx.data, item.id)?.comparison ?? '';
          return fillValueTemplate(raw, ctx, item.value, item.dateFormat);
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
