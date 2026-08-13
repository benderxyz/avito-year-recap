import { zodResolver } from '@hookform/resolvers/zod';
import { useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from '@tanstack/react-router';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import MetricFormFields from '@/entities/metric/MetricFormFields';
import {
  type MetricFormValues,
  metricFormDefaults,
  metricFormSchema,
  metricToFormValues,
  toMetricCreate,
  toMetricWrite,
} from '@/entities/metric/metric-form-schema';
import {
  getGetApiAdminMetricsKeyQueryKey,
  getGetApiAdminMetricsQueryKey,
  useDeleteApiAdminMetricsKey,
  useGetApiAdminMetricsKey,
  usePostApiAdminMetrics,
  usePutApiAdminMetricsKey,
} from '@/shared/api/generated/metrics/metrics';
import { routes } from '@/shared/lib/routes';
import ResourceFormPage from '@/shared/ui/ResourceFormPage';

export default function MetricPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const params = useParams({ strict: false });
  const metricKey = typeof params.key === 'string' ? params.key : undefined;

  const isCreate = !metricKey;

  const { control, handleSubmit, reset } = useForm<MetricFormValues>({
    resolver: zodResolver(metricFormSchema),
    defaultValues: metricFormDefaults,
  });

  const query = useGetApiAdminMetricsKey(metricKey ?? '', {
    query: { enabled: Boolean(metricKey) },
  });

  const createMutation = usePostApiAdminMetrics();
  const updateMutation = usePutApiAdminMetricsKey();
  const deleteMutation = useDeleteApiAdminMetricsKey();

  useEffect(() => {
    if (query.data) {
      reset(metricToFormValues(query.data));
    }
  }, [query.data, reset]);

  async function onSubmit(values: MetricFormValues) {
    if (isCreate) {
      const created = await createMutation.mutateAsync({ data: toMetricCreate(values) });
      await queryClient.invalidateQueries({ queryKey: getGetApiAdminMetricsQueryKey() });
      await navigate({ to: routes.metricByKey, params: { key: created.key } });
      return;
    }

    await updateMutation.mutateAsync({ key: metricKey, data: toMetricWrite(values) });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminMetricsQueryKey() });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminMetricsKeyQueryKey(metricKey) });
  }

  async function onDelete() {
    if (!metricKey) {
      return;
    }

    await deleteMutation.mutateAsync({ key: metricKey });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminMetricsQueryKey() });
    await navigate({ to: routes.metrics });
  }

  return (
    <ResourceFormPage
      title={isCreate ? t('metrics.createTitle') : t('metrics.editTitle')}
      isLoading={!isCreate && query.isLoading}
      isError={!isCreate && query.isError}
      errorMessage={t('metrics.loadError')}
      onSubmit={handleSubmit((values) => void onSubmit(values))}
      onDelete={isCreate ? undefined : () => void onDelete()}
      isSubmitting={createMutation.isPending || updateMutation.isPending}
      isDeleting={deleteMutation.isPending}
      saveLabel={t('metrics.save')}
      deleteLabel={t('metrics.delete')}
    >
      <MetricFormFields control={control} keyReadOnly={!isCreate} />
    </ResourceFormPage>
  );
}
