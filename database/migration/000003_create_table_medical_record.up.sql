CREATE TABLE "medicalRecord" (
     "identityNumber" bigint,
     "symptoms" text,
     "medications" text,
     "createAt" timestamp,
     "createBy" uuid
);

ALTER TABLE "medicalRecord" ADD FOREIGN KEY ("createBy") REFERENCES "staff" ("id");

ALTER TABLE "medicalRecord" ADD FOREIGN KEY ("identityNumber") REFERENCES "patient" ("identityNumber");
