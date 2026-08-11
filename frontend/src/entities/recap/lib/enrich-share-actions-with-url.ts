import {
  ESceneActionType,
  ESceneType,
  EShareKind,
  type RecapPayload,
  type SceneAction,
} from '@recap-engine/core';

type RecapPayloadWithShareUrl = RecapPayload & {
  features?: RecapPayload['features'] & { shareUrl?: string };
};

function enrichShareAction(action: SceneAction, shareUrl: string | undefined): SceneAction {
  if (action.type !== ESceneActionType.Share || action.share.kind !== EShareKind.Link) {
    return action;
  }

  const url = action.share.url || shareUrl;
  if (!url) {
    return action;
  }

  return {
    ...action,
    share: {
      kind: EShareKind.Link,
      url,
    },
  };
}

export function enrichShareActionsWithUrl(payload: RecapPayload): RecapPayload {
  const shareUrl = (payload as RecapPayloadWithShareUrl).features?.shareUrl;
  return {
    ...payload,
    story: payload.story.map((item) => {
      if (item.type !== ESceneType.Outro || !item.actions) {
        return item;
      }

      return {
        ...item,
        actions: item.actions.map((action) => enrichShareAction(action, shareUrl)),
      };
    }),
  };
}
