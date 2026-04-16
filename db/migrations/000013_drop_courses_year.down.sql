-- Restore courses.year. Existing rows get the current year as a safe default
-- since we cannot recover the original value once dropped.
ALTER TABLE courses ADD COLUMN year INTEGER NOT NULL DEFAULT EXTRACT(YEAR FROM CURRENT_DATE);
ALTER TABLE courses ALTER COLUMN year DROP DEFAULT;
