import styles from './RecapYearBanner.module.css';

type RecapYearBannerProps = {
  year: number;
  onOpen: () => void;
};

export function RecapYearBanner({ year, onOpen }: RecapYearBannerProps) {
  return (
    <button type="button" className={styles.banner} onClick={onOpen}>
      <div className={styles.copy}>
        <span className={styles.eyebrow}>Персональная подборка</span>
        <strong className={styles.title}>Итоги {year}</strong>
        <span className={styles.text}>Посмотрите, каким был ваш год на Авито</span>
      </div>
      <span className={styles.cta}>Открыть</span>
    </button>
  );
}
