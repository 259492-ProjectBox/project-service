-- Insert File extensions 
INSERT INTO "file_extensions" ("extension_name","mime_type") VALUES
  ('jpeg','image/jpeg'),
  ('text','text/plain'),
  ('png','image/png'),
  ('word','application/msword'),
  ('zip','application/vnd.rar'),
  ('pdf','application/pdf'),
  ('powerpoint','application/vnd.ms-powerpoint'),
  ('excel','application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'),
  ('autocad','application/acad');
  
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
    ('academic year', '2568', 1),
    ('academic year', '2568', 2),
    ('academic year', '2568', 3),
    ('academic year', '2568', 4),
    ('academic year', '2568', 5),
    ('academic year', '2568', 6),
    ('academic year', '2568', 7),
    ('academic year', '2568', 8),
    ('academic year', '2568', 9),
    ('academic year', '2568', 10),
    ('semester', '2', 1),
    ('semester', '2', 2),
    ('semester', '2', 3),
    ('semester', '2', 4),
    ('semester', '2', 5),
    ('semester', '2', 6),
    ('semester', '2', 7),
    ('semester', '2', 8),
    ('semester', '2', 9),
    ('semester', '2', 10),
    ('highest_academic_year', '2568', NULL),
    ('lowest_academic_year', '2568', NULL);

-- Insert Courses
INSERT INTO "courses" ("course_no", "course_name", "program_id") VALUES
  ('261492', 'CPE Project', 1);
  
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

INSERT INTO "resource_types" ("type_name") VALUES 
('file'),
('url');

INSERT INTO "project_configs" ("title", "is_active", "program_id") VALUES
  ('title_th',true, 1),
  ('title_en', true, 1),
  ('abstract_text', true, 1),
  ('academic_year', true, 1),
  ('semester', true, 1),
  ('section_id', true, 1),
  ('course_id', true, 1);
