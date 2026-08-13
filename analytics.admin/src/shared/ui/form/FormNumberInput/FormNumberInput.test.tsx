import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import FormNumberInput from './FormNumberInput';

type FormValues = {
  sortOrder: number | string;
};

function SortField() {
  const { control, watch } = useForm<FormValues>({ defaultValues: { sortOrder: '' } });

  return (
    <>
      <FormNumberInput name="sortOrder" control={control} label="Sort order" />
      <output>{String(watch('sortOrder'))}</output>
    </>
  );
}

describe('FormNumberInput', () => {
  it('writes a number into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<SortField />);

    await user.type(getByLabelText('Sort order'), '12');

    expect(getByRole('status')).toHaveTextContent('12');
  });
});
