import { parseAsBoolean, parseAsString, parseAsStringEnum } from 'nuqs';
import { GetApiAdminBadgesVisibility } from '@/shared/api/generated/model/getApiAdminBadgesVisibility';

export const badgeFilterParsers = {
  search: parseAsString.withDefault(''),
  enabled: parseAsBoolean,
  visibility: parseAsStringEnum([
    GetApiAdminBadgesVisibility.private,
    GetApiAdminBadgesVisibility.public,
    GetApiAdminBadgesVisibility.both,
  ]),
  metric: parseAsString.withDefault(''),
};
