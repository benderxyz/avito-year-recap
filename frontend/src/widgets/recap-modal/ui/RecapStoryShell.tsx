import { ActionIcon, Alert, Center, Loader, Text } from '@mantine/core';
import type { PreparedRecap, ResolvedTheme } from '@recap-engine/core';
import { Recap, type RecapEvent } from '@recap-engine/react';
import styles from './RecapStoryShell.module.css';

type RecapStoryShellProps = {
  title: string;
  onClose: () => void;
  isLoading: boolean;
  isError: boolean;
  error: unknown;
  prepared: PreparedRecap | null;
  theme: ResolvedTheme;
  onEvent: (event: RecapEvent) => void;
};

export function RecapStoryShell({
  title,
  onClose,
  isLoading,
  isError,
  error,
  prepared,
  theme,
  onEvent,
}: RecapStoryShellProps) {
  return (
    <div className={styles.storyShell}>
      <ActionIcon
        className={styles.closeButton}
        variant="subtle"
        color="gray"
        size="lg"
        radius="xl"
        aria-label="Закрыть"
        onClick={onClose}
      >
        ×
      </ActionIcon>

      <div className={styles.storyFrame}>
        <Text id="recap-modal-title" className={styles.storyTitle}>
          {title}
        </Text>

        {isLoading ? (
          <Center className={styles.storyState}>
            <Loader color="avito" />
          </Center>
        ) : null}

        {isError ? (
          <div className={styles.storyState}>
            <Alert color="red" title="Не удалось загрузить итоги года">
              {error instanceof Error ? error.message : 'Unknown error'}
            </Alert>
          </div>
        ) : null}

        {!isLoading && !isError && prepared ? (
          <div className={styles.recapRoot}>
            <Recap {...prepared} theme={theme} onEvent={onEvent} className={styles.recapWidget} />
          </div>
        ) : null}

        {!isLoading && !isError && !prepared ? (
          <Center className={styles.storyState}>
            <Text>Нет данных для отображения</Text>
          </Center>
        ) : null}
      </div>
    </div>
  );
}

export { styles as recapStoryShellStyles };
