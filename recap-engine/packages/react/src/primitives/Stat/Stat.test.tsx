import { jest } from '@jest/globals';
import { screen, waitFor } from '@testing-library/react';
import { renderWithRecap } from '../../test/renderWithRecap';
import { Stat } from './Stat';

describe('Stat', () => {
  it('formats values, resolves plurals, applies classes and completes reduced motion', async () => {
    const onMotionComplete = jest.fn();
    const { value } = renderWithRecap(
      <Stat
        value={21}
        unit={{ one: 'товар', few: 'товара', many: 'товаров' }}
        valueFormat={{ minimumIntegerDigits: 3 }}
        classNames={{ root: 'root-extra', value: 'value-extra', unit: 'unit-extra' }}
        onMotionComplete={onMotionComplete}
      />,
    );

    const displayedValue = screen.getByText('021');
    expect(displayedValue.parentElement).toHaveClass('value-extra');
    expect(displayedValue.parentElement).toHaveAttribute('data-target-value', '021');
    expect(screen.getByText('товар')).toHaveClass('unit-extra');
    expect(displayedValue.parentElement?.parentElement).toHaveClass('root-extra');
    await waitFor(() => expect(onMotionComplete).toHaveBeenCalledTimes(1));
    expect(value.notifyBlockMotionComplete).toHaveBeenCalledTimes(1);
  });
});
