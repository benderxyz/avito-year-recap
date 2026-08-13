import { z } from 'zod';
import type { MetricCreate } from '@/shared/api/generated/model/metricCreate';
import { MetricCreateCurrency } from '@/shared/api/generated/model/metricCreateCurrency';
import { MetricCreateSourceField } from '@/shared/api/generated/model/metricCreateSourceField';
import { MetricCreateValueType } from '@/shared/api/generated/model/metricCreateValueType';
import type { MetricDefinition } from '@/shared/api/generated/model/metricDefinition';
import type { MetricWrite } from '@/shared/api/generated/model/metricWrite';

export const metricFormSchema = z.object({
  key: z.string().min(1),
  valueType: z.enum([
    MetricCreateValueType.number,
    MetricCreateValueType.money,
    MetricCreateValueType.percentile,
    MetricCreateValueType.ratio,
    MetricCreateValueType.string,
    MetricCreateValueType.date,
  ]),
  sourceField: z.enum([
    MetricCreateSourceField.value,
    MetricCreateSourceField.percentile,
    MetricCreateSourceField.share,
  ]),
  sourceKey: z.string(),
  currency: z.enum([MetricCreateCurrency.RUB]).nullable(),
  percentileKey: z.string().nullable(),
  comparisonMinPercentile: z.number().nullable(),
  sortOrder: z.number(),
  enabled: z.boolean(),
  isPublic: z.boolean(),
  includeInLlm: z.boolean(),
});

export type MetricFormValues = z.infer<typeof metricFormSchema>;

export const metricFormDefaults: MetricFormValues = {
  key: '',
  valueType: MetricCreateValueType.number,
  sourceField: MetricCreateSourceField.value,
  sourceKey: '',
  currency: null,
  percentileKey: null,
  comparisonMinPercentile: null,
  sortOrder: 0,
  enabled: true,
  isPublic: false,
  includeInLlm: false,
};

export function metricToFormValues(metric: MetricDefinition): MetricFormValues {
  return {
    key: metric.key,
    valueType: metric.valueType,
    sourceField: metric.sourceField,
    sourceKey: metric.sourceKey,
    currency: metric.currency ?? null,
    percentileKey: metric.percentileKey ?? null,
    comparisonMinPercentile: metric.comparisonMinPercentile ?? null,
    sortOrder: metric.sortOrder,
    enabled: metric.enabled,
    isPublic: metric.isPublic,
    includeInLlm: metric.includeInLlm,
  };
}

function emptyToNull(value: string | null) {
  return value ? value : null;
}

export function toMetricCreate(values: MetricFormValues): MetricCreate {
  return {
    key: values.key,
    valueType: values.valueType,
    sourceField: values.sourceField,
    sourceKey: values.sourceKey,
    currency: values.currency,
    percentileKey: emptyToNull(values.percentileKey),
    comparisonMinPercentile: values.comparisonMinPercentile,
    sortOrder: values.sortOrder,
    enabled: values.enabled,
    isPublic: values.isPublic,
    includeInLlm: values.includeInLlm,
  };
}

export function toMetricWrite(values: MetricFormValues): MetricWrite {
  const created = toMetricCreate(values);
  return {
    valueType: created.valueType,
    sourceField: created.sourceField,
    sourceKey: created.sourceKey,
    currency: created.currency,
    percentileKey: created.percentileKey,
    comparisonMinPercentile: created.comparisonMinPercentile,
    sortOrder: created.sortOrder,
    enabled: created.enabled,
    isPublic: created.isPublic,
    includeInLlm: created.includeInLlm,
  };
}
