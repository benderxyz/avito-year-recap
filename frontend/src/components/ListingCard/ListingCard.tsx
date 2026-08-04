import styles from './ListingCard.module.css';

export type Listing = {
  id: string;
  title: string;
  price: string;
  location: string;
  image: string;
};

type ListingCardProps = {
  listing: Listing;
};

export function ListingCard({ listing }: ListingCardProps) {
  return (
    <article className={styles.card}>
      <div className={styles.media}>
        <img src={listing.image} alt="" loading="lazy" />
      </div>
      <div className={styles.body}>
        <div className={styles.price}>{listing.price}</div>
        <div className={styles.title}>{listing.title}</div>
        <div className={styles.meta}>{listing.location}</div>
      </div>
    </article>
  );
}
