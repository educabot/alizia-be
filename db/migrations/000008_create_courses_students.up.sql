-- Courses (student groups: "2do 1era", "3ro 2da", etc.)
CREATE TABLE courses (
    id          BIGSERIAL    PRIMARY KEY,
    organization_id UUID    NOT NULL REFERENCES organizations(id),
    name        VARCHAR(255) NOT NULL,
    year        INTEGER      NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Students (enrolled in a course)
CREATE TABLE students (
    id         BIGSERIAL    PRIMARY KEY,
    course_id  BIGINT       NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Course-Subject assignments (course + subject + teacher + school year)
CREATE TABLE course_subjects (
    id              BIGSERIAL    PRIMARY KEY,
    organization_id UUID         NOT NULL REFERENCES organizations(id),
    course_id       BIGINT       NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    subject_id      BIGINT       NOT NULL REFERENCES subjects(id),
    teacher_id      BIGINT       NOT NULL REFERENCES users(id),
    school_year     INTEGER      NOT NULL,
    start_date      DATE,
    end_date        DATE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, subject_id, school_year)
);

CREATE INDEX idx_courses_organization_id ON courses(organization_id);
CREATE INDEX idx_students_course_id ON students(course_id);
CREATE INDEX idx_course_subjects_course_id ON course_subjects(course_id);
CREATE INDEX idx_course_subjects_subject_id ON course_subjects(subject_id);
CREATE INDEX idx_course_subjects_teacher_id ON course_subjects(teacher_id);
CREATE INDEX idx_course_subjects_org_year ON course_subjects(organization_id, school_year);
