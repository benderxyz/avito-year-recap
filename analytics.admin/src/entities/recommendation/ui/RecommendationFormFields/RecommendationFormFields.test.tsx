import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import { type RecommendationFormValues, recommendationFormDefaults } from '../../model/form-schema';
import RecommendationFormFields from './RecommendationFormFields';

function Harness() {
  const { control, watch } = useForm<RecommendationFormValues>({
    defaultValues: recommendationFormDefaults,
  });

  return (
    <>
      <RecommendationFormFields control={control} />
      <output>{watch('title')}</output>
    </>
  );
}

describe('RecommendationFormFields', () => {
  it('writes the title into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<Harness />);

    await user.type(getByLabelText('Title'), 'More orders');

    expect(getByRole('status')).toHaveTextContent('More orders');
  });
});
