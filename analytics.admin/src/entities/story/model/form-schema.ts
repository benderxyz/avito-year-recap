import { z } from 'zod';
import type { Predicate } from '@/shared/api/generated/model/predicate';
import { PredicateOp } from '@/shared/api/generated/model/predicateOp';
import type { StoryCreate } from '@/shared/api/generated/model/storyCreate';
import type { StoryCreatePayload } from '@/shared/api/generated/model/storyCreatePayload';
import { StoryCreateSceneType } from '@/shared/api/generated/model/storyCreateSceneType';
import { StoryCreateVisibility } from '@/shared/api/generated/model/storyCreateVisibility';
import type { StoryRule } from '@/shared/api/generated/model/storyRule';
import type { StoryWrite } from '@/shared/api/generated/model/storyWrite';

const predicateOpEnum = z.enum([
  PredicateOp.gt,
  PredicateOp.gte,
  PredicateOp.eq,
  PredicateOp.exists,
]);

export const storyFormSchema = z.object({
  id: z.string().min(1),
  sceneType: z.enum([
    StoryCreateSceneType.intro,
    StoryCreateSceneType.stat,
    StoryCreateSceneType.insight,
    StoryCreateSceneType.achievement,
    StoryCreateSceneType.upsell,
    StoryCreateSceneType.blocks,
    StoryCreateSceneType.outro,
    StoryCreateSceneType.custom,
  ]),
  visibility: z.enum([
    StoryCreateVisibility.private,
    StoryCreateVisibility.public,
    StoryCreateVisibility.both,
  ]),
  payload: z.string().min(1).refine(isJsonObject, { message: 'Invalid JSON object' }),
  sortOrder: z.number(),
  enabled: z.boolean(),
  when: z.object({
    metric: z.string(),
    op: predicateOpEnum,
    value: z.number().nullable(),
  }),
});

export type StoryFormValues = z.infer<typeof storyFormSchema>;

function isJsonObject(value: string) {
  try {
    const parsed: unknown = JSON.parse(value);
    return typeof parsed === 'object' && parsed !== null && !Array.isArray(parsed);
  } catch {
    return false;
  }
}

function stringifyPayload(payload: Record<string, unknown>) {
  return JSON.stringify(payload, null, 2);
}

export const storyFormDefaults: StoryFormValues = {
  id: '',
  sceneType: StoryCreateSceneType.intro,
  visibility: StoryCreateVisibility.private,
  payload: stringifyPayload({ id: '', type: StoryCreateSceneType.intro }),
  sortOrder: 0,
  enabled: true,
  when: {
    metric: '',
    op: PredicateOp.gte,
    value: null,
  },
};

export function storyToFormValues(story: StoryRule): StoryFormValues {
  return {
    id: story.id,
    sceneType: story.sceneType,
    visibility: story.visibility,
    payload: stringifyPayload(story.payload),
    sortOrder: story.sortOrder,
    enabled: story.enabled,
    when: {
      metric: story.when?.metric ?? '',
      op: story.when?.op ?? PredicateOp.gte,
      value: story.when?.value ?? null,
    },
  };
}

function parsePayload(value: string): StoryCreatePayload {
  return JSON.parse(value) as StoryCreatePayload;
}

function toOptionalPredicate(when: StoryFormValues['when']): Predicate | null {
  if (!when.metric) {
    return null;
  }

  if (when.op === PredicateOp.exists || when.value === null) {
    return { metric: when.metric, op: when.op };
  }

  return { metric: when.metric, op: when.op, value: when.value };
}

export function toStoryCreate(values: StoryFormValues): StoryCreate {
  return {
    id: values.id,
    sceneType: values.sceneType,
    visibility: values.visibility,
    payload: {
      ...parsePayload(values.payload),
      id: values.id,
      type: values.sceneType,
    },
    sortOrder: values.sortOrder,
    enabled: values.enabled,
    when: toOptionalPredicate(values.when),
  };
}

export function toStoryWrite(values: StoryFormValues): StoryWrite {
  const created = toStoryCreate(values);
  return {
    sceneType: created.sceneType,
    visibility: created.visibility,
    payload: created.payload,
    sortOrder: created.sortOrder,
    enabled: created.enabled,
    when: created.when,
  };
}
