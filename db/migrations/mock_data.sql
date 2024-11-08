-- Insert roles
INSERT INTO "roles" (role_name) VALUES ('Advisor'), ('Committee');

-- Insert majors
INSERT INTO "majors" (major_name) VALUES ('Computer Science'), ('Information Technology'), ('Software Engineering');

-- Insert students
INSERT INTO "students" (id, prefix, first_name, last_name, email, major_id) VALUES
('640610304', 'Ms.', 'Alice', 'Smith', 'alice@example.com', 1),
('640610305', 'Mr.', 'Bob', 'Johnson', 'bob@example.com', 1),
('640610306', 'Mr.', 'Charlie', 'Brown', 'charlie@example.com', 2),
('640610307', 'Mr.', 'David', 'Wilson', 'david@example.com', 3);

-- Insert employees (with proper major_id)
INSERT INTO "employees" (prefix, first_name, last_name, email, role_id) VALUES
('Dr.', 'Jane', 'Doe', 'jane.doe@example.com', 1),
('Dr.', 'John', 'Smith', 'john.smith@example.com', 2);

-- Insert courses (make sure major_id is removed if not in the table structure)
INSERT INTO "courses" (course_no, course_name) VALUES
('259191', 'Introduction to Computer Science'),
('259192', 'Data Structures'),
('259193', 'Software Engineering I'),
('259194', 'Database Systems');

-- Insert sections
INSERT INTO "sections" (course_id, section_number, semester) VALUES
(1, 'A', 1),
(1, 'B', 2),
(2, 'A', 1),
(3, 'C', 1);

-- Insert projects
INSERT INTO "projects" 
(old_project_no, project_no, title_th, title_en, abstract, relation_description, advisor_id, course_id, section_id, semester, academic_year, major_id, project_status, created_at) 
VALUES
(NULL, 'P001-1/67', 'โปรเจค 1', 'Project 1', 'This is an abstract for project 1.', 'Relation description 1', 1, 1, 1, 1, 2567, 1, 'อยู่ในการพิจารณา', CURRENT_TIMESTAMP),
('P001-1/67', 'P002-1/67', 'โปรเจค 2', 'Project 2', 'This is an abstract for project 2.', 'Relation description 2', 2, 2, 2, 1, 2567, 1, 'อยู่ในการพิจารณา', CURRENT_TIMESTAMP);

-- Insert project students
INSERT INTO "project_students" (project_id, student_id) VALUES
(1, '640610304'),
(1, '640610305'),
(2, '640610306');

-- Insert project employees
INSERT INTO "project_employees" (project_id, employee_id) VALUES (1, 1), (2, 2);

-- Insert comments
INSERT INTO "comments" (project_id, comment_text, created_at) VALUES
(1, 'Great project! Looking forward to seeing the final results.', CURRENT_TIMESTAMP),
(2, 'Interesting approach to the problem.', CURRENT_TIMESTAMP);

-- Insert resource types
INSERT INTO "resource_types" (resource_type) VALUES
('Report'),
('Presentation'),
('Video');

-- Insert resources
INSERT INTO "resources" (title, project_id, resource_type_id, url, created_at) VALUES
('Project Report', 1, 1, 'http://example.com/report1', CURRENT_TIMESTAMP),
('Project Presentation', 1, 2, 'http://example.com/presentation1', CURRENT_TIMESTAMP),
('Project Video', 2, 3, 'http://example.com/video1', CURRENT_TIMESTAMP);

-- Insert important dates
INSERT INTO "important_dates" (major_id, event_date, title, description) VALUES
(1, '2024-12-01', 'Submission Deadline', 'Final project submission deadline for CS majors.'),
(2, '2024-11-15', 'Midterm Exam', 'Midterm exam for IT majors.');