export enum EButtonVariant {
  Primary = 'primary',
  Secondary = 'secondary',
  Ghost = 'ghost',
}

export enum ESceneActionType {
  Next = 'next',
  Prev = 'prev',
  Link = 'link',
  Share = 'share',
  GoTo = 'goto',
  Custom = 'custom',
}

export enum EShareKind {
  Story = 'story',
  Link = 'link',
}

export enum ELinkTarget {
  Blank = '_blank',
  Self = '_self',
}

export enum ERecapEventType {
  SceneEnter = 'scene_enter',
  SceneExit = 'scene_exit',
  Action = 'action',
  Complete = 'complete',
  Share = 'share',
}

type NextSceneAction = {
  type: ESceneActionType.Next;
  label?: string;
  variant?: EButtonVariant;
};

type PreviousSceneAction = {
  type: ESceneActionType.Prev;
  label?: string;
  variant?: EButtonVariant;
};

type LinkSceneAction = {
  type: ESceneActionType.Link;
  label: string;
  href: string;
  target?: ELinkTarget;
  variant?: EButtonVariant;
};

type ShareSceneAction = {
  type: ESceneActionType.Share;
  label: string;
  share: {
    kind: EShareKind;
    title?: string;
    text?: string;
    url?: string;
  };
  variant?: EButtonVariant;
};

type GoToSceneAction = {
  type: ESceneActionType.GoTo;
  label: string;
  sceneId: string;
  variant?: EButtonVariant;
};

type CustomSceneAction = {
  type: ESceneActionType.Custom;
  label: string;
  id: string;
  variant?: EButtonVariant;
};

export type SceneAction =
  | NextSceneAction
  | PreviousSceneAction
  | LinkSceneAction
  | ShareSceneAction
  | GoToSceneAction
  | CustomSceneAction;

type SceneEnterRecapEvent = {
  type: ERecapEventType.SceneEnter;
  sceneId: string;
  index: number;
};

type SceneExitRecapEvent = {
  type: ERecapEventType.SceneExit;
  sceneId: string;
  index: number;
};

type ActionRecapEvent = {
  type: ERecapEventType.Action;
  sceneId: string;
  action: SceneAction;
};

type CompleteRecapEvent = { type: ERecapEventType.Complete };

type ShareRecapEvent = {
  type: ERecapEventType.Share;
  sceneId: string;
  payload: {
    kind: EShareKind;
    title?: string;
    text?: string;
    url?: string;
  };
};

export type RecapEvent =
  | SceneEnterRecapEvent
  | SceneExitRecapEvent
  | ActionRecapEvent
  | CompleteRecapEvent
  | ShareRecapEvent;
