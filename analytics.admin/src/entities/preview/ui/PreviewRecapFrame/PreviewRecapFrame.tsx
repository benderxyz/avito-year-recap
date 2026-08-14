import { Alert, Center, Loader, Text } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { PreparedRecap, RecapEngineModule } from '../../lib/load-recap-engine';
import styles from './PreviewRecapFrame.module.css';

type PreviewRecapFrameProps = {
  engine: RecapEngineModule | null;
  engineVersion: string | null;
  prepared: PreparedRecap | null;
  theme: unknown;
  autoplay: boolean;
  loop: boolean;
  gestures: boolean;
  tapNav: boolean;
  holdToPause: boolean;
  reducedMotion: boolean;
  isLoading: boolean;
  error: string | null;
};

export default function PreviewRecapFrame({
  engine,
  engineVersion,
  prepared,
  theme,
  autoplay,
  loop,
  gestures,
  tapNav,
  holdToPause,
  reducedMotion,
  isLoading,
  error,
}: PreviewRecapFrameProps) {
  const { t } = useTranslation();
  const Recap = engine?.Recap;

  return (
    <div className={styles.frame}>
      {isLoading ? (
        <Center className={styles.state}>
          <Loader />
        </Center>
      ) : null}

      {!isLoading && error ? (
        <div className={styles.state}>
          <Alert color="red">{error}</Alert>
        </div>
      ) : null}

      {!isLoading && !error && Recap && prepared ? (
        <Recap
          key={engineVersion}
          {...prepared}
          theme={theme}
          autoplay={autoplay}
          loop={loop}
          gestures={gestures}
          tapNav={tapNav}
          holdToPause={holdToPause}
          reducedMotion={reducedMotion}
          className={styles.widget}
        />
      ) : null}

      {!isLoading && !error && (!Recap || !prepared) ? (
        <Center className={styles.state}>
          <Text>{t('preview.empty')}</Text>
        </Center>
      ) : null}
    </div>
  );
}
