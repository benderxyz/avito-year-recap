import type { Control } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import FormNumberInput from '@/shared/ui/form/FormNumberInput';
import FormSelect from '@/shared/ui/form/FormSelect';
import FormSwitch from '@/shared/ui/form/FormSwitch';
import FormTextarea from '@/shared/ui/form/FormTextarea';
import FormTextInput from '@/shared/ui/form/FormTextInput';
import { getBadgeVisibilityOptions, getPredicateOpOptions } from '../../lib/enum-labels';
import type { BadgeFormValues } from '../../model/form-schema';

type BadgeFormFieldsProps = {
  control: Control<BadgeFormValues>;
  idReadOnly?: boolean;
};

export default function BadgeFormFields({ control, idReadOnly = false }: BadgeFormFieldsProps) {
  const { t } = useTranslation();

  return (
    <>
      <FormTextInput
        name="id"
        control={control}
        label={t('badges.fields.id')}
        disabled={idReadOnly}
      />
      <FormTextInput name="title" control={control} label={t('badges.fields.title')} />
      <FormTextarea name="description" control={control} label={t('badges.fields.description')} />
      <FormTextInput name="iconUrl" control={control} label={t('badges.fields.iconUrl')} />
      <FormSelect
        name="visibility"
        control={control}
        label={t('badges.fields.visibility')}
        data={getBadgeVisibilityOptions(t)}
      />
      <FormTextInput name="when.metric" control={control} label={t('badges.fields.whenMetric')} />
      <FormSelect
        name="when.op"
        control={control}
        label={t('badges.fields.whenOp')}
        data={getPredicateOpOptions(t)}
      />
      <FormNumberInput name="when.value" control={control} label={t('badges.fields.whenValue')} />
      <FormSwitch name="enabled" control={control} label={t('badges.fields.enabled')} />
    </>
  );
}
