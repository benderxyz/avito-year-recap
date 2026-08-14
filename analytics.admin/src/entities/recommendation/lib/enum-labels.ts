import type { TFunction } from 'i18next';
import type { GroupWhenMatch } from '@/shared/api/generated/model/groupWhenMatch';
import { GroupWhenMatch as GroupWhenMatchValues } from '@/shared/api/generated/model/groupWhenMatch';
import type { PredicateOp } from '@/shared/api/generated/model/predicateOp';
import { PredicateOp as PredicateOpValues } from '@/shared/api/generated/model/predicateOp';

const MATCH_MODE_KEYS = {
  all: 'recommendations.enums.matchMode.all',
  any: 'recommendations.enums.matchMode.any',
} as const satisfies Record<GroupWhenMatch, string>;

const PREDICATE_OP_KEYS = {
  gt: 'recommendations.enums.predicateOp.gt',
  gte: 'recommendations.enums.predicateOp.gte',
  eq: 'recommendations.enums.predicateOp.eq',
  exists: 'recommendations.enums.predicateOp.exists',
} as const satisfies Record<PredicateOp, string>;

export function getMatchModeLabel(t: TFunction, value: GroupWhenMatch) {
  return t(MATCH_MODE_KEYS[value]);
}

export function getPredicateOpLabel(t: TFunction, value: PredicateOp) {
  return t(PREDICATE_OP_KEYS[value]);
}

export function getMatchModeOptions(t: TFunction) {
  return Object.values(GroupWhenMatchValues).map((value) => ({
    value,
    label: getMatchModeLabel(t, value),
  }));
}

export function getPredicateOpOptions(t: TFunction) {
  return Object.values(PredicateOpValues).map((value) => ({
    value,
    label: getPredicateOpLabel(t, value),
  }));
}
