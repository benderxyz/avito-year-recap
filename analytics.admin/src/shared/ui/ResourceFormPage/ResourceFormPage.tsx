import { Alert, Button, Group, Loader, Stack, Title } from '@mantine/core';
import type { ReactNode, SubmitEventHandler } from 'react';

type ResourceFormPageProps = {
  title: ReactNode;
  isLoading?: boolean;
  isError?: boolean;
  errorMessage?: string;
  onSubmit: SubmitEventHandler<HTMLFormElement>;
  onDelete?: () => void;
  isSubmitting?: boolean;
  isDeleting?: boolean;
  saveLabel?: string;
  deleteLabel?: string;
  children: ReactNode;
};

export default function ResourceFormPage({
  title,
  isLoading = false,
  isError = false,
  errorMessage = 'Failed to load',
  onSubmit,
  onDelete,
  isSubmitting = false,
  isDeleting = false,
  saveLabel = 'Save',
  deleteLabel = 'Delete',
  children,
}: ResourceFormPageProps) {
  return (
    <Stack gap="md">
      <Title order={2}>{title}</Title>
      {isError ? <Alert color="red">{errorMessage}</Alert> : null}
      {isLoading ? (
        <Loader aria-label="Loading" />
      ) : (
        <form onSubmit={onSubmit}>
          <Stack gap="md">
            {children}
            <Group>
              <Button type="submit" loading={isSubmitting}>
                {saveLabel}
              </Button>
              {onDelete ? (
                <Button
                  type="button"
                  color="red"
                  variant="light"
                  loading={isDeleting}
                  onClick={onDelete}
                >
                  {deleteLabel}
                </Button>
              ) : null}
            </Group>
          </Stack>
        </form>
      )}
    </Stack>
  );
}
