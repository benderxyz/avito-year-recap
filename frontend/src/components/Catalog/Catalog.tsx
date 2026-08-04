import { ListingCard, type Listing } from '../ListingCard/ListingCard';
import styles from './Catalog.module.css';

const LISTINGS: Listing[] = [
  {
    id: '1',
    title: 'iPhone 15 Pro 256GB',
    price: '89 900 ₽',
    location: 'Москва',
    image: 'https://picsum.photos/seed/iphone/400/300',
  },
  {
    id: '2',
    title: 'Велосипед горный Trek',
    price: '34 500 ₽',
    location: 'Санкт-Петербург',
    image: 'https://picsum.photos/seed/bike/400/300',
  },
  {
    id: '3',
    title: 'Диван угловой серый',
    price: '27 000 ₽',
    location: 'Казань',
    image: 'https://picsum.photos/seed/sofa/400/300',
  },
  {
    id: '4',
    title: 'MacBook Air M2',
    price: '78 000 ₽',
    location: 'Москва',
    image: 'https://picsum.photos/seed/macbook/400/300',
  },
  {
    id: '5',
    title: 'Кроссовки Nike Dunk',
    price: '9 900 ₽',
    location: 'Екатеринбург',
    image: 'https://picsum.photos/seed/shoes/400/300',
  },
  {
    id: '6',
    title: 'Гитара Fender Stratocaster',
    price: '52 000 ₽',
    location: 'Новосибирск',
    image: 'https://picsum.photos/seed/guitar/400/300',
  },
  {
    id: '7',
    title: 'Пылесос Dyson V15',
    price: '41 000 ₽',
    location: 'Москва',
    image: 'https://picsum.photos/seed/dyson/400/300',
  },
  {
    id: '8',
    title: 'Стол письменный дуб',
    price: '12 500 ₽',
    location: 'Краснодар',
    image: 'https://picsum.photos/seed/desk/400/300',
  },
];

export function Catalog() {
  return (
    <section className={styles.section}>
      <div className={styles.head}>
        <h1 className={styles.title}>Рекомендации для вас</h1>
        <p className={styles.subtitle}>Объявления рядом и по интересам</p>
      </div>
      <div className={styles.grid}>
        {LISTINGS.map((listing) => (
          <ListingCard key={listing.id} listing={listing} />
        ))}
      </div>
    </section>
  );
}
