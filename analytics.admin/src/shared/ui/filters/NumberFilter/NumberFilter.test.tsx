import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import NumberFilter from './NumberFilter';

function Harness() {
  const [value, setValue] = useState<number | null>(null);
  return (
    <>
      <NumberFilter value={value} onChange={setValue} label="Min priority" />
      <output>{value ?? 'empty'}</output>
    </>
  );
}

describe('NumberFilter', () => {
  it('notifies about the entered number', async () => {
    const user = userEvent.setup();
    const { getByLabelText, getByRole } = renderWithProviders(<Harness />);

    await user.type(getByLabelText('Min priority'), '10');

    expect(getByRole('status')).toHaveTextContent('10');
  });
});
