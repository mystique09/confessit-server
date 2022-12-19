CREATE TYPE "satisfaction" AS ENUM (
  'LIKE',
  'DISLIKE'
);

CREATE TABLE "users" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" date NOT NULL DEFAULT (now()),
  "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "user_identities" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "user_id" uuid NOT NULL,
  "identity_hash" uuid NOT NULL
);

CREATE TABLE "messages" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "receiver_id" uuid NOT NULL,
  "content" varchar NOT NULL,
  "seen" boolean NOT NULL DEFAULT false,
  "created_at" date NOT NULL DEFAULT (now()),
  "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "content" varchar NOT NULL,
  "user_identity_id" uuid NOT NULL,
  "created_at" date NOT NULL DEFAULT (now()),
  "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "content" varchar NOT NULL,
  "user_identity_id" uuid NOT NULL,
  "post_id" uuid NOT NULL,
  "parent_id" uuid,
  "created_at" date NOT NULL DEFAULT (now()),
  "updated_at" date NOT NULL DEFAULT (now())
);

CREATE TABLE "likes" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "user_identity_id" uuid NOT NULL,
  "post_id" uuid NOT NULL,
  "type" satisfaction
);

CREATE TABLE "comment_likes" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "user_identity_id" uuid NOT NULL,
  "comment_id" uuid NOT NULL,
  "type" satisfaction
);

CREATE TABLE "sessions" (
  "id" uuid UNIQUE PRIMARY KEY NOT NULL,
  "user_id" uuid NOT NULL,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "created_at" date NOT NULL DEFAULT (now()),
  "expires_at" date NOT NULL
);

ALTER TABLE "user_identities" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "messages" ADD FOREIGN KEY ("receiver_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "posts" ADD FOREIGN KEY ("user_identity_id") REFERENCES "user_identities" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "comments" ADD FOREIGN KEY ("user_identity_id") REFERENCES "user_identities" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "comments" ADD FOREIGN KEY ("parent_id") REFERENCES "comments" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "likes" ADD FOREIGN KEY ("user_identity_id") REFERENCES "user_identities" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "likes" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "comment_likes" ADD FOREIGN KEY ("user_identity_id") REFERENCES "user_identities" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "comment_likes" ADD FOREIGN KEY ("comment_id") REFERENCES "comments" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;