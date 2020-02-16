# Performing migrations

```makefile
make db # runs database in the background
make migration create_products_table
```

Then add sql to both up & down migrations files found under: `./api/internal/schema/migrations/`.

```sql
-- 000001_create_products_table.up.sql

CREATE TABLE products (
    id UUID not null unique,
    name varchar(100) not null,
    price real not null,
    description varchar(100) not null,
    created timestamp without time zone default (now() at time zone 'utc')
);
```

```sql
-- 000001_create_products_table.down.sql

DROP TABLE IF EXISTS products;
```

Make another migration to add tags to products:

```
make migration add_tags_to_products
```

```sql

-- 000002_add_tags_to_products.up.sql

ALTER TABLE products
ADD COLUMN tags varchar(255);
```

```sql
-- 000002_add_tags_to_products.down.sql

ALTER TABLE products
DROP Column tags;
```

Migrate up to the latest migration

```makefile
make up # you can migrate down with "make down"
```

Display which version you have selected:

```makefile
make version
```

[Learn more about my go-migrate postgres helper](https://github.com/ivorscott/go-migrate-postgres-helper)

Next we need to seed the database:

```makefile
make seed products
```

This adds an empty products.sql seed file found under `./api/internal/schema/`. Add some rows:

```sql
-- ./api/internal/schema/products.sql

INSERT INTO products (id, name, price, description) VALUES
('cbef5139-323f-48b8-b911-dc9be7d0bc07','Xbox One X', 499.00, 'Eighth-generation home video game console developed by Microsoft.'),
('ce93a886-3a0e-456b-b7f5-8652d2de1e8f','Playsation 4', 299.00, 'Eighth-generation home video game console developed by Sony Interactive Entertainment.'),
('faa25b57-7031-4b37-8a89-de013418deb0','Nintendo Switch', 299.00, 'Hybrid console that can be used as a stationary and portable device developed by Nintendo.')
ON CONFLICT DO NOTHING;
```

Appending "ON CONFLICT DO NOTHING;" to the end of the sql command prevents conflicts if the seed file is executed to the database more than once. **Note:** This behavior works because the products table has at least one table column with a unique constraint.

Finally, add the products seed file to the database

```
make insert products
```

Enter the database and examine its state

```makefile
make debug-db
```
