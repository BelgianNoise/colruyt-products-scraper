
-- EVERYTHING AUTH
ALTER USER postgres WITH PASSWORD 'xxxxxx';
REVOKE ALL ON SCHEMA public FROM PUBLIC;

CREATE ROLE viewer WITH
	NOLOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD 'xxxxxx';

CREATE ROLE webapp WITH
	LOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD 'xxxxxx';

GRANT viewer TO webapp;

CREATE ROLE creator WITH
	NOLOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1;

CREATE ROLE "go-backend" WITH
	LOGIN
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	NOREPLICATION
	CONNECTION LIMIT -1
	PASSWORD 'xxxxxx';

GRANT creator TO "go-backend";

-- CREATE NEW DATABASE
CREATE DATABASE "colruyt-products"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

-- ====================
-- SWITCH TO NEW DATABASE HERE
-- ====================

REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO viewer;
GRANT USAGE ON SCHEMA public TO creator;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public
GRANT SELECT ON TABLES TO viewer;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public
GRANT INSERT, SELECT, DELETE, UPDATE ON TABLES TO creator;

ALTER SCHEMA public
    RENAME TO products;

CREATE TABLE products.product
(
    id text NOT NULL,
    name text,
    long_name text,
    short_name text,
    content text,
    full_image text,
    square_image text,
    thumbnail text,
    commercial_article_number text,
    technical_article_number text,
    alcohol_volume text,
    country_of_origin text,
    fic_code text,
    is_biffe boolean,
    is_bio boolean,
    is_exclusively_sold_in_luxembourg boolean,
    is_new boolean,
    is_private_label boolean,
    is_weight_article boolean,
    order_unit text,
    recent_quantity_of_stock_units text,
    weightconversion_factor text,
    brand text,
    business_domain text,
    is_available boolean,
    seo_brand text,
    top_category_id text,
    top_category_name text,
    walk_route_sequence_number integer,
    PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);

ALTER TABLE IF EXISTS products.product
    OWNER to postgres;

CREATE TABLE products.price
(
    id SERIAL NOT NULL,
    product_id text NOT NULL,
    basic_price numeric NOT NULL,
    quantity_price numeric,
    quantity_price_quantity text,
    is_red_price boolean NOT NULL,
    in_promo boolean NOT NULL,
    in_pre_condition_promo boolean NOT NULL,
    is_price_available boolean NOT NULL,
    measurement_unit text NOT NULL,
    measurement_unit_price numeric NOT NULL,
    recommended_quantity text NOT NULL,
    "time" timestamp with time zone NOT NULL,
    promotion text DEFAULT NULL,
    promo_codes text DEFAULT NULL,
    PRIMARY KEY (product_id, "time"),
    FOREIGN KEY (product_id)
        REFERENCES products.product (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)
WITH (
    OIDS = FALSE
);

ALTER TABLE IF EXISTS products.price
    OWNER to postgres;

GRANT USAGE ON SEQUENCE products.price_id_seq TO creator;

CREATE TABLE products.promotion
(
    promotion_id text NOT NULL,
    active_end_date text NOT NULL,
    active_start_date text NOT NULL,
    benefit text NOT NULL,
    linked_products text NOT NULL,
    commercial_promotion_id text NOT NULL,
    folder_id text NOT NULL,
    max_times integer NOT NULL,
    personalised boolean NOT NULL,
    promotion_kind text NOT NULL,
    promotion_type text NOT NULL,
    publication_end_date text NOT NULL,
    publication_start_date text NOT NULL,
    PRIMARY KEY (promotion_id)
)
WITH (
    OIDS = FALSE
);

ALTER TABLE IF EXISTS products.promotion
    OWNER to postgres;
