CREATE TABLE IF NOT EXISTS ACTIONS_CONFIGURATIONS
(
    ID               VARCHAR(36) PRIMARY KEY,
    METRIC_ACTION_ID VARCHAR(36),
    REPEATABLE       BOOLEAN                             NOT NULL,
    NUMBER_OF_CYCLES SMALLINT,
    CREATED_AT       TIMESTAMP DEFAULT clock_timestamp() NOT NULL,
    DELETED_AT       TIMESTAMP DEFAULT clock_timestamp() NOT NULL,
    CONSTRAINT FK_METRIC_METRIC FOREIGN KEY (METRIC_ACTION_ID) REFERENCES METRICS_GROUP_ACTIONS (ID)
)
