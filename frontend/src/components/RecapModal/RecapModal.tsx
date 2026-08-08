import {
  createTheme,
  ERecapEventType,
  ESceneActionType,
  metricList,
  prepareRecap,
  Recap,
  type RecapEvent,
  type RecapPayload,
  registerScene,
} from '@recap-engine/react';
import { useEffect, useMemo } from 'react';
import { mockRecapPayload } from '../../mock/recapPayload';
import styles from './RecapModal.module.css';

const theme = createTheme();

registerScene<RecapPayload>('top-categories', ({ data, props }) => {
  const metricKey = typeof props?.metricKey === 'string' ? props.metricKey : 'favoriteCategories';
  const title = typeof props?.title === 'string' ? props.title : 'Топ категорий';
  const items = metricList(data.metrics, metricKey);

  return (
    <div className={styles.customScene}>
      <p className={styles.customEyebrow}>За год</p>
      <h2 className={styles.customTitle}>{title}</h2>
      <ul className={styles.customList}>
        {items.map((item, index) => (
          <li key={item.id} className={styles.customItem}>
            <span className={styles.customRank}>{index + 1}</span>
            <span className={styles.customLabel}>{item.label}</span>
            <span className={styles.customValue}>{item.value ?? 0}</span>
          </li>
        ))}
      </ul>
    </div>
  );
});

type RecapModalProps = {
  onClose: () => void;
};

export function RecapModal({ onClose }: RecapModalProps) {
  const prepared = useMemo(() => prepareRecap(mockRecapPayload), []);

  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
    };

    document.body.style.overflow = 'hidden';
    window.addEventListener('keydown', onKeyDown);

    return () => {
      document.body.style.overflow = '';
      window.removeEventListener('keydown', onKeyDown);
    };
  }, [onClose]);

  const onEvent = (event: RecapEvent) => {
    if (event.type === ERecapEventType.Complete) {
      onClose();
      return;
    }

    if (
      event.type === ERecapEventType.Action &&
      event.action.type === ESceneActionType.Custom &&
      event.action.id === 'close-recap'
    ) {
      onClose();
    }
  };

  return (
    <div className={styles.overlay} role="dialog" aria-modal="true" aria-label="Итоги 2026">
      <button type="button" className={styles.backdrop} aria-label="Закрыть" onClick={onClose} />
      <div className={styles.frame}>
        <button type="button" className={styles.close} onClick={onClose} aria-label="Закрыть">
          ×
        </button>
        <Recap {...prepared} autoplay theme={theme} onEvent={onEvent} className={styles.recap} />
      </div>
    </div>
  );
}
