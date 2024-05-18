CREATE TYPE "role" AS ENUM (
  'it',
  'nurse'
);


CREATE TABLE "staff" (
     "id" uuid PRIMARY KEY,
     "nip" bigint,
     "name" varchar,
     "role" role,
     "password" varchar,
     "identityCardScanImg" varchar,
     "createdAt" timestamp
);