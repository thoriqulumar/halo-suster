CREATE TYPE "gender" AS ENUM (
  'male',
  'female'
);

CREATE TABLE "patient" (
   "identityNumber" bigint PRIMARY KEY,
   "phoneNumber" varchar,
   "name" varchar,
   "birthDate" date,
   "gender" gender,
   "identityCardScanImg" varchar,
   "createAt" timestamp
);
