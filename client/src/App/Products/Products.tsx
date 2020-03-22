import React from "react";
import { IProduct } from "./types";

const classes = ["a", "b", "c"];

const Products: React.FC<{ products: IProduct[] }> = ({ products }) => (
  <main className="products-feature">
    <header>Store</header>
    <div className="products">
      {products.map((product, index) => {
        return (
          index < 3 && (
            <div
              key={product.id}
              className={`products__card ${classes[index]}`}
            >
              <b className="products__card-name">{product.name}</b>
              <p className="products__card-description">
                {product.description}
              </p>
              <p className="products__card-price">{product.price}</p>
            </div>
          )
        );
      })}
    </div>
  </main>
);

export default Products;
