CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "subscription" boolean NOT NULL DEFAULT False,
  "user_cart" bigint[],
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "product_name" varchar UNIQUE NOT NULL,
  "description" text NOT NULL,
  "price" float NOT NULL,
  "quantity" bigint NOT NULL DEFAULT 0,
  "discount" float,
  "rating" float,
  "size_options" varchar[],
  "color_options" varchar[],
  "category" varchar NOT NULL,
  "brand" varchar,
  "image_url" varchar[],
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "amount" float NOT NULL,
  "status" varchar NOT NULL DEFAULT 'pending',
  "shipping_address" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT fk_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "color" varchar,
  "size" varchar,
  "quantity" integer NOT NULL,
  CONSTRAINT fk_order_id FOREIGN KEY ("order_id") REFERENCES "orders" ("id"),
  CONSTRAINT fk_product_id FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

CREATE TABLE "reviews" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "rating" integer NOT NULL,
  "review" text NOT NULL,
  CONSTRAINT fk_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  CONSTRAINT fk_product_id FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

CREATE TABLE "blogs" (
  "id" bigserial PRIMARY KEY,
  "author" bigint NOT NULL,
  "title" varchar NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT fk_author FOREIGN KEY ("author") REFERENCES "users" ("id")
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("subscription");

CREATE INDEX ON "products" ("product_name");

CREATE INDEX ON "products" ("price");

CREATE INDEX ON "products" ("category");

CREATE INDEX ON "products" ("brand");

CREATE INDEX ON "products" ("category", "brand");

CREATE INDEX ON "orders" ("user_id");

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "order_items" ("product_id");

CREATE INDEX ON "reviews" ("product_id");

COMMENT ON COLUMN "users"."user_cart" IS 'list of product id in the cart';

COMMENT ON COLUMN "users"."role" IS 'user or admin';

COMMENT ON COLUMN "products"."discount" IS 'admins may have discount. Float of percentage ie 14.5';

COMMENT ON COLUMN "products"."rating" IS 'calculate when reviews is created. 1-5';

COMMENT ON COLUMN "products"."image_url" IS 'list of file paths to the product images';

