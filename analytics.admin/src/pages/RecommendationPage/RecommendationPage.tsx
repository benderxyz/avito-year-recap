import { zodResolver } from '@hookform/resolvers/zod';
import { useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from '@tanstack/react-router';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  getGetApiAdminRecommendationsIdQueryKey,
  getGetApiAdminRecommendationsQueryKey,
  RecommendationFormFields,
  type RecommendationFormValues,
  recommendationFormDefaults,
  recommendationFormSchema,
  recommendationToFormValues,
  toRecommendationCreate,
  toRecommendationWrite,
  useDeleteApiAdminRecommendationsId,
  useGetApiAdminRecommendationsId,
  usePostApiAdminRecommendations,
  usePutApiAdminRecommendationsId,
} from '@/entities/recommendation';
import { routes } from '@/shared/lib/routes';
import ResourceFormPage from '@/shared/ui/ResourceFormPage';

export default function RecommendationPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const params = useParams({ strict: false });
  const recommendationId = typeof params.id === 'string' ? params.id : undefined;

  const isCreate = !recommendationId;

  const { control, handleSubmit, reset } = useForm<RecommendationFormValues>({
    resolver: zodResolver(recommendationFormSchema),
    defaultValues: recommendationFormDefaults,
  });

  const query = useGetApiAdminRecommendationsId(recommendationId ?? '', {
    query: { enabled: Boolean(recommendationId) },
  });

  const createMutation = usePostApiAdminRecommendations();
  const updateMutation = usePutApiAdminRecommendationsId();
  const deleteMutation = useDeleteApiAdminRecommendationsId();

  useEffect(() => {
    if (query.data) {
      reset(recommendationToFormValues(query.data));
    }
  }, [query.data, reset]);

  async function onSubmit(values: RecommendationFormValues) {
    if (isCreate) {
      const created = await createMutation.mutateAsync({ data: toRecommendationCreate(values) });
      await queryClient.invalidateQueries({ queryKey: getGetApiAdminRecommendationsQueryKey() });
      await navigate({ to: routes.recommendationById, params: { id: created.id } });
      return;
    }

    await updateMutation.mutateAsync({
      id: recommendationId,
      data: toRecommendationWrite(values),
    });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminRecommendationsQueryKey() });
    await queryClient.invalidateQueries({
      queryKey: getGetApiAdminRecommendationsIdQueryKey(recommendationId),
    });
  }

  async function onDelete() {
    if (!recommendationId) {
      return;
    }

    await deleteMutation.mutateAsync({ id: recommendationId });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminRecommendationsQueryKey() });
    await navigate({ to: routes.recommendations });
  }

  return (
    <ResourceFormPage
      title={isCreate ? t('recommendations.createTitle') : t('recommendations.editTitle')}
      isLoading={!isCreate && query.isLoading}
      isError={!isCreate && query.isError}
      errorMessage={t('recommendations.loadError')}
      onSubmit={handleSubmit((values) => void onSubmit(values))}
      onDelete={isCreate ? undefined : () => void onDelete()}
      isSubmitting={createMutation.isPending || updateMutation.isPending}
      isDeleting={deleteMutation.isPending}
      saveLabel={t('recommendations.save')}
      deleteLabel={t('recommendations.delete')}
    >
      <RecommendationFormFields control={control} idReadOnly={!isCreate} />
    </ResourceFormPage>
  );
}
