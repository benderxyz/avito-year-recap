import type { TFunction } from 'i18next';
import type { BadgeCreateVisibility } from '@/shared/api/generated/model/badgeCreateVisibility';
import { BadgeCreateVisibility as BadgeCreateVisibilityValues } from '@/shared/api/generated/model/badgeCreateVisibility';
import type { PredicateOp } from '@/shared/api/generated/model/predicateOp';
import { PredicateOp as PredicateOpValues } from '@/shared/api/generated/model/predicateOp';

const VISIBILITY_KEYS = {
  private: 'badges.enums.visibility.private',
  public: 'badges.enums.visibility.public',
  both: 'badges.enums.visibility.both',
} as const satisfies Record<BadgeCreateVisibility, string>;

const PREDICATE_OP_KEYS = {
  gt: 'badges.enums.predicateOp.gt',
  gte: 'badges.enums.predicateOp.gte',
  eq: 'badges.enums.predicateOp.eq',
  exists: 'badges.enums.predicateOp.exists',
} as const satisfies Record<PredicateOp, string>;

export function getBadgeVisibilityLabel(t: TFunction, value: BadgeCreateVisibility) {
  return t(VISIBILITY_KEYS[value]);
}

export function getPredicateOpLabel(t: TFunction, value: PredicateOp) {
  return t(PREDICATE_OP_KEYS[value]);
}

export function getBadgeVisibilityOptions(t: TFunction) {
  return Object.values(BadgeCreateVisibilityValues).map((value) => ({
    value,
    label: getBadgeVisibilityLabel(t, value),
  }));
}

export function getPredicateOpOptions(t: TFunction) {
  return Object.values(PredicateOpValues).map((value) => ({
    value,
    label: getPredicateOpLabel(t, value),
  }));
}
