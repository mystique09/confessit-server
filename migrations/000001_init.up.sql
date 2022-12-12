CREATE TABLE "user" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" date DEFAULT (now()),
  "updated_at" date DEFAULT (now())
);

CREATE TABLE "message" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "receiver_id" uuid NOT NULL,
  "content" varchar,
  "seen" boolean DEFAULT false,
  "created_at" date DEFAULT (now()),
  "updated_at" date DEFAULT (now())
);

ALTER TABLE "message" ADD FOREIGN KEY ("receiver_id") REFERENCES "user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;