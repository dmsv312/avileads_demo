BEGIN;

ALTER TABLE audit_logs
    ADD COLUMN rk_id INTEGER NULL;

ALTER TABLE audit_logs
    ADD CONSTRAINT fk_audit_logs_rk
        FOREIGN KEY (rk_id)
            REFERENCES vacancy_avito (id)
            ON DELETE SET NULL
            ON UPDATE CASCADE;

CREATE INDEX idx_audit_logs_rk_id ON audit_logs (rk_id);

COMMIT;