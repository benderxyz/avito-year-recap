import type { TFunction } from 'i18next';
import type { MetricCreateCurrency } from '@/shared/api/generated/model/metricCreateCurrency';
import { MetricCreateCurrency as MetricCreateCurrencyValues } from '@/shared/api/generated/model/metricCreateCurrency';
import type { MetricCreateSourceField } from '@/shared/api/generated/model/metricCreateSourceField';
import { MetricCreateSourceField as MetricCreateSourceFieldValues } from '@/shared/api/generated/model/metricCreateSourceField';
import type { MetricCreateValueType } from '@/shared/api/generated/model/metricCreateValueType';
import { MetricCreateValueType as MetricCreateValueTypeValues } from '@/shared/api/generated/model/metricCreateValueType';

const VALUE_TYPE_KEYS = {
  number: 'metrics.enums.valueType.number',
  money: 'metrics.enums.valueType.money',
  percentile: 'metrics.enums.valueType.percentile',
  ratio: 'metrics.enums.valueType.ratio',
  string: 'metrics.enums.valueType.string',
  date: 'metrics.enums.valueType.date',
} as const satisfies Record<MetricCreateValueType, string>;

const SOURCE_FIELD_KEYS = {
  value: 'metrics.enums.sourceField.value',
  percentile: 'metrics.enums.sourceField.percentile',
  share: 'metrics.enums.sourceField.share',
} as const satisfies Record<MetricCreateSourceField, string>;

const CURRENCY_KEYS = {
  RUB: 'metrics.enums.currency.RUB',
} as const satisfies Record<Exclude<MetricCreateCurrency, null>, string>;

export function getMetricValueTypeLabel(t: TFunction, value: MetricCreateValueType) {
  return t(VALUE_TYPE_KEYS[value]);
}

export function getMetricSourceFieldLabel(t: TFunction, value: MetricCreateSourceField) {
  return t(SOURCE_FIELD_KEYS[value]);
}

export function getMetricCurrencyLabel(t: TFunction, value: Exclude<MetricCreateCurrency, null>) {
  return t(CURRENCY_KEYS[value]);
}

export function getMetricValueTypeOptions(t: TFunction) {
  return Object.values(MetricCreateValueTypeValues).map((value) => ({
    value,
    label: getMetricValueTypeLabel(t, value),
  }));
}

export function getMetricSourceFieldOptions(t: TFunction) {
  return Object.values(MetricCreateSourceFieldValues).map((value) => ({
    value,
    label: getMetricSourceFieldLabel(t, value),
  }));
}

export function getMetricCurrencyOptions(t: TFunction) {
  return Object.values(MetricCreateCurrencyValues).map((value) => ({
    value,
    label: getMetricCurrencyLabel(t, value),
  }));
}
