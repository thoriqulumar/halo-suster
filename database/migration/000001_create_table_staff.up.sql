CREATE TYPE "role" AS ENUM (
  'it',
  'nurse'
);

CREATE TYPE "status" AS ENUM (
    'active',
    'deleted'
);


CREATE TABLE "staff" (
     "id" uuid PRIMARY KEY,
     "nip" bigint,
     "name" varchar,
     "role" role,
     "password" varchar,
     "identityCardScanImg" varchar,
     "status" status,
     "createdAt" timestamp
);