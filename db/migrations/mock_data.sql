-- Insert Majors
INSERT INTO "majors" ("major_name") VALUES
  ('Computer Science'),
  ('Electrical Engineering'),
  ('Mechanical Engineering'),
  ('Civil Engineering');

INSERT INTO "project_status" ("status_name", "major_id") VALUES 
('Pending', 1),
('Approved', 2),
('Rejected', 3),
('In Progress', 1),
('Completed', 4);

INSERT INTO "configs" ("config_name", "value", "major_id")
VALUES 
    ('academic year', '2025', 1),
    ('academic year', '2025', 2),
    ('academic year', '2025', 3),
    ('academic year', '2025', 4);

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
INSERT INTO "projects" ("project_no", "title_th", "title_en", "abstract_text", "academic_year", "semester", "section_id", "status_id", "course_id", "major_id") VALUES
  ('P002', 'โครงการศึกษา', 'Study Project', 'Abstract of Study Project', 2024, 1, 'A', 1, 1, 1),
  ('P004', 'โครงการวิศวกรรม', 'Engineering Project', 'Abstract of Engineering Project', 2024, 2, 'B', 2, 2, 2);

-- Insert Project Students
INSERT INTO "project_students" ("project_id", "student_id") VALUES
  (1, '640610304'),
  (1, '640610305'),
  (2, '640610306'),
  (2, '640610307');


-- Insert Employees
INSERT INTO "employees" ("prefix", "first_name", "last_name", "email", "major_id") VALUES
  ('Dr.', 'Alice', 'Walker', 'alice.walker@example.com', 1),
  ('Mr.', 'Bob', 'Taylor', 'bob.taylor@example.com', 2);

-- Insert Project Employees
INSERT INTO "project_employees" ("project_id", "employee_id") VALUES
  (1, 1),
  (2, 2);

-- Insert Resource Types
INSERT INTO "resource_types" ("type_name","mime_type") VALUES
  ('jpeg','image/jpeg'),
  ('text','text/plain'),
  ('png','image/png'),
  ('word','application/msword'),
  ('zip','application/vnd.rar'),
  ('pdf','application/pdf'),
  ('powerpoint','application/vnd.ms-powerpoint');

  -- Insert Asset Resources
  INSERT INTO "asset_resources" ("description", "major_id") VALUES
    ('A comprehensive guide to comp  uter science', 1),
    ('Tools and equipment for electrical engineering labs', 2);

  -- Insert Project Resources
  INSERT INTO "project_resources" ("project_id") VALUES
    (1),
    (2);

  -- Insert Resources
  INSERT INTO "resources" ("title", "resource_type_id", "url" , "asset_resource_id" , "project_resource_id") VALUES
    ('Computer Science Textbook', 1, 'https://example.com/cs-textbook',1 , NULL),
    ('Electrical Engineering Lab Equipment', 2, 'https://example.com/ee-lab', NULL, 2);

-- Insert Project Resource Configs
INSERT INTO "project_resource_configs" ("title", "resource_type_id", "major_id") VALUES
  ('Default Resource Config', 1, 1),
  ('Lab Resource Config', 2, 2);

-- Insert Project Resource Configs
INSERT INTO "project_configs" ("title", "is_active", "major_id") VALUES
  ('title_th',true, 1),
  ('title_en', true, 2),
  ('abstract_text', true, 2),
  ('academic_year', true, 2),
  ('semester', true, 2),
  ('section_id', true, 2),
  ('course_id', true, 2);

-- Insert Calendar Events
INSERT INTO "calendar" ("major_id", "start_date", "end_date", "description" , "title") VALUES
  (1, '2024-01-15', '2024-01-15', 'Final exam for CS 101 course' , 'CS 101 Final Exam'),
  (2, '2024-02-20', '2024-02-20', 'Electrical Engineering Lab practical session' , 'EE Lab Practical'),
  (3, '2024-03-25', '2024-03-25', 'Mechanical Engineering project presentation' , 'ME Project Presentation'),
  (4, '2024-04-30', '2024-04-30', 'Civil Engineering project submission' , 'CE Project Submission');
