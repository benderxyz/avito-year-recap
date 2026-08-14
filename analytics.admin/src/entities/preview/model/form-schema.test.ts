import { describe, expect, it } from 'vitest';
import { previewFormDefaults, previewFormSchema } from './form-schema';

const validPreview = {
  ...previewFormDefaults,
  reactVersion: '2.0.1',
};

describe('previewFormSchema', () => {
  it('rejects an empty react version', () => {
    const result = previewFormSchema.safeParse(previewFormDefaults);

    expect(result.success).toBe(false);
  });

  it('accepts a null seed', () => {
    const result = previewFormSchema.safeParse(validPreview);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.seed).toBeNull();
    }
  });
});
