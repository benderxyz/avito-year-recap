import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import { type MetricFormValues, metricFormDefaults } from '../../model/form-schema';
import MetricFormFields from './MetricFormFields';

function Harness() {
  const { control, watch } = useForm<MetricFormValues>({ defaultValues: metricFormDefaults });

  return (
    <>
      <MetricFormFields control={control} />
      <output>{watch('key')}</output>
    </>
  );
}

describe('MetricFormFields', () => {
  it('writes the key into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<Harness />);

    await user.type(getByLabelText('Key'), 'orders_count');

    expect(getByRole('status')).toHaveTextContent('orders_count');
  });
});
