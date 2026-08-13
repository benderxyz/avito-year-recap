import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import { describe, expect, it } from 'vitest';
import { renderWithProviders } from '@/shared/lib/test/render-with-providers';
import SearchFilter from './SearchFilter';

function Harness() {
  const [value, setValue] = useState('');
  return <SearchFilter value={value} onChange={setValue} />;
}

describe('SearchFilter', () => {
  it('notifies about text changes', async () => {
    const user = userEvent.setup();
    const { getByLabelText } = renderWithProviders(<Harness />);

    await user.type(getByLabelText('Search'), 'abc');

    expect(getByLabelText('Search')).toHaveValue('abc');
  });
});
