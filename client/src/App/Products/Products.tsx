import React from "react";

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
}

interface Props {
  products: Product[];
}
export const Products: React.FC<Props> = ({ products }) => (
  <main className="products-feature">
    <header>Store</header>
    <div className="products">
      {products.map(product => (
        <div key={product.id} className="products__card">
          <b className="products__card-name">{product.name}</b>
          <p className="products__card-description">{product.description}</p>
          <p className="products__card-price">{product.price}</p>
        </div>
      ))}
    </div>
  </main>
);
