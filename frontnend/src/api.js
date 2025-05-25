export const fetchProducts = async () => {
  const res = await fetch("/api/products");
  const data = await res.json();
  return data.products;
};
