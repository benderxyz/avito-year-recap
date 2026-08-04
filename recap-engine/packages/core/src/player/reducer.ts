import {
  EPlayerActionType,
  EPlayerPhase,
  type PlayerAction,
  type PlayerState,
} from '../types/player';

export function createInitialPlayerState(total: number, initialIndex = 0): PlayerState {
  return {
    index: Math.min(Math.max(initialIndex, 0), Math.max(total - 1, 0)),
    direction: 1,
    phase: EPlayerPhase.Enter,
    total,
  };
}

export function playerReducer(state: PlayerState, action: PlayerAction): PlayerState {
  switch (action.type) {
    case EPlayerActionType.SetTotal:
      return {
        ...state,
        total: action.total,
        index: Math.min(state.index, Math.max(action.total - 1, 0)),
      };

    case EPlayerActionType.SetPhase:
      return { ...state, phase: action.phase };

    case EPlayerActionType.Next: {
      if (state.index >= state.total - 1) return state;
      return {
        ...state,
        index: state.index + 1,
        direction: 1,
        phase: EPlayerPhase.Enter,
      };
    }

    case EPlayerActionType.Prev: {
      if (state.index <= 0) return state;
      return {
        ...state,
        index: state.index - 1,
        direction: -1,
        phase: EPlayerPhase.Enter,
      };
    }

    case EPlayerActionType.GoTo: {
      if (action.index < 0 || action.index >= state.total || action.index === state.index) {
        return state;
      }

      return {
        ...state,
        index: action.index,
        direction: action.index > state.index ? 1 : -1,
        phase: EPlayerPhase.Enter,
      };
    }

    default: {
      const _exhaustive: never = action;
      return _exhaustive;
    }
  }
}
