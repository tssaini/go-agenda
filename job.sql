-- CREATE TABLE pet (name VARCHAR(20), owner VARCHAR(20), species VARCHAR(20), sex CHAR(1), birth DATE, death DATE);
CREATE TABLE agendaJob (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(40) NOT NULL, 
    nextRun DATETIME,
    locked BOOLEAN, 
    lastRun DATETIME
);

INSERT INTO agendaJob (name, nextRun, locked, lastRun) VALUES ('testJob','2020-07-25 09:00:00', FALSE, '2020-07-26 09:00:00');