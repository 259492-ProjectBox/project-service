-- Insert File extensions 
INSERT INTO "file_extensions" ("extension_name","mime_type") VALUES
  ('jpeg','image/jpeg'),
  ('text','text/plain'),
  ('png','image/png'),
  ('word','application/msword'),
  ('zip','application/vnd.rar'),
  ('pdf','application/pdf'),
  ('powerpoint','application/vnd.ms-powerpoint');
  
-- Insert programs
INSERT INTO "programs" ("abbreviation", "program_name_en", "program_name_th") 
VALUES 
   ('CPE', 'Computer Engineering', 'วิศวกรรมคอมพิวเตอร์'),
   ('ISNE', 'Information Systems and Network Engineering', 'วิศวกรรมระบบสารสนเทศและเครือข่าย'),
   ('RE', 'Robotics Engineering and Artificial Intelligence', 'วิศวกรรมหุ่นยนต์และปัญญาประดิษฐ์'),
   ('EE', 'Electrical Engineering', 'วิศวกรรมไฟฟ้า'),
   ('ME', 'Mechanical Engineering', 'วิศวกรรมเครื่องกล'),
   ('CE', 'Civil Engineering', 'วิศวกรรมโยธา'),
   ('IE', 'Industrial Engineering', 'วิศวกรรมอุตสาหการ'),
   ('INE', 'Integrated Engineering', 'วิศวกรรมบูรณาการ'),
   ('ENVI', 'Environmental Engineering', 'วิศวกรรมสิ่งแวดล้อม'),
   ('MPE', 'Mining and Petroleum Engineering', 'วิศวกรรมเหมืองแร่และปิโตรเลียม');

INSERT INTO "configs" ("config_name", "value", "program_id")
VALUES 
    ('academic year', '2025', 1),
    ('academic year', '2025', 2),
    ('academic year', '2025', 3),
    ('academic year', '2025', 4),
    ('semester', '2', 1),
    ('semester', '2', 2),
    ('semester', '2', 3),
    ('semester', '2', 4);

-- Insert Courses
INSERT INTO "courses" ("course_no", "course_name", "program_id") VALUES
  ('261492', 'CPE Project', 1),
  ('CS101', 'Introduction to Computer Science', 1 ),
  ('EE101', 'Introduction to Electrical Engineering', 2 ),
  ('ME101', 'Introduction to Mechanical Engineering', 3),
  ('CE101', 'Introduction to Civil Engineering', 4);

-- Insert Students
INSERT INTO "students" 
  ("student_id", "first_name", "last_name", "email", "sec_lab" , "semester","academic_year", "course_id" ,"program_id") 
VALUES
  ('640610303', 'John', 'Doe', 'john.doe@example.com','001' , 2, 2025,1,  1),
  ('640610305', 'James', 'Brown', 'james.brown@example.com','002' , 2,  2025,2,  1),
  ('640610306', 'Jane', 'Smith', 'jane.smith@example.com','002' , 2, 2025,4,   2),
  ('640610307', 'Emily', 'Clark', 'emily.clark@example.com','003' , 1,2025, 4, 2),
  ('640610308', 'Mark', 'Johnson', 'mark.johnson@example.com','003' , 2,2025, 1,  3),
  ('640610309', 'David', 'Wilson', 'david.wilson@example.com','001' , 1,2025, 2, 3),
  ('640610310', 'Emily', 'Davis', 'emily.davis@example.com', '003' ,1,2025, 3,4),
  ('640610311', 'Sophia', 'Taylor', 'sophia.taylor@example.com','002' ,1, 2025,4,  4);

