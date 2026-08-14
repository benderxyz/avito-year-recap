import { Button, Stack, TextInput } from '@mantine/core';
import type { Control } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import FormNumberInput from '@/shared/ui/form/FormNumberInput';
import FormSelect from '@/shared/ui/form/FormSelect';
import FormSwitch from '@/shared/ui/form/FormSwitch';
import { PREVIEW_PALETTE_IDS } from '../../lib/palettes';
import type { PreviewFormValues } from '../../model/form-schema';

type PreviewControlsProps = {
  control: Control<PreviewFormValues>;
  reactVersions: string[];
  coreVersion: string | null;
  versionsLoading?: boolean;
  isSubmitting: boolean;
  onPreview: () => void;
};

export default function PreviewControls({
  control,
  reactVersions,
  coreVersion,
  versionsLoading = false,
  isSubmitting,
  onPreview,
}: PreviewControlsProps) {
  const { t } = useTranslation();

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        onPreview();
      }}
    >
      <Stack gap="md">
        <FormSelect
          name="reactVersion"
          control={control}
          label={t('preview.fields.reactVersion')}
          data={reactVersions.map((version) => ({ value: version, label: version }))}
          searchable
          disabled={versionsLoading}
          w="100%"
        />
        <TextInput
          label={t('preview.fields.coreVersion')}
          value={coreVersion ?? '—'}
          disabled
          w="100%"
        />
        <FormNumberInput name="year" control={control} label={t('preview.fields.year')} w="100%" />
        <FormSelect
          name="mode"
          control={control}
          label={t('preview.fields.mode')}
          data={[
            { value: 'private', label: t('preview.modes.private') },
            { value: 'public', label: t('preview.modes.public') },
          ]}
          w="100%"
        />
        <FormNumberInput name="seed" control={control} label={t('preview.fields.seed')} w="100%" />
        <FormSelect
          name="themeId"
          control={control}
          label={t('preview.fields.theme')}
          data={PREVIEW_PALETTE_IDS.map((id) => ({
            value: id,
            label: t(`preview.themes.${id}`),
          }))}
          w="100%"
        />
        <FormSwitch name="autoplay" control={control} label={t('preview.fields.autoplay')} />
        <FormSwitch name="loop" control={control} label={t('preview.fields.loop')} />
        <FormSwitch name="gestures" control={control} label={t('preview.fields.gestures')} />
        <FormSwitch name="tapNav" control={control} label={t('preview.fields.tapNav')} />
        <FormSwitch name="holdToPause" control={control} label={t('preview.fields.holdToPause')} />
        <FormSwitch
          name="reducedMotion"
          control={control}
          label={t('preview.fields.reducedMotion')}
        />
        <Button type="submit" loading={isSubmitting}>
          {t('preview.run')}
        </Button>
      </Stack>
    </form>
  );
}
