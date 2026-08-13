import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import FormTextarea from './FormTextarea';

type FormValues = {
  text: string;
};

function TextField() {
  const { control, watch } = useForm<FormValues>({ defaultValues: { text: '' } });

  return (
    <>
      <FormTextarea name="text" control={control} label="Text" />
      <output>{watch('text')}</output>
    </>
  );
}

describe('FormTextarea', () => {
  it('writes the value into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<TextField />);

    await user.type(getByLabelText('Text'), 'Body');

    expect(getByRole('status')).toHaveTextContent('Body');
  });
});
