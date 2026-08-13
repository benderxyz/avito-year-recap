import type { CSSProperties } from 'react';
import styles from './RecapModalOverlayCircles.module.css';

type RecapModalOverlayCirclesProps = {
  style?: CSSProperties;
};

type CircleConfig = {
  id: string;
  size: number;
  color: string;
  top?: string;
  right?: string;
  bottom?: string;
  left?: string;
  enterDelay: number;
  pulseDelay: number;
};

const CIRCLES: CircleConfig[] = [
  {
    id: 'top-left',
    size: 220,
    color: '#00AAFF',
    top: '-60px',
    left: '-40px',
    enterDelay: 0,
    pulseDelay: 0,
  },
  {
    id: 'right-top',
    size: 140,
    color: '#FF6163',
    top: '12%',
    right: '-50px',
    enterDelay: 70,
    pulseDelay: 400,
  },
  {
    id: 'bottom-left',
    size: 250,
    color: '#97CF26',
    bottom: '-80px',
    left: '8%',
    enterDelay: 140,
    pulseDelay: 800,
  },
  {
    id: 'bottom-right',
    size: 100,
    color: '#A169F7',
    bottom: '15%',
    right: '5%',
    enterDelay: 210,
    pulseDelay: 1200,
  },
  {
    id: 'left-center',
    size: 180,
    color: '#00AAFF',
    top: '55%',
    left: '-90px',
    enterDelay: 280,
    pulseDelay: 1600,
  },
];

export function RecapModalOverlayCircles({ style }: RecapModalOverlayCirclesProps) {
  return (
    <div className={styles.root} style={style} aria-hidden="true">
      {CIRCLES.map((circle) => (
        <div
          key={circle.id}
          className={styles.circle}
          style={
            {
              '--circle-size': `${circle.size}px`,
              '--circle-color': circle.color,
              '--circle-enter-delay': `${circle.enterDelay}ms`,
              '--circle-pulse-delay': `${550 + circle.enterDelay + circle.pulseDelay}ms`,
              top: circle.top,
              right: circle.right,
              bottom: circle.bottom,
              left: circle.left,
            } as CSSProperties
          }
        />
      ))}
    </div>
  );
}
