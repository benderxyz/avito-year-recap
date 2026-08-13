import userEvent from '@testing-library/user-event';
import { useForm } from 'react-hook-form';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import FormTextInput from './FormTextInput';

type FormValues = {
  title: string;
};

function TitleField({ onReady }: { onReady?: (setError: () => void) => void }) {
  const { control, watch, setError } = useForm<FormValues>({ defaultValues: { title: '' } });
  onReady?.(() => setError('title', { message: 'Required' }));

  return (
    <>
      <FormTextInput name="title" control={control} label="Title" />
      <output>{watch('title')}</output>
    </>
  );
}

describe('FormTextInput', () => {
  it('writes the value into react-hook-form', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<TitleField />);

    await user.type(getByLabelText('Title'), 'Hello');

    expect(getByRole('status')).toHaveTextContent('Hello');
  });

  it('shows a field error from the form state', async () => {
    const user = userEvent.setup();
    let triggerError = () => {};
    const { getByText, getByRole } = renderWithProviders(
      <>
        <TitleField
          onReady={(setError) => {
            triggerError = setError;
          }}
        />
        <button type="button" onClick={() => triggerError()}>
          set-error
        </button>
      </>,
    );

    await user.click(getByRole('button', { name: 'set-error' }));

    expect(getByText('Required')).toBeInTheDocument();
  });
});
