export enum EPlayerPhase {
  Enter = 'enter',
  Active = 'active',
  Exit = 'exit',
}

export type PlayerState = {
  index: number;
  direction: 1 | -1;
  phase: EPlayerPhase;
  total: number;
};

export enum EPlayerActionType {
  Next = 'next',
  Prev = 'prev',
  GoTo = 'go_to',
  SetPhase = 'set_phase',
  SetTotal = 'set_total',
}

type NextPlayerAction = { type: EPlayerActionType.Next };

type PreviousPlayerAction = { type: EPlayerActionType.Prev };

type GoToPlayerAction = { type: EPlayerActionType.GoTo; index: number };

type SetPhasePlayerAction = {
  type: EPlayerActionType.SetPhase;
  phase: EPlayerPhase;
};

type SetTotalPlayerAction = {
  type: EPlayerActionType.SetTotal;
  total: number;
};

export type PlayerAction =
  | NextPlayerAction
  | PreviousPlayerAction
  | GoToPlayerAction
  | SetPhasePlayerAction
  | SetTotalPlayerAction;
