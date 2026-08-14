import { z } from 'zod';
import type { GroupWhen } from '@/shared/api/generated/model/groupWhen';
import { GroupWhenMatch } from '@/shared/api/generated/model/groupWhenMatch';
import type { Predicate } from '@/shared/api/generated/model/predicate';
import { PredicateOp } from '@/shared/api/generated/model/predicateOp';
import type { RecommendationCreate } from '@/shared/api/generated/model/recommendationCreate';
import type { RecommendationRule } from '@/shared/api/generated/model/recommendationRule';
import type { RecommendationWrite } from '@/shared/api/generated/model/recommendationWrite';

const predicateOpEnum = z.enum([
  PredicateOp.gt,
  PredicateOp.gte,
  PredicateOp.eq,
  PredicateOp.exists,
]);

const predicateSchema = z.object({
  metric: z.string().min(1),
  op: predicateOpEnum,
  value: z.number().nullable(),
});

export const recommendationFormSchema = z.object({
  id: z.string().min(1),
  title: z.string().min(1),
  text: z.string().min(1),
  callout: z.string(),
  linkLabel: z.string().min(1),
  path: z.string().min(1),
  enabled: z.boolean(),
  priority: z.number(),
  when: z.object({
    match: z.enum([GroupWhenMatch.all, GroupWhenMatch.any]),
    predicates: z.array(predicateSchema),
  }),
});

export type RecommendationFormValues = z.infer<typeof recommendationFormSchema>;

export const emptyPredicate: RecommendationFormValues['when']['predicates'][number] = {
  metric: '',
  op: PredicateOp.gte,
  value: null,
};

export const recommendationFormDefaults: RecommendationFormValues = {
  id: '',
  title: '',
  text: '',
  callout: '',
  linkLabel: '',
  path: '',
  enabled: true,
  priority: 0,
  when: {
    match: GroupWhenMatch.all,
    predicates: [],
  },
};

export function recommendationToFormValues(
  recommendation: RecommendationRule,
): RecommendationFormValues {
  return {
    id: recommendation.id,
    title: recommendation.title,
    text: recommendation.text,
    callout: recommendation.callout,
    linkLabel: recommendation.linkLabel,
    path: recommendation.path,
    enabled: recommendation.enabled,
    priority: recommendation.priority,
    when: {
      match: recommendation.when.match ?? GroupWhenMatch.all,
      predicates: recommendation.when.predicates.map((predicate) => ({
        metric: predicate.metric,
        op: predicate.op,
        value: predicate.value ?? null,
      })),
    },
  };
}

function toPredicate(when: RecommendationFormValues['when']['predicates'][number]): Predicate {
  if (when.op === PredicateOp.exists || when.value === null) {
    return { metric: when.metric, op: when.op };
  }

  return { metric: when.metric, op: when.op, value: when.value };
}

function toGroupWhen(when: RecommendationFormValues['when']): GroupWhen {
  return {
    match: when.match,
    predicates: when.predicates.map(toPredicate),
  };
}

export function toRecommendationCreate(values: RecommendationFormValues): RecommendationCreate {
  return {
    id: values.id,
    title: values.title,
    text: values.text,
    callout: values.callout,
    linkLabel: values.linkLabel,
    path: values.path,
    enabled: values.enabled,
    priority: values.priority,
    when: toGroupWhen(values.when),
  };
}

export function toRecommendationWrite(values: RecommendationFormValues): RecommendationWrite {
  const created = toRecommendationCreate(values);
  return {
    title: created.title,
    text: created.text,
    callout: created.callout,
    linkLabel: created.linkLabel,
    path: created.path,
    enabled: created.enabled,
    priority: created.priority,
    when: created.when,
  };
}
