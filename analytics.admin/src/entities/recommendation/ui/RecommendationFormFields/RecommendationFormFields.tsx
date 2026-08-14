import { Button, Group, Stack } from '@mantine/core';
import type { Control } from 'react-hook-form';
import { useFieldArray } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import FormNumberInput from '@/shared/ui/form/FormNumberInput';
import FormSelect from '@/shared/ui/form/FormSelect';
import FormSwitch from '@/shared/ui/form/FormSwitch';
import FormTextarea from '@/shared/ui/form/FormTextarea';
import FormTextInput from '@/shared/ui/form/FormTextInput';
import { getMatchModeOptions, getPredicateOpOptions } from '../../lib/enum-labels';
import { emptyPredicate, type RecommendationFormValues } from '../../model/form-schema';

type RecommendationFormFieldsProps = {
  control: Control<RecommendationFormValues>;
  idReadOnly?: boolean;
};

export default function RecommendationFormFields({
  control,
  idReadOnly = false,
}: RecommendationFormFieldsProps) {
  const { t } = useTranslation();
  const { fields, append, remove } = useFieldArray({ control, name: 'when.predicates' });

  return (
    <>
      <FormTextInput
        name="id"
        control={control}
        label={t('recommendations.fields.id')}
        disabled={idReadOnly}
      />
      <FormTextInput name="title" control={control} label={t('recommendations.fields.title')} />
      <FormTextarea name="text" control={control} label={t('recommendations.fields.text')} />
      <FormTextInput name="callout" control={control} label={t('recommendations.fields.callout')} />
      <FormTextInput
        name="linkLabel"
        control={control}
        label={t('recommendations.fields.linkLabel')}
      />
      <FormTextInput name="path" control={control} label={t('recommendations.fields.path')} />
      <FormSelect
        name="when.match"
        control={control}
        label={t('recommendations.fields.whenMatch')}
        data={getMatchModeOptions(t)}
      />
      <Stack gap="sm">
        {fields.map((field, index) => (
          <Group key={field.id} align="flex-end">
            <FormTextInput
              name={`when.predicates.${index}.metric`}
              control={control}
              label={t('recommendations.fields.whenMetric')}
            />
            <FormSelect
              name={`when.predicates.${index}.op`}
              control={control}
              label={t('recommendations.fields.whenOp')}
              data={getPredicateOpOptions(t)}
            />
            <FormNumberInput
              name={`when.predicates.${index}.value`}
              control={control}
              label={t('recommendations.fields.whenValue')}
            />
            <Button type="button" variant="subtle" onClick={() => remove(index)}>
              {t('recommendations.removePredicate')}
            </Button>
          </Group>
        ))}
        <Button type="button" variant="light" onClick={() => append(emptyPredicate)}>
          {t('recommendations.addPredicate')}
        </Button>
      </Stack>
      <FormSwitch name="enabled" control={control} label={t('recommendations.fields.enabled')} />
    </>
  );
}
