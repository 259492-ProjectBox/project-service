CREATE TABLE "programs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "program_name_th" varchar UNIQUE,
  "program_name_en" varchar UNIQUE,
  "abbreviation" varchar UNIQUE
);

CREATE TABLE "project_number_counter" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "number" int,
  "academic_year" int,
  "semester" int,
  "program_id" int
);

CREATE TABLE "configs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "config_name" varchar,
  "value" varchar,
  "program_id" int NULL,
  CONSTRAINT "unique_config_name_program_id" UNIQUE ("config_name", "program_id")
);

CREATE TABLE "students" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "student_id" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "sec_lab" varchar,
  "email" varchar,
  "semester" int,
  "academic_year" int,
  "program_id" int,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT "unique_student_combination" UNIQUE ("student_id", "semester", "academic_year", "program_id")
);

CREATE TABLE "project_roles" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "role_name_th" varchar,
  "role_name_en" varchar,
  "program_id" int
);

CREATE TABLE "projects" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "project_no" varchar UNIQUE,
  "title_th" varchar,
  "title_en" varchar,
  "abstract_text" text,
  "academic_year" int,
  "semester" int,
  "is_public" bool DEFAULT false,
  "section_id" varchar,
  "is_public" bool DEFAULT true,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp,
  "program_id" int
);

CREATE TABLE "file_extensions" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "extension_name" varchar,
  "mime_type" varchar UNIQUE
);

CREATE TABLE "project_students" (
  "project_id" int,
  "student_id" int
);

CREATE TABLE "project_staffs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "project_id" int,
  "staff_id" int,
  "project_role_id" int
);

CREATE TABLE "staffs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "prefix_th" varchar,
  "prefix_en" varchar,
  "first_name_th" varchar,
  "last_name_th" varchar,
  "first_name_en" varchar,
  "last_name_en" varchar,
  "email" varchar,
  "is_active" bool DEFAULT true,
  "program_id" int,
  CONSTRAINT "unique_email_program" UNIQUE ("email",  "program_id")
);

CREATE TABLE "resource_types" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "type_name" varchar
);

CREATE TABLE "asset_resources" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar,  
  "description" varchar,
  "program_id" int
);

CREATE TABLE "project_resources" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar,
  "icon_name" varchar,
  "url" varchar,
  "resource_type_id" int,
  "project_id" int,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "project_resource_configs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar ,
  "icon_name" varchar,
  "resource_type_id" int,
  "is_active" bool DEFAULT true,
  "program_id" int,
  CONSTRAINT "unique_title_programId" UNIQUE ("title",  "program_id")
);

CREATE TABLE "project_configs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar,
  "is_active" bool DEFAULT true,
  "program_id" int
);

CREATE TABLE "project_keywords" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "keyword_id" varchar,
  "project_id" int
);

CREATE TABLE "keywords" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "keyword" varchar,
  "program_id" int
);

CREATE UNIQUE INDEX ON "keywords" ("keyword", "program_id");

CREATE UNIQUE INDEX ON "students" ("student_id", "semester", "academic_year", "program_id");

CREATE UNIQUE INDEX ON "configs" ("config_name", "program_id");

CREATE UNIQUE INDEX ON "project_number_counter" ("academic_year", "semester", "program_id");

CREATE UNIQUE INDEX ON "project_staffs" ("project_id", "staff_id", "project_role_id");

ALTER TABLE "project_keywords" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ;

ALTER TABLE "project_keywords" ADD FOREIGN KEY ("keyword_id") REFERENCES "keywords" ("id") ;

ALTER TABLE "keywords" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "configs" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_number_counter" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "projects" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_students" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_students" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_staffs" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_staffs" ADD FOREIGN KEY ("staff_id") REFERENCES "staffs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_staffs" ADD FOREIGN KEY ("project_role_id") REFERENCES "project_roles" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "staffs" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "students" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resources" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "asset_resources" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resource_configs" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resource_configs" ADD FOREIGN KEY ("resource_type_id") REFERENCES "resource_types" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
