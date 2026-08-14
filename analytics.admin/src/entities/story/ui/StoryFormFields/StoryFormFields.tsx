import type { Control } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import FormJsonInput from '@/shared/ui/form/FormJsonInput';
import FormNumberInput from '@/shared/ui/form/FormNumberInput';
import FormSelect from '@/shared/ui/form/FormSelect';
import FormSwitch from '@/shared/ui/form/FormSwitch';
import FormTextInput from '@/shared/ui/form/FormTextInput';
import {
  getPredicateOpOptions,
  getStorySceneTypeOptions,
  getStoryVisibilityOptions,
} from '../../lib/enum-labels';
import type { StoryFormValues } from '../../model/form-schema';

type StoryFormFieldsProps = {
  control: Control<StoryFormValues>;
  idReadOnly?: boolean;
};

export default function StoryFormFields({ control, idReadOnly = false }: StoryFormFieldsProps) {
  const { t } = useTranslation();

  return (
    <>
      <FormTextInput
        name="id"
        control={control}
        label={t('stories.fields.id')}
        disabled={idReadOnly}
      />
      <FormSelect
        name="sceneType"
        control={control}
        label={t('stories.fields.sceneType')}
        data={getStorySceneTypeOptions(t)}
      />
      <FormSelect
        name="visibility"
        control={control}
        label={t('stories.fields.visibility')}
        data={getStoryVisibilityOptions(t)}
      />
      <FormJsonInput
        name="payload"
        control={control}
        label={t('stories.fields.payload')}
        autosize
        minRows={8}
        formatOnBlur
      />
      <FormTextInput name="when.metric" control={control} label={t('stories.fields.whenMetric')} />
      <FormSelect
        name="when.op"
        control={control}
        label={t('stories.fields.whenOp')}
        data={getPredicateOpOptions(t)}
      />
      <FormNumberInput name="when.value" control={control} label={t('stories.fields.whenValue')} />
      <FormSwitch name="enabled" control={control} label={t('stories.fields.enabled')} />
    </>
  );
}
