import { z } from 'zod';
import type { BadgeCreate } from '@/shared/api/generated/model/badgeCreate';
import { BadgeCreateVisibility } from '@/shared/api/generated/model/badgeCreateVisibility';
import type { BadgeRule } from '@/shared/api/generated/model/badgeRule';
import type { BadgeWrite } from '@/shared/api/generated/model/badgeWrite';
import type { Predicate } from '@/shared/api/generated/model/predicate';
import { PredicateOp } from '@/shared/api/generated/model/predicateOp';

export const badgeFormSchema = z.object({
  id: z.string().min(1),
  title: z.string().min(1),
  description: z.string().min(1),
  iconUrl: z.string().nullable(),
  sortOrder: z.number(),
  enabled: z.boolean(),
  visibility: z.enum([
    BadgeCreateVisibility.private,
    BadgeCreateVisibility.public,
    BadgeCreateVisibility.both,
  ]),
  when: z.object({
    metric: z.string().min(1),
    op: z.enum([PredicateOp.gt, PredicateOp.gte, PredicateOp.eq, PredicateOp.exists]),
    value: z.number().nullable(),
  }),
});

export type BadgeFormValues = z.infer<typeof badgeFormSchema>;

export const badgeFormDefaults: BadgeFormValues = {
  id: '',
  title: '',
  description: '',
  iconUrl: null,
  sortOrder: 0,
  enabled: true,
  visibility: BadgeCreateVisibility.private,
  when: {
    metric: '',
    op: PredicateOp.gte,
    value: null,
  },
};

export function badgeToFormValues(badge: BadgeRule): BadgeFormValues {
  return {
    id: badge.id,
    title: badge.title,
    description: badge.description,
    iconUrl: badge.iconUrl ?? null,
    sortOrder: badge.sortOrder,
    enabled: badge.enabled,
    visibility: badge.visibility,
    when: {
      metric: badge.when.metric,
      op: badge.when.op,
      value: badge.when.value ?? null,
    },
  };
}

function emptyToNull(value: string | null) {
  return value ? value : null;
}

function toPredicate(when: BadgeFormValues['when']): Predicate {
  if (when.op === PredicateOp.exists || when.value === null) {
    return { metric: when.metric, op: when.op };
  }

  return { metric: when.metric, op: when.op, value: when.value };
}

export function toBadgeCreate(values: BadgeFormValues): BadgeCreate {
  return {
    id: values.id,
    title: values.title,
    description: values.description,
    iconUrl: emptyToNull(values.iconUrl),
    sortOrder: values.sortOrder,
    enabled: values.enabled,
    visibility: values.visibility,
    when: toPredicate(values.when),
  };
}

export function toBadgeWrite(values: BadgeFormValues): BadgeWrite {
  const created = toBadgeCreate(values);
  return {
    title: created.title,
    description: created.description,
    iconUrl: created.iconUrl,
    sortOrder: created.sortOrder,
    enabled: created.enabled,
    visibility: created.visibility,
    when: created.when,
  };
}
