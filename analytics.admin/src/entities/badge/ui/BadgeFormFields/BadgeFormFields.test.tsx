import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import { type BadgeFormValues, badgeFormDefaults } from '../../model/form-schema';
import BadgeFormFields from './BadgeFormFields';

function Harness() {
  const { control, watch } = useForm<BadgeFormValues>({ defaultValues: badgeFormDefaults });

  return (
    <>
      <BadgeFormFields control={control} />
      <output>{watch('title')}</output>
    </>
  );
}

describe('BadgeFormFields', () => {
  it('writes the title into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<Harness />);

    await user.type(getByLabelText('Title'), 'Top seller');

    expect(getByRole('status')).toHaveTextContent('Top seller');
  });
});
