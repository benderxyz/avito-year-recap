import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import { type StoryFormValues, storyFormDefaults } from '../../model/form-schema';
import StoryFormFields from './StoryFormFields';

function Harness() {
  const { control, watch } = useForm<StoryFormValues>({ defaultValues: storyFormDefaults });

  return (
    <>
      <StoryFormFields control={control} />
      <output>{watch('id')}</output>
    </>
  );
}

describe('StoryFormFields', () => {
  it('writes the id into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<Harness />);

    await user.type(getByLabelText('Id'), 'year_intro');

    expect(getByRole('status')).toHaveTextContent('year_intro');
  });
});
