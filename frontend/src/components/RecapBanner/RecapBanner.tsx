import styles from './RecapBanner.module.css';

type RecapBannerProps = {
  onOpen: () => void;
};

export function RecapBanner({ onOpen }: RecapBannerProps) {
  return (
    <button type="button" className={styles.banner} onClick={onOpen}>
      <div className={styles.copy}>
        <span className={styles.eyebrow}>Персональная подборка</span>
        <strong className={styles.title}>Итоги 2026</strong>
        <span className={styles.text}>Посмотрите, каким был ваш год на Авито</span>
      </div>
      <span className={styles.cta}>Открыть</span>
    </button>
  );
}
