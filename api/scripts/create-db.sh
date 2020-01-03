#!/bin/bash
set -e

if [ ! -f "/docker-entrypoint-initdb.d/backup.sql" ]; then
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE products (
        id serial primary key,
	    name varchar(100) not null,
	    price real not null,
	    description varchar(100) not null,
	    created timestamp without time zone default (now() at time zone 'utc')
    );
    INSERT INTO products (name, price, description)
    VALUES
        ('Xbox One XXXXXX', 499.00, 'Eighth-generation home video game console developed by Microsoft.'),
        ('Playsation 444444', 299.00, 'Eighth-generation home video game console developed by Sony Interactive Entertainment.'),
        ('Nintendo Switch', 299.00, 'Hybrid console that can be used as a stationary and portable device developed by Nintendo.');
    SELECT name, price from products;
EOSQL
else 
echo -e /docker-entrypoint-initdb.d/create-db.sh was ignored because a backup exists "\n"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
SELECT name, price from products;
EOSQL
fi
