import type { Badge } from '../types/badge';
import type { RecapPayload } from '../types/payload';
import type { SceneDefinition } from '../types/scenes';
import { buildScenesFromStory } from './buildScenesFromStory';

export type PreparedRecap = {
  data: RecapPayload;
  scenes: SceneDefinition<RecapPayload>[];
  badges: Badge[];
  locale: string;
};

/** Props bag for `<Recap {...prepareRecap(payload)} theme={…} />`. */
export function prepareRecap(payload: RecapPayload): PreparedRecap {
  return {
    data: payload,
    scenes: buildScenesFromStory(payload),
    locale: payload.meta.locale,
    badges: (payload.badges ?? []).map((b) => ({
      id: b.id,
      title: b.title,
      description: b.description,
      icon: b.iconUrl,
    })),
  };
}
