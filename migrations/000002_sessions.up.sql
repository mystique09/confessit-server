CREATE TABLE "session" (
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

ALTER TABLE "session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;