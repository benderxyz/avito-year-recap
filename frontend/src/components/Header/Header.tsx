import styles from './Header.module.css';

export function Header() {
  return (
    <header className={styles.header}>
      <div className={styles.inner}>
        <a className={styles.logo} href="/" aria-label="Авито">
          Авито
        </a>
        <div className={styles.search}>
          <input
            className={styles.input}
            type="search"
            placeholder="Поиск по объявлениям"
            defaultValue=""
            readOnly
          />
          <button type="button" className={styles.searchBtn}>
            Найти
          </button>
        </div>
        <nav className={styles.nav} aria-label="Профиль">
          <span className={styles.navItem}>Избранное</span>
          <span className={styles.navItem}>Сообщения</span>
          <span className={styles.avatar}>А</span>
        </nav>
      </div>
    </header>
  );
}
