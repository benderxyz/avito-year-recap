import { useQuery } from '@tanstack/react-query';
import { filterProductsByPrice, mapDummyJsonProduct } from '@/entities/product';
import {
  getProductsSkip,
  PRODUCTS_PAGE_SIZE,
  parsePriceFilter,
  useProductFilters,
} from '@/features/product-filters';
import { fetchProductCategories, fetchProducts } from '@/shared/api';

export function useProductsCatalog() {
  const [filters] = useProductFilters();

  const productsQuery = useQuery({
    queryKey: [
      'products',
      {
        q: filters.q,
        category: filters.category,
        sort: filters.sort,
        order: filters.order,
        page: filters.page,
      },
    ],
    queryFn: () =>
      fetchProducts({
        q: filters.q || undefined,
        category: filters.category || undefined,
        sortBy: filters.sort,
        order: filters.order,
        limit: PRODUCTS_PAGE_SIZE,
        skip: getProductsSkip(filters.page),
      }),
  });

  const categoriesQuery = useQuery({
    queryKey: ['product-categories'],
    queryFn: fetchProductCategories,
    staleTime: 60 * 60 * 1000,
  });

  const products = (productsQuery.data?.products ?? [])
    .map(mapDummyJsonProduct)
    .filter(
      (product) =>
        filterProductsByPrice(
          [product],
          parsePriceFilter(filters.priceMin),
          parsePriceFilter(filters.priceMax),
        ).length > 0,
    );

  return {
    filters,
    products,
    total: productsQuery.data?.total ?? 0,
    categories: categoriesQuery.data ?? [],
    isLoading: productsQuery.isLoading,
    isError: productsQuery.isError,
    error: productsQuery.error,
  };
}
