import { zodResolver } from '@hookform/resolvers/zod';
import { useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from '@tanstack/react-router';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  BadgeFormFields,
  type BadgeFormValues,
  badgeFormDefaults,
  badgeFormSchema,
  badgeToFormValues,
  getGetApiAdminBadgesIdQueryKey,
  getGetApiAdminBadgesQueryKey,
  toBadgeCreate,
  toBadgeWrite,
  useDeleteApiAdminBadgesId,
  useGetApiAdminBadgesId,
  usePostApiAdminBadges,
  usePutApiAdminBadgesId,
} from '@/entities/badge';
import { routes } from '@/shared/lib/routes';
import ResourceFormPage from '@/shared/ui/ResourceFormPage';

export default function BadgePage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const params = useParams({ strict: false });
  const badgeId = typeof params.id === 'string' ? params.id : undefined;

  const isCreate = !badgeId;

  const { control, handleSubmit, reset } = useForm<BadgeFormValues>({
    resolver: zodResolver(badgeFormSchema),
    defaultValues: badgeFormDefaults,
  });

  const query = useGetApiAdminBadgesId(badgeId ?? '', {
    query: { enabled: Boolean(badgeId) },
  });

  const createMutation = usePostApiAdminBadges();
  const updateMutation = usePutApiAdminBadgesId();
  const deleteMutation = useDeleteApiAdminBadgesId();

  useEffect(() => {
    if (query.data) {
      reset(badgeToFormValues(query.data));
    }
  }, [query.data, reset]);

  async function onSubmit(values: BadgeFormValues) {
    if (isCreate) {
      const created = await createMutation.mutateAsync({ data: toBadgeCreate(values) });
      await queryClient.invalidateQueries({ queryKey: getGetApiAdminBadgesQueryKey() });
      await navigate({ to: routes.badgeById, params: { id: created.id } });
      return;
    }

    await updateMutation.mutateAsync({ id: badgeId, data: toBadgeWrite(values) });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminBadgesQueryKey() });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminBadgesIdQueryKey(badgeId) });
  }

  async function onDelete() {
    if (!badgeId) {
      return;
    }

    await deleteMutation.mutateAsync({ id: badgeId });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminBadgesQueryKey() });
    await navigate({ to: routes.badges });
  }

  return (
    <ResourceFormPage
      title={isCreate ? t('badges.createTitle') : t('badges.editTitle')}
      isLoading={!isCreate && query.isLoading}
      isError={!isCreate && query.isError}
      errorMessage={t('badges.loadError')}
      onSubmit={handleSubmit((values) => void onSubmit(values))}
      onDelete={isCreate ? undefined : () => void onDelete()}
      isSubmitting={createMutation.isPending || updateMutation.isPending}
      isDeleting={deleteMutation.isPending}
      saveLabel={t('badges.save')}
      deleteLabel={t('badges.delete')}
    >
      <BadgeFormFields control={control} idReadOnly={!isCreate} />
    </ResourceFormPage>
  );
}
