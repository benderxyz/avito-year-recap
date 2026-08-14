import { zodResolver } from '@hookform/resolvers/zod';
import { Box, Flex, Group } from '@mantine/core';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  fetchReactPackument,
  getApiAdminPreview,
  getCoreDependency,
  getGetApiAdminPreviewQueryKey,
  loadRecapEngine,
  PREVIEW_PALETTES,
  type PreparedRecap,
  PreviewControls,
  type PreviewFormValues,
  PreviewRecapFrame,
  parseReactPackument,
  previewFormDefaults,
  previewFormSchema,
  type RecapEngineModule,
} from '@/entities/preview';
import CatalogPage from '@/shared/ui/CatalogPage';

type PreviewSession = {
  version: string;
  engine: RecapEngineModule;
  prepared: PreparedRecap;
};

export default function PreviewPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { control, handleSubmit, watch, setValue } = useForm<PreviewFormValues>({
    resolver: zodResolver(previewFormSchema),
    defaultValues: previewFormDefaults,
  });

  const reactVersion = watch('reactVersion');
  const themeId = watch('themeId');
  const autoplay = watch('autoplay');
  const loop = watch('loop');
  const gestures = watch('gestures');
  const tapNav = watch('tapNav');
  const holdToPause = watch('holdToPause');
  const reducedMotion = watch('reducedMotion');

  const packumentQuery = useQuery({
    queryKey: ['npm', '@recap-engine/react'],
    queryFn: fetchReactPackument,
  });

  const parsed = packumentQuery.data ? parseReactPackument(packumentQuery.data) : null;

  useEffect(() => {
    if (parsed?.latest && !reactVersion) {
      setValue('reactVersion', parsed.latest);
    }
  }, [parsed, reactVersion, setValue]);

  const coreVersion = parsed && reactVersion ? getCoreDependency(parsed, reactVersion) : null;

  const [session, setSession] = useState<PreviewSession | null>(null);
  const [isPreviewing, setIsPreviewing] = useState(false);
  const [previewError, setPreviewError] = useState<string | null>(null);

  async function onPreview(values: PreviewFormValues) {
    setIsPreviewing(true);
    setPreviewError(null);

    const params = {
      year: values.year,
      mode: values.mode,
      ...(values.seed === null ? {} : { seed: values.seed }),
    };

    const [payloadResult, engineResult] = await Promise.allSettled([
      queryClient.fetchQuery({
        queryKey: getGetApiAdminPreviewQueryKey(params),
        queryFn: () => getApiAdminPreview(params),
      }),
      loadRecapEngine(values.reactVersion),
    ]);

    if (engineResult.status === 'rejected') {
      setPreviewError(t('preview.errors.engine'));
      setSession(null);
      setIsPreviewing(false);
      return;
    }

    if (payloadResult.status === 'rejected') {
      setPreviewError(t('preview.errors.backend'));
      setSession(null);
      setIsPreviewing(false);
      return;
    }

    try {
      const prepared = engineResult.value.prepareRecap(payloadResult.value);
      setSession({
        version: values.reactVersion,
        engine: engineResult.value,
        prepared,
      });
    } catch {
      setPreviewError(t('preview.errors.prepareRecap'));
      setSession(null);
    } finally {
      setIsPreviewing(false);
    }
  }

  const theme = session ? session.engine.createTheme({ colors: PREVIEW_PALETTES[themeId] }) : null;

  return (
    <CatalogPage
      title={t('preview.title')}
      isError={packumentQuery.isError}
      errorMessage={t('preview.errors.npm')}
    >
      <Group align="flex-start" gap="xl" wrap="wrap">
        <Box w={{ base: '100%', sm: 340 }} flex="none">
          <PreviewControls
            control={control}
            reactVersions={parsed?.versions ?? []}
            coreVersion={coreVersion}
            versionsLoading={packumentQuery.isLoading}
            isSubmitting={isPreviewing}
            onPreview={handleSubmit((values) => void onPreview(values))}
          />
        </Box>
        <Flex
          flex={1}
          justify="center"
          miw={280}
          pos="sticky"
          top="calc(var(--app-shell-header-offset, 60px) + var(--mantine-spacing-md))"
          style={{ alignSelf: 'flex-start' }}
        >
          <PreviewRecapFrame
            engine={session?.engine ?? null}
            engineVersion={session?.version ?? null}
            prepared={session?.prepared ?? null}
            theme={theme}
            autoplay={autoplay}
            loop={loop}
            gestures={gestures}
            tapNav={tapNav}
            holdToPause={holdToPause}
            reducedMotion={reducedMotion}
            isLoading={isPreviewing}
            error={previewError}
          />
        </Flex>
      </Group>
    </CatalogPage>
  );
}
