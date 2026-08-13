import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import FormJsonInput from './FormJsonInput';

type FormValues = {
  payload: string;
};

function PayloadField() {
  const { control, watch } = useForm<FormValues>({ defaultValues: { payload: '' } });

  return (
    <>
      <FormJsonInput name="payload" control={control} label="Payload" />
      <output>{watch('payload')}</output>
    </>
  );
}

describe('FormJsonInput', () => {
  it('writes JSON text into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<PayloadField />);

    await user.type(getByLabelText('Payload'), '{{"a":1}');

    expect(getByRole('status')).toHaveTextContent('{"a":1}');
  });
});
