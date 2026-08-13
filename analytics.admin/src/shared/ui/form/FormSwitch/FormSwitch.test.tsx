import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import FormSwitch from './FormSwitch';

type FormValues = {
  enabled: boolean;
};

function EnabledField() {
  const { control, watch } = useForm<FormValues>({ defaultValues: { enabled: false } });

  return (
    <>
      <FormSwitch name="enabled" control={control} label="Enabled" />
      <output>{String(watch('enabled'))}</output>
    </>
  );
}

describe('FormSwitch', () => {
  it('writes a boolean into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByRole } = renderWithProviders(<EnabledField />);

    await user.click(getByRole('switch', { name: 'Enabled' }));

    expect(getByRole('status')).toHaveTextContent('true');
  });
});
