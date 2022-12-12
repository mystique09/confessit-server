CREATE TABLE "user" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

CREATE TABLE "message" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "to" uuid,
  "content" varchar,
  "seen" boolean DEFAULT false,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz DEFAULT (now())
);

ALTER TABLE "message" ADD FOREIGN KEY ("to") REFERENCES "user" ("id");