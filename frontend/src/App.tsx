import { useCallback, useState } from 'react';
import { Catalog } from './components/Catalog/Catalog';
import { Header } from './components/Header/Header';
import { RecapBanner } from './components/RecapBanner/RecapBanner';
import { RecapModal } from './components/RecapModal/RecapModal';
import styles from './App.module.css';

export default function App() {
  const [recapOpen, setRecapOpen] = useState(false);

  const openRecap = useCallback(() => setRecapOpen(true), []);
  const closeRecap = useCallback(() => setRecapOpen(false), []);

  return (
    <div className={styles.page}>
      <Header />
      <main className={styles.main}>
        <RecapBanner onOpen={openRecap} />
        <Catalog />
      </main>
      {recapOpen ? <RecapModal onClose={closeRecap} /> : null}
    </div>
  );
}
