-- Insert Majors
INSERT INTO "majors" ("major_name") VALUES
  ('Computer Science'),
  ('Electrical Engineering'),
  ('Mechanical Engineering'),
  ('Civil Engineering');

-- Insert Courses
INSERT INTO "courses" ("course_no", "course_name", "major_id" , "semester") VALUES
  ('CS101', 'Introduction to Computer Science', 1 , 1),
  ('EE101', 'Introduction to Electrical Engineering', 2 , 1),
  ('ME101', 'Introduction to Mechanical Engineering', 3,2),
  ('CE101', 'Introduction to Civil Engineering', 4,2);

-- Insert Students
INSERT INTO "students" ("id", "prefix", "first_name", "last_name", "email", "major_id") VALUES
  ('640610304', 'Mr.', 'John', 'Doe', 'john.doe@example.com', 1),
  ('640610305', 'Ms.', 'Jane', 'Smith', 'jane.smith@example.com', 2),
  ('640610306', 'Mr.', 'Mark', 'Johnson', 'mark.johnson@example.com', 3),
  ('640610307', 'Ms.', 'Emily', 'Davis', 'emily.davis@example.com', 4);

-- Insert Projects
INSERT INTO "projects" ("old_project_no", "project_no", "title_th", "title_en", "abstract", "relation_description", "academic_year", "semester", "section_id", "is_approved", "course_id", "major_id") VALUES
  (NULL, 'P002', 'โครงการศึกษา', 'Study Project', 'Abstract of Study Project', 'Project relations description', 2024, 1, 'A', true, 1, 1),
  ('P002', 'P004', 'โครงการวิศวกรรม', 'Engineering Project', 'Abstract of Engineering Project', 'Engineering project relations', 2024, 2, 'B', false, 2, 2);

-- Insert Project Students
INSERT INTO "project_students" ("project_id", "student_id") VALUES
  (1, '640610304'),
  (1, '640610305'),
  (2, '640610306'),
  (2, '640610307');
-- Insert Project Employee Types
INSERT INTO "project_employee_types" ("type_name") VALUES
  ('Manager'),
  ('Developer'),
  ('Designer');

-- Insert Roles
INSERT INTO "roles" ("role_name") VALUES
  ('Advisor'),
  ('Coordinator');

-- Insert Employees
INSERT INTO "employees" ("prefix", "first_name", "last_name", "email", "major_id", "role_id") VALUES
  ('Dr.', 'Alice', 'Walker', 'alice.walker@example.com', 1, 1),
  ('Mr.', 'Bob', 'Taylor', 'bob.taylor@example.com', 2, 2);

-- Insert Project Employees
INSERT INTO "project_employees" ("project_id", "employee_id") VALUES
  (1, 1),
  (2, 2);

-- Insert Resource Types
INSERT INTO "resource_types" ("resource_type") VALUES
  ('Textbook'),
  ('Lab Equipment');

-- Insert Resources
INSERT INTO "resources" ("title", "resource_type_id", "url") VALUES
  ('Computer Science Textbook', 1, 'https://example.com/cs-textbook'),
  ('Electrical Engineering Lab Equipment', 2, 'https://example.com/ee-lab');

-- Insert Asset Resources
INSERT INTO "asset_resources" ("resource_id", "description", "major_id") VALUES
  (1, 'A comprehensive guide to comp  uter science', 1),
  (2, 'Tools and equipment for electrical engineering labs', 2);

-- Insert Project Resources
INSERT INTO "project_resources" ("resource_id", "project_id") VALUES
  (1, 1),
  (2, 2);

-- Insert Project Resource Configs
INSERT INTO "project_resource_configs" ("title", "resource_type_id", "major_id") VALUES
  ('Default Resource Config', 1, 1),
  ('Lab Resource Config', 2, 2);

-- Insert Calendar Events
INSERT INTO "calendar" ("major_id", "event_date", "title", "description") VALUES
  (1, '2024-01-15', 'CS 101 Exam', 'Final exam for CS 101 course'),
  (2, '2024-02-20', 'EE 101 Lab', 'Electrical Engineering Lab practical session');
