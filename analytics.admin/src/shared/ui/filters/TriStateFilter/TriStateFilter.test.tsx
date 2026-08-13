import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import TriStateFilter from './TriStateFilter';

function Harness() {
  const [value, setValue] = useState<boolean | undefined>(undefined);
  return (
    <>
      <TriStateFilter value={value} onChange={setValue} />
      <output>{value === undefined ? 'all' : String(value)}</output>
    </>
  );
}

describe('TriStateFilter', () => {
  it('maps Yes to true', async () => {
    const user = userEvent.setup();
    const { getByRole, getByText } = renderWithProviders(<Harness />);

    await user.click(getByRole('textbox', { name: 'Enabled' }));
    await user.click(getByRole('option', { name: 'Yes', hidden: true }));

    expect(getByText('true')).toBeInTheDocument();
  });
});
