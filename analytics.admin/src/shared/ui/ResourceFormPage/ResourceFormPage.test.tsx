import userEvent from '@testing-library/user-event';
import type { FormEvent } from 'react';
import { describe, expect, it, vi } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import ResourceFormPage from './ResourceFormPage';

describe('ResourceFormPage', () => {
  it('submits the form', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn((event: FormEvent<HTMLFormElement>) => event.preventDefault());
    const { getByRole } = renderWithProviders(
      <ResourceFormPage title="Badge" onSubmit={onSubmit}>
        <input name="title" defaultValue="x" />
      </ResourceFormPage>,
    );

    await user.click(getByRole('button', { name: 'Save' }));

    expect(onSubmit).toHaveBeenCalledOnce();
  });

  it('calls onDelete from the delete button', async () => {
    const user = userEvent.setup();
    const onDelete = vi.fn();
    const { getByRole } = renderWithProviders(
      <ResourceFormPage
        title="Badge"
        onSubmit={(event) => event.preventDefault()}
        onDelete={onDelete}
      >
        <div>Fields</div>
      </ResourceFormPage>,
    );

    await user.click(getByRole('button', { name: 'Delete' }));

    expect(onDelete).toHaveBeenCalledOnce();
  });

  it('hides the form while loading', () => {
    const { queryByRole, queryByText } = renderWithProviders(
      <ResourceFormPage title="Badge" isLoading onSubmit={(event) => event.preventDefault()}>
        <div>Fields</div>
      </ResourceFormPage>,
    );

    expect(queryByText('Fields')).not.toBeInTheDocument();
    expect(queryByRole('button', { name: 'Save' })).not.toBeInTheDocument();
  });
});
