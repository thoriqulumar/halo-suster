CREATE TABLE "medicalRecord" (
     "identityNumber" bigint,
     "symptoms" text,
     "medications" text,
     "createdAt" timestamp,
     "createdBy" uuid
);

ALTER TABLE "medicalRecord" ADD FOREIGN KEY ("createdBy") REFERENCES "staff" ("id");

ALTER TABLE "medicalRecord" ADD FOREIGN KEY ("identityNumber") REFERENCES "patient" ("identityNumber");
