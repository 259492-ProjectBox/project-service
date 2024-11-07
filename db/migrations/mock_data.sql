INSERT INTO "roles" (id, role_name) VALUES
(1 ,'Advisor'),
(2, 'Committee');

INSERT INTO "majors" (id ,major_name) VALUES
(1,'Computer Science'),
(2,'Information Technology'),
(3,'Software Engineering');

SELECT * FROM "majors";
-- Insert students (note that student_id is now the primary key)
INSERT INTO "students" (student_id, student_name, email, major_id) VALUES
('640610304', 'Alice Smith', 'alice@example.com', 1),
('640610305', 'Bob Johnson', 'bob@example.com', 1),
('640610306', 'Charlie Brown', 'charlie@example.com', 2),
('640610307', 'David Wilson', 'david@example.com', 3);


-- Insert courses
INSERT INTO "courses" (id ,course_no, course_name) VALUES
(1,101, 'Introduction to Computer Science'),
(2, 102, 'Data Structures'),
(3,201, 'Software Engineering I'),
(4, 202, 'Database Systems');

-- Insert sections
INSERT INTO "sections" (id,course_id, section_number, semester) VALUES
(1,1, 'A', 1),
(2,1, 'B', 2),
(3,2, 'A', 1),
(4,3, 'C', 1);

-- Insert employees
INSERT INTO "employees" (id,employee_name, email, role_id, created_at) VALUES
(1,'Dr. Jane Doe', 'jane.doe@example.com', 1, CURRENT_TIMESTAMP),
(2,'Dr. John Smith', 'john.smith@example.com', 2, CURRENT_TIMESTAMP);

-- Insert projects
INSERT INTO "projects" 
(id ,old_project_no, project_no, title_th, title_en, abstract, relation_description, advisor_id, course_id, section_id, semester, academic_year, major_id, project_status, created_at) 
VALUES
(1,NULL, 'P001-1/67', 'โปรเจค 1', 'Project 1', 'This is an abstract for project 1.', 'Relation description 1', 1, 1, 1, 1, 2567, 1, 'อยู่ในการพิจารณา', CURRENT_TIMESTAMP),
(2,'P001-1/67', 'P002-1/67', 'โปรเจค 2', 'Project 2', 'This is an abstract for project 2.', 'Relation description 2', 2, 2, 2, 1, 2567, 1, 'อยู่ในการพิจารณา', CURRENT_TIMESTAMP);

-- Insert project students
INSERT INTO "project_students" (project_id, student_id) VALUES
(1, '640610304'),
(1, '640610305'),
(2, '640610306');

-- Insert project employees
INSERT INTO "project_employees" (project_id, employee_id , role_id) VALUES
(1, 1 ,1),
(2, 2 , 1);

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
