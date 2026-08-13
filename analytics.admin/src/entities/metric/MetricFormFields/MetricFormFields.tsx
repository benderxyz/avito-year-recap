import type { Control } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import FormNumberInput from '@/shared/ui/form/FormNumberInput';
import FormSelect from '@/shared/ui/form/FormSelect';
import FormSwitch from '@/shared/ui/form/FormSwitch';
import FormTextInput from '@/shared/ui/form/FormTextInput';
import {
  getMetricCurrencyOptions,
  getMetricSourceFieldOptions,
  getMetricValueTypeOptions,
} from '../metric-enum-labels';
import type { MetricFormValues } from '../metric-form-schema';

type MetricFormFieldsProps = {
  control: Control<MetricFormValues>;
  keyReadOnly?: boolean;
};

export default function MetricFormFields({ control, keyReadOnly = false }: MetricFormFieldsProps) {
  const { t } = useTranslation();

  return (
    <>
      <FormTextInput
        name="key"
        control={control}
        label={t('metrics.fields.key')}
        disabled={keyReadOnly}
      />
      <FormSelect
        name="valueType"
        control={control}
        label={t('metrics.fields.valueType')}
        data={getMetricValueTypeOptions(t)}
      />
      <FormSelect
        name="sourceField"
        control={control}
        label={t('metrics.fields.sourceField')}
        data={getMetricSourceFieldOptions(t)}
      />
      <FormTextInput name="sourceKey" control={control} label={t('metrics.fields.sourceKey')} />
      <FormSelect
        name="currency"
        control={control}
        label={t('metrics.fields.currency')}
        data={getMetricCurrencyOptions(t)}
        clearable
      />
      <FormTextInput
        name="percentileKey"
        control={control}
        label={t('metrics.fields.percentileKey')}
      />
      <FormNumberInput
        name="comparisonMinPercentile"
        control={control}
        label={t('metrics.fields.comparisonMinPercentile')}
      />
      <FormNumberInput name="sortOrder" control={control} label={t('metrics.fields.sortOrder')} />
      <FormSwitch name="enabled" control={control} label={t('metrics.fields.enabled')} />
      <FormSwitch name="isPublic" control={control} label={t('metrics.fields.isPublic')} />
      <FormSwitch name="includeInLlm" control={control} label={t('metrics.fields.includeInLlm')} />
    </>
  );
}