-- Insert mock data into project_roles
INSERT INTO "project_roles" ("role_name_th", "role_name_en", "program_id") VALUES 
('อาจารย์ที่ปรึกษา', 'Advisor', 1),
('กรรมการภายนอก', 'Co Advisor', 1),
('กรรมการ', 'Committee', 1),
('กรรมการภายนอก', 'External Committee', 1),
('อาจารย์ที่ปรึกษา', 'Advisor', 2),
('กรรมการภายนอก', 'Co Advisor', 2),
('กรรมการ', 'Committee', 2),
('กรรมการภายนอก', 'External Committee', 2),
('อาจารย์ที่ปรึกษา', 'Advisor', 3),
('กรรมการภายนอก', 'Co Advisor', 3),
('กรรมการ', 'Committee', 3),
('กรรมการภายนอก', 'External Committee', 3),
('อาจารย์ที่ปรึกษา', 'Advisor', 4),
('กรรมการภายนอก', 'Co Advisor', 4),
('กรรมการ', 'Committee', 4),
('กรรมการภายนอก', 'External Committee', 4),
('อาจารย์ที่ปรึกษา', 'Advisor', 5),
('กรรมการภายนอก', 'Co Advisor', 5),
('กรรมการ', 'Committee', 5),
('กรรมการภายนอก', 'External Committee', 5),
('อาจารย์ที่ปรึกษา', 'Advisor', 6),
('กรรมการภายนอก', 'Co Advisor', 6),
('กรรมการ', 'Committee', 6),
('กรรมการภายนอก', 'External Committee', 6),
('อาจารย์ที่ปรึกษา', 'Advisor', 7),
('กรรมการภายนอก', 'Co Advisor', 7),
('กรรมการ', 'Committee', 7),
('กรรมการภายนอก', 'External Committee', 7),
('อาจารย์ที่ปรึกษา', 'Advisor', 8),
('กรรมการภายนอก', 'Co Advisor', 8),
('กรรมการ', 'Committee', 8),
('กรรมการภายนอก', 'External Committee', 8),
('อาจารย์ที่ปรึกษา', 'Advisor', 9),
('กรรมการภายนอก', 'Co Advisor', 9),
('กรรมการ', 'Committee', 9),
('กรรมการภายนอก', 'External Committee', 9),
('อาจารย์ที่ปรึกษา', 'Advisor', 10),
('กรรมการภายนอก', 'Co Advisor', 10),
('กรรมการ', 'Committee', 10),
('กรรมการภายนอก', 'External Committee', 10);


-- Insert Projects
INSERT INTO "projects" ("project_no", "title_th", "title_en", "abstract_text", "academic_year", "semester","is_public" ,"section_id",  "course_id", "program_id") VALUES
  ('P002', 'โครงการศึกษา', 'Study Project', 'Abstract of Study Project', 2024, 1, true,'A',1, 1),
  ('P004', 'โครงการวิศวกรรม', 'Engineering Project', 'Abstract of Engineering Project', 2024, 2,true, 'B', 2, 2);

-- Insert Project Students
INSERT INTO "project_students" ("project_id", "student_id") VALUES
  (1, 1),
  (1, 2),
  (2, 3),
  (2, 4);

INSERT INTO "staffs" ("prefix_th", "prefix_en", "first_name_th", "last_name_th", "first_name_en", "last_name_en", "email", "program_id") VALUES 
  ('ดร.', 'Dr.', 'อลิซ', 'เทย์เลอร์', 'Alice', 'Taylor', 'alice.taylor@example.com', 1),
  ('ศร.', 'Prof.', 'บ็อบ', 'วิลเลียมส์', 'Bob', 'Williams', 'bob.williams@example.com', 2);

INSERT INTO "project_staffs" ("project_id", "staff_id", "project_role_id") VALUES 
  (1, 1, 1),
  (1, 1, 2),
  (1, 2, 2),
  (2, 2, 1),
  (2, 2, 2),
  (2, 1, 2);

INSERT INTO "resource_types" ("type_name") VALUES 
('file'),
('url');

INSERT INTO "asset_resources" ("title", "description", "program_id")
VALUES
  ('A comprehensive guide to computer science', 'A detailed guide covering various computer science topics and resources.', 1),
  ('Tools and equipment for electrical engineering labs', 'A list of essential tools and equipment required for electrical engineering labs.', 2);


-- Corrected INSERT INTO project_resources
INSERT INTO "project_resources" ("title", "icon_name", "url", "resource_type_id", "project_id")
VALUES
  ('Project Documentation', 'doc-icon', 'https://example.com/docs', 2, 1),
  ('API Reference', 'api-icon', 'https://example.com/api', 2, 1),
  ('Project GitHub Repository', 'github-icon', 'https://github.com/example/repo', 2, 2),
  ('Video Tutorial', 'video-icon', 'https://example.com/video', 2, 2);

-- Corrected INSERT INTO project_resource_configs
INSERT INTO "project_resource_configs" ("title", "icon_name", "resource_type_id", "is_active", "program_id")
VALUES
  ('Default Resource Config', 'default-icon', 1, true, 1),
  ('Lab Resource Config', 'lab-icon', 2, true, 1);


INSERT INTO "project_configs" ("title", "is_active", "program_id") VALUES
  ('title_th',true, 1),
  ('title_en', true, 1),
  ('abstract_text', true, 1),
  ('academic_year', true, 1),
  ('semester', true, 1),
  ('section_id', true, 1),
  ('course_id', true, 1);
