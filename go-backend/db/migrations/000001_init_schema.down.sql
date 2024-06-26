ALTER TABLE orders DROP CONSTRAINT fk_user_id;
ALTER TABLE order_items DROP CONSTRAINT fk_order_id;
ALTER TABLE order_items DROP CONSTRAINT fk_product_id;
ALTER TABLE reviews DROP CONSTRAINT fk_user_id;
ALTER TABLE reviews DROP CONSTRAINT fk_product_id;
ALTER TABLE blogs DROP CONSTRAINT fk_author;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS blogs;