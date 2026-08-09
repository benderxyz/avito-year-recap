export type Product = {
  id: string;
  title: string;
  price: number;
  priceLabel: string;
  category: string;
  thumbnail: string;
  rating: number;
  brand: string;
};

export type ProductFilters = {
  q: string;
  category: string;
  sort: string;
  order: 'asc' | 'desc';
  priceMin: number | null;
  priceMax: number | null;
  page: number;
};
