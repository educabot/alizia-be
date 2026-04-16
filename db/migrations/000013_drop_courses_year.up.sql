-- Drop courses.year: school year now lives only in course_subjects.school_year,
-- so a course can be reused across multiple academic years via different
-- course_subjects assignments instead of duplicating the year on the course.
ALTER TABLE courses DROP COLUMN year;
