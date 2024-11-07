CREATE TABLE IF NOT EXISTS "roles" (
    "id" SERIAL PRIMARY KEY,
    "role_name" VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "majors" (
    "id" SERIAL PRIMARY KEY, 
    "major_name" VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "students" (
    "id" VARCHAR(255) PRIMARY KEY,
    "student_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "major_id" INTEGER NOT NULL,
    FOREIGN KEY ("major_id") REFERENCES "majors"("id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "courses" (
    "id" SERIAL PRIMARY KEY,
    "course_no" INTEGER NOT NULL,
    "course_name" VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "sections" (
    "id" SERIAL PRIMARY KEY,
    "course_id" INTEGER NOT NULL,
    "section_number" VARCHAR(255) NOT NULL,
    "semester" INTEGER NOT NULL,
    FOREIGN KEY ("course_id") REFERENCES "courses"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "employees" (
    "id" SERIAL PRIMARY KEY,
    "employee_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "role_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "projects" (
    "id" SERIAL PRIMARY KEY,
    "old_project_no" VARCHAR(255),
    "project_no" VARCHAR(255) NOT NULL,
    "title_th" VARCHAR(255),
    "title_en" VARCHAR(255),
    "abstract" TEXT,
    "relation_description" TEXT NOT NULL,
    "advisor_id" INTEGER NOT NULL,
    "course_id" INTEGER NOT NULL,
    "section_id" INTEGER,
    "semester" INTEGER NOT NULL,
    "academic_year" INTEGER NOT NULL,
    "major_id" INTEGER NOT NULL,
    "project_status" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("advisor_id") REFERENCES "employees"("id") ON DELETE SET NULL,
    FOREIGN KEY ("course_id") REFERENCES "courses"("id") ON DELETE SET NULL,
    FOREIGN KEY ("section_id") REFERENCES "sections"("id") ON DELETE SET NULL,
    FOREIGN KEY ("major_id") REFERENCES "majors"("id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "project_employees" (
    "project_id" INTEGER NOT NULL,
    "employee_id" INTEGER NOT NULL,
    FOREIGN KEY ("project_id") REFERENCES "projects"("id") ON DELETE SET NULL,
    FOREIGN KEY ("employee_id") REFERENCES "employees"("id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "project_students" (
    "project_id" INTEGER NOT NULL,
    "student_id" VARCHAR(255) NOT NULL,
    FOREIGN KEY ("project_id") REFERENCES "projects"("id") ON DELETE SET NULL,
    FOREIGN KEY ("student_id") REFERENCES "students"("id") ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS "comments" (
    "id" SERIAL PRIMARY KEY,
    "project_id" INTEGER NOT NULL,
    "comment_text" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("project_id") REFERENCES "projects"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "resource_types" (
    "id" SERIAL PRIMARY KEY,
    "resource_type" VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "resources" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(255),
    "project_id" INTEGER NOT NULL,
    "resource_type_id" INTEGER NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("project_id") REFERENCES "projects"("id") ON DELETE CASCADE,
    FOREIGN KEY ("resource_type_id") REFERENCES "resource_types"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "important_dates" (
    "id" SERIAL PRIMARY KEY,
    "major_id" INTEGER NOT NULL,
    "event_date" DATE NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    FOREIGN KEY ("major_id") REFERENCES "majors"("id") ON DELETE CASCADE
);