CREATE DATABASE mydb;
USE mydb;
CREATE TABLE `mytable` ( id INT PRIMARY KEY, name VARCHAR(20) );
INSERT INTO `mytable` VALUES ( 1, 'Will' );
INSERT INTO mytable VALUES ( 2, 'Marry' );
INSERT INTO mytable VALUES ( 3, 'Dean' );
SELECT id, name FROM mytable WHERE id = 1;
UPDATE mytable SET name = 'Willy' WHERE id = 1;
SELECT id, name FROM mytable;
DELETE FROM mytable WHERE id = 1;
SELECT id, name FROM mytable;
DROP DATABASE mydb;
SELECT COUNT(1) FROM mytable; gives the NUMBER OF records IN the TABLE



SAVEPOINT `1`;
RELEASE SAVEPOINT "1";
ROLLBACK TO '1';
