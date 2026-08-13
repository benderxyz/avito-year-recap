import { zodResolver } from '@hookform/resolvers/zod';
import { useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from '@tanstack/react-router';
import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  getGetApiAdminStoriesIdQueryKey,
  getGetApiAdminStoriesQueryKey,
  StoryFormFields,
  type StoryFormValues,
  storyFormDefaults,
  storyFormSchema,
  storyToFormValues,
  toStoryCreate,
  toStoryWrite,
  useDeleteApiAdminStoriesId,
  useGetApiAdminStoriesId,
  usePostApiAdminStories,
  usePutApiAdminStoriesId,
} from '@/entities/story';
import { routes } from '@/shared/lib/routes';
import ResourceFormPage from '@/shared/ui/ResourceFormPage';

export default function StoryPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const params = useParams({ strict: false });
  const storyId = typeof params.id === 'string' ? params.id : undefined;

  const isCreate = !storyId;

  const { control, handleSubmit, reset } = useForm<StoryFormValues>({
    resolver: zodResolver(storyFormSchema),
    defaultValues: storyFormDefaults,
  });

  const query = useGetApiAdminStoriesId(storyId ?? '', {
    query: { enabled: Boolean(storyId) },
  });

  const createMutation = usePostApiAdminStories();
  const updateMutation = usePutApiAdminStoriesId();
  const deleteMutation = useDeleteApiAdminStoriesId();

  useEffect(() => {
    if (query.data) {
      reset(storyToFormValues(query.data));
    }
  }, [query.data, reset]);

  async function onSubmit(values: StoryFormValues) {
    if (isCreate) {
      const created = await createMutation.mutateAsync({ data: toStoryCreate(values) });
      await queryClient.invalidateQueries({ queryKey: getGetApiAdminStoriesQueryKey() });
      await navigate({ to: routes.storyById, params: { id: created.id } });
      return;
    }

    await updateMutation.mutateAsync({ id: storyId, data: toStoryWrite(values) });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminStoriesQueryKey() });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminStoriesIdQueryKey(storyId) });
  }

  async function onDelete() {
    if (!storyId) {
      return;
    }

    await deleteMutation.mutateAsync({ id: storyId });
    await queryClient.invalidateQueries({ queryKey: getGetApiAdminStoriesQueryKey() });
    await navigate({ to: routes.stories });
  }

  return (
    <ResourceFormPage
      title={isCreate ? t('stories.createTitle') : t('stories.editTitle')}
      isLoading={!isCreate && query.isLoading}
      isError={!isCreate && query.isError}
      errorMessage={t('stories.loadError')}
      onSubmit={handleSubmit((values) => void onSubmit(values))}
      onDelete={isCreate ? undefined : () => void onDelete()}
      isSubmitting={createMutation.isPending || updateMutation.isPending}
      isDeleting={deleteMutation.isPending}
      saveLabel={t('stories.save')}
      deleteLabel={t('stories.delete')}
    >
      <StoryFormFields control={control} idReadOnly={!isCreate} />
    </ResourceFormPage>
  );
}
