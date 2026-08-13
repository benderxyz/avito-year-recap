import { describe, expect, it } from 'vitest';
import { storyFormDefaults, storyFormSchema, toStoryCreate } from './form-schema';

const validStory = {
  ...storyFormDefaults,
  id: 'year_intro',
  payload: JSON.stringify({ id: 'year_intro', type: 'intro', title: 'Your year' }),
};

describe('storyFormSchema', () => {
  it('rejects an empty id', () => {
    const result = storyFormSchema.safeParse(storyFormDefaults);

    expect(result.success).toBe(false);
  });

  it('accepts a valid story', () => {
    const result = storyFormSchema.safeParse(validStory);

    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.id).toBe('year_intro');
    }
  });

  it('rejects invalid payload JSON', () => {
    const result = storyFormSchema.safeParse({
      ...validStory,
      payload: '{',
    });

    expect(result.success).toBe(false);
  });

  it('omits when when metric is empty and syncs payload id/type', () => {
    const created = toStoryCreate({
      ...validStory,
      sceneType: 'stat',
      payload: JSON.stringify({ title: 'Orders' }),
    });

    expect(created.when).toBeNull();
    expect(created.payload).toEqual({
      title: 'Orders',
      id: 'year_intro',
      type: 'stat',
    });
  });
});
