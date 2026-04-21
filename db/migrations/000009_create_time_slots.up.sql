-- Time slots (weekly schedule grid)
CREATE TABLE time_slots (
    id          BIGSERIAL    PRIMARY KEY,
    course_id   BIGINT       NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    day_of_week SMALLINT     NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time  TIME         NOT NULL,
    end_time    TIME         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Time slot subjects (1 subject = normal class, 2 subjects = shared class)
CREATE TABLE time_slot_subjects (
    id                BIGSERIAL PRIMARY KEY,
    time_slot_id      BIGINT NOT NULL REFERENCES time_slots(id) ON DELETE CASCADE,
    course_subject_id BIGINT NOT NULL REFERENCES course_subjects(id) ON DELETE CASCADE,
    UNIQUE(time_slot_id, course_subject_id)
);

CREATE INDEX idx_time_slots_course_id ON time_slots(course_id);
CREATE INDEX idx_time_slots_course_day ON time_slots(course_id, day_of_week);
CREATE INDEX idx_tss_slot ON time_slot_subjects(time_slot_id);
CREATE INDEX idx_tss_cs ON time_slot_subjects(course_subject_id);

-- Trigger: validate course_subject belongs to same course as time_slot
CREATE OR REPLACE FUNCTION validate_time_slot_subject() RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM course_subjects cs
        JOIN time_slots ts ON ts.course_id = cs.course_id
        WHERE cs.id = NEW.course_subject_id AND ts.id = NEW.time_slot_id
    ) THEN
        RAISE EXCEPTION 'course_subject does not belong to the same course as the time_slot';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_time_slot_subject
    BEFORE INSERT OR UPDATE ON time_slot_subjects
    FOR EACH ROW
    EXECUTE FUNCTION validate_time_slot_subject();
