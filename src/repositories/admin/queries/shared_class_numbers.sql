-- shared_class_numbers.sql
-- Returns the weekly schedule of a course_subject as ordered rows
-- (weekly_position, is_shared) where is_shared=true when the time_slot
-- contains more than one course_subject (i.e. a shared class).
--
-- The organization_id filter is intentional defense-in-depth: the usecase
-- already resolves the course_subject scoped to the tenant, but joining on
-- course_subjects.organization_id here means a future refactor that drops the
-- pre-check still cannot leak another tenant's schedule.
--
-- Param $1: course_subject_id
-- Param $2: organization_id
WITH subject_slots AS (
    SELECT ts.id AS slot_id,
           ROW_NUMBER() OVER (ORDER BY ts.day_of_week, ts.start_time) AS weekly_position
    FROM time_slots ts
    JOIN time_slot_subjects tss ON tss.time_slot_id = ts.id
    JOIN course_subjects cs ON cs.id = tss.course_subject_id
    WHERE tss.course_subject_id = ?
      AND cs.organization_id = ?
)
SELECT ss.weekly_position,
       ((SELECT count(*) FROM time_slot_subjects WHERE time_slot_id = ss.slot_id) > 1) AS is_shared
FROM subject_slots ss
ORDER BY ss.weekly_position;
