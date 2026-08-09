import { Modal } from '@mantine/core';
import {
  metricList,
  prepareRecap,
  type RecapEvent,
  type RecapPayload,
  registerScene,
} from '@recap-engine/react';
import { useEffect, useMemo } from 'react';
import {
  buildRecapTheme,
  shouldCloseRecapOnEvent,
  useRecapQuery,
  useResolvedRecapColorScheme,
} from '@/entities/recap';
import styles from './RecapModal.module.css';
import { RecapStoryShell, recapStoryShellStyles } from './RecapStoryShell';

let topCategoriesSceneRegistered = false;

function ensureTopCategoriesScene() {
  if (topCategoriesSceneRegistered) {
    return;
  }

  registerScene<RecapPayload>('top-categories', ({ data, props }) => {
    const metricKey = typeof props?.metricKey === 'string' ? props.metricKey : 'favoriteCategories';
    const title = typeof props?.title === 'string' ? props.title : 'Топ категорий';
    const items = metricList(data.metrics, metricKey);

    return (
      <div className={recapStoryShellStyles.customScene}>
        <p className={recapStoryShellStyles.customEyebrow}>За год</p>
        <h2 className={recapStoryShellStyles.customTitle}>{title}</h2>
        <ul className={recapStoryShellStyles.customList}>
          {items.map((item, index) => (
            <li key={item.id} className={recapStoryShellStyles.customItem}>
              <span className={recapStoryShellStyles.customRank}>{index + 1}</span>
              <span className={recapStoryShellStyles.customLabel}>{item.label}</span>
              <span className={recapStoryShellStyles.customValue}>{item.value ?? 0}</span>
            </li>
          ))}
        </ul>
      </div>
    );
  });

  topCategoriesSceneRegistered = true;
}

type RecapModalProps = {
  userId: string;
  opened: boolean;
  onClose: () => void;
};

export function RecapModal({ userId, opened, onClose }: RecapModalProps) {
  const { data, isLoading, isError, error } = useRecapQuery(userId, opened);
  const resolvedColorScheme = useResolvedRecapColorScheme();
  const theme = useMemo(() => buildRecapTheme(resolvedColorScheme), [resolvedColorScheme]);

  useEffect(() => {
    ensureTopCategoriesScene();
  }, []);

  const prepared = useMemo(() => (data ? prepareRecap(data) : null), [data]);

  const onEvent = (event: RecapEvent) => {
    if (shouldCloseRecapOnEvent(event)) {
      onClose();
    }
  };

  return (
    <Modal
      opened={opened}
      onClose={onClose}
      centered
      withCloseButton={false}
      padding={0}
      radius="xl"
      size="auto"
      overlayProps={{
        blur: 12,
        backgroundOpacity: 0.62,
      }}
      classNames={{
        inner: styles.inner,
        content: styles.content,
        body: styles.body,
      }}
      aria-labelledby="recap-modal-title"
    >
      <RecapStoryShell
        title={`Итоги ${data?.meta.year ?? ''}`}
        onClose={onClose}
        isLoading={isLoading}
        isError={isError}
        error={error}
        prepared={prepared}
        theme={theme}
        onEvent={onEvent}
      />
    </Modal>
  );
}
