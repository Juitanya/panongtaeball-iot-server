-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create "discord_toggle_histories" table
CREATE TABLE "public"."discord_toggle_histories" ("id" bigserial NOT NULL, "action_by" bigserial NOT NULL, "payload" json NOT NULL, "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));
-- Create index "ix_discord_toggle_histories_action_by" to table: "discord_toggle_histories"
CREATE INDEX "ix_discord_toggle_histories_action_by" ON "public"."discord_toggle_histories" ("action_by");
-- Create "function_histories" table
CREATE TABLE "public"."function_histories" ("id" bigserial NOT NULL, "associate_with" character varying NOT NULL, "called_by_function" character varying NOT NULL, "line" bigint NOT NULL, "file_location" text NOT NULL, "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));
-- Create "user_otps" table
CREATE TABLE "public"."user_otps" ("id" bigserial NOT NULL, "user_id" bigserial NOT NULL, "otp" character varying NOT NULL, "expired_at" timestamp(3) NOT NULL DEFAULT (CURRENT_TIMESTAMP + '01:00:00'::interval), "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP);
-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "email" character varying NULL, "phone_number" character varying NULL, "password_hash" character varying NOT NULL, "admin" boolean NOT NULL DEFAULT false, "verify_by" integer NOT NULL DEFAULT 0, "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP, "deleted_at" timestamp(3) NULL, PRIMARY KEY ("id"));
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_phone" to table: "users"
CREATE INDEX "idx_users_phone" ON "public"."users" ("phone_number");
-- Create index "unique_users_email" to table: "users"
CREATE UNIQUE INDEX "unique_users_email" ON "public"."users" ("email") WHERE ((deleted_at IS NULL) AND ((email)::text <> ''::text));
-- Create index "unique_users_phone" to table: "users"
CREATE UNIQUE INDEX "unique_users_phone" ON "public"."users" ("phone_number") WHERE ((deleted_at IS NULL) AND ((phone_number)::text <> ''::text));
