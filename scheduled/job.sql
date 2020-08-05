-- CREATE TABLE pet (name VARCHAR(20), owner VARCHAR(20), species VARCHAR(20), sex CHAR(1), birth DATE, death DATE);


CREATE TABLE IF NOT EXISTS agendaJob2 (
    name VARCHAR(40) PRIMARY KEY, 
    nextRun DATETIME, 
    lastRun DATETIME,
    scheduled BOOLEAN,
    jobRunning BOOLEAN,
    lastErr VARCHAR(255)
);

INSERT INTO agendaJob (name, nextRun, locked, lastRun) VALUES ('testJob','2020-07-25 09:00:00', FALSE, '2020-07-26 09:00:00');
