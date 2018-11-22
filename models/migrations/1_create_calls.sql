-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "public"."webhook_calls" (
  "id" serial8,
  "timestamp" timestamp(6) NOT NULL,
  "host" text NOT NULL,
  "path" text NOT NULL,
  "method" text NOT NULL,
  "content_length" int4 NOT NULL,
  "from" text NOT NULL,
  "body" text,
  PRIMARY KEY ("id")
)
;
CREATE TABLE "public"."webhook_calls_headers" (
  "id" serial8,
  "call_id" int8 NOT NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("call_id") REFERENCES "public"."webhook_calls" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
;
CREATE TABLE "public"."webhook_calls_queries" (
  "id" serial8,
  "call_id" int8 NOT NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("call_id") REFERENCES "public"."webhook_calls" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
;
CREATE TABLE "public"."webhook_calls_form_params" (
  "id" serial8,
  "call_id" int8 NOT NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("call_id") REFERENCES "public"."webhook_calls" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
;
CREATE TABLE "public"."webhook_calls_files" (
  "id" serial8,
  "key" text NOT NULL,
  "call_id" int8 NOT NULL,
  "filename" text NOT NULL,
  "size" int8 NOT NULL,
  "contents" bytea NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("call_id") REFERENCES "public"."webhook_calls" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
;
CREATE TABLE "public"."webhook_calls_files_headers" (
  "id" serial8,
  "file_id" int8 NOT NULL,
  "key" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("file_id") REFERENCES "public"."webhook_calls_files" ("id") ON DELETE CASCADE ON UPDATE CASCADE
)
;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE "public"."webhook_calls_form_params";
DROP TABLE "public"."webhook_calls_queries";
DROP TABLE "public"."webhook_calls_header";
DROP TABLE "public"."webhook_calls";
