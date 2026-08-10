import { MantineProvider } from '@mantine/core';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { mantineTheme } from '@/shared/ui';
import { RecapYearBanner } from '@/widgets/recap-year-banner';

describe('RecapYearBanner', () => {
  it('renders year and calls onOpen', () => {
    const onOpen = vi.fn();

    render(
      <MantineProvider theme={mantineTheme}>
        <RecapYearBanner year={2026} onOpen={onOpen} />
      </MantineProvider>,
    );

    expect(screen.getByText('Итоги 2026')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button'));
    expect(onOpen).toHaveBeenCalledOnce();
  });
});
