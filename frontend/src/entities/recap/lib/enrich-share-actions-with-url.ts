import {
  ESceneActionType,
  ESceneType,
  type RecapPayload,
  type SceneAction,
} from '@recap-engine/core';

type RecapPayloadWithShareUrl = RecapPayload & {
  features?: RecapPayload['features'] & { shareUrl?: string };
};

function enrichShareAction(action: SceneAction, shareUrl: string): SceneAction {
  if (action.type !== ESceneActionType.Share || action.share.url) {
    return action;
  }

  return {
    ...action,
    share: {
      ...action.share,
      url: shareUrl,
    },
  };
}

export function enrichShareActionsWithUrl(payload: RecapPayload): RecapPayload {
  const shareUrl = (payload as RecapPayloadWithShareUrl).features?.shareUrl;
  if (!shareUrl) {
    return payload;
  }

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
