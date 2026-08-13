import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import FormSelect from './FormSelect';

type FormValues = {
  visibility: string | null;
};

function VisibilityField() {
  const { control, watch } = useForm<FormValues>({ defaultValues: { visibility: null } });

  return (
    <>
      <FormSelect
        name="visibility"
        control={control}
        label="Visibility"
        data={[
          { value: 'private', label: 'Private' },
          { value: 'public', label: 'Public' },
        ]}
      />
      <output>{watch('visibility') ?? ''}</output>
    </>
  );
}

describe('FormSelect', () => {
  it('writes the selected option into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByRole, getByText } = renderWithProviders(<VisibilityField />);

    await user.click(getByRole('textbox', { name: 'Visibility' }));
    await user.click(getByRole('option', { name: 'Public', hidden: true }));

    expect(getByText('public')).toBeInTheDocument();
  });
});
