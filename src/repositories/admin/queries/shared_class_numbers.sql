-- shared_class_numbers.sql
-- Returns the weekly schedule of a course_subject as ordered rows
-- (weekly_position, is_shared) where is_shared=true when the time_slot
-- contains more than one course_subject (i.e. a shared class).
--
-- Param $1: course_subject_id
WITH subject_slots AS (
    SELECT ts.id AS slot_id,
           ROW_NUMBER() OVER (ORDER BY ts.day_of_week, ts.start_time) AS weekly_position
    FROM time_slots ts
    JOIN time_slot_subjects tss ON tss.time_slot_id = ts.id
    WHERE tss.course_subject_id = ?
)
SELECT ss.weekly_position,
       ((SELECT count(*) FROM time_slot_subjects WHERE time_slot_id = ss.slot_id) > 1) AS is_shared
FROM subject_slots ss
ORDER BY ss.weekly_position;
