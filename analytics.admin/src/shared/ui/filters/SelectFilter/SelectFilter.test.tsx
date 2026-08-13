import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import SelectFilter from './SelectFilter';

function Harness() {
  const [value, setValue] = useState<string | null>(null);
  return (
    <>
      <SelectFilter
        value={value}
        onChange={setValue}
        label="Visibility"
        data={[
          { value: 'private', label: 'Private' },
          { value: 'public', label: 'Public' },
        ]}
      />
      <output>{value ?? 'all'}</output>
    </>
  );
}

describe('SelectFilter', () => {
  it('notifies about the selected option', async () => {
    const user = userEvent.setup();
    const { getByRole, getByText } = renderWithProviders(<Harness />);

    await user.click(getByRole('textbox', { name: 'Visibility' }));
    await user.click(getByRole('option', { name: 'Public', hidden: true }));

    expect(getByText('public')).toBeInTheDocument();
  });
});
