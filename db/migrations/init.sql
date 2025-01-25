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
  "course_id" int
);

CREATE TABLE "configs" (
  "config_name" varchar,
  "value" varchar,
  "program_id" int
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
  CONSTRAINT "unique_student_combination" UNIQUE ("student_id", "semester", "academic_year", "program_id")
);

CREATE TABLE "courses" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "course_no" varchar,
  "course_name" varchar,
  "program_id" int
);

CREATE TABLE "project_roles" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "role_name" varchar,
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
  "section_id" varchar,
  "is_public" bool DEFAULT true,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp,
  "course_id" int,
  "program_id" int
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
  "prefix" varchar,
  "first_name" varchar,
  "last_name" varchar,
  "email" varchar UNIQUE,
  "program_id" int
);

CREATE TABLE "file_extensions" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "extension_name" varchar,
  "mime_type" varchar UNIQUE
);

CREATE TABLE "resource_types" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "type_name" varchar
);

CREATE TABLE "resources" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar,
  "project_resource_id" int,
  "asset_resource_id" int,
  "file_extension_id" int,
  "resource_type_id" int,
  "url" varchar,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "asset_resources" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "description" varchar,
  "program_id" int
);

CREATE TABLE "project_resources" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "project_id" int
);

CREATE TABLE "project_resource_configs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar,
  "icon_name" varchar,
  "max_file_size" int,
  "file_extension_id" int,
  "resource_type_id" int,
  "is_active" bool DEFAULT true,
  "program_id" int
);

CREATE TABLE "project_configs" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar,
  "is_active" bool DEFAULT true,
  "program_id" int
);

CREATE UNIQUE INDEX ON "project_number_counter" ("academic_year", "semester", "course_id");

CREATE UNIQUE INDEX ON "project_staffs" ("project_id", "staff_id", "project_role_id");

ALTER TABLE "configs" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_number_counter" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "courses" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "projects" ADD FOREIGN KEY ("course_id") REFERENCES "courses" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "projects" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_students" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_students" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_staffs" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_staffs" ADD FOREIGN KEY ("staff_id") REFERENCES "staffs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_staffs" ADD FOREIGN KEY ("project_role_id") REFERENCES "project_roles" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "staffs" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "students" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "resources" ADD FOREIGN KEY ("resource_type_id") REFERENCES "resource_types" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "resources" ADD FOREIGN KEY ("file_extension_id") REFERENCES "file_extensions" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "resources" ADD FOREIGN KEY ("project_resource_id") REFERENCES "project_resources" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resources" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "resources" ADD FOREIGN KEY ("asset_resource_id") REFERENCES "asset_resources" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "asset_resources" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resource_configs" ADD FOREIGN KEY ("program_id") REFERENCES "programs" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resource_configs" ADD FOREIGN KEY ("resource_type_id") REFERENCES "resource_types" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "project_resource_configs" ADD FOREIGN KEY ("file_extension_id") REFERENCES "file_extensions" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
