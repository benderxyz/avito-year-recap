import { describe, expect, it } from '@jest/globals';
import { renderHook } from '@testing-library/react';
import { createRecapValue, recapWrapper } from '../../test/recapTestUtils';
import { useRecap } from './RecapProvider';

describe('RecapProvider', () => {
  it('provides the complete recap value to consumers', () => {
    const value = createRecapValue({
      data: { userId: '42' },
      progress: 0.75,
      reducedMotion: true,
    });

    const { result } = renderHook(() => useRecap<{ userId: string }>(), {
      wrapper: recapWrapper(value),
    });

    expect(result.current).toBe(value);
    expect(result.current.data.userId).toBe('42');
    expect(result.current.progress).toBe(0.75);
    expect(result.current.reducedMotion).toBe(true);
  });

  it('throws a clear error outside the provider', () => {
    expect(() => renderHook(() => useRecap())).toThrow('useRecap must be used within <Recap>');
  });
});
