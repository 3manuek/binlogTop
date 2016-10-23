CREATE TABLE `test1` (
  `i` int(11) NOT NULL,
  PRIMARY KEY (`i`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `test2` (
  `i` int(11) NOT NULL,
  PRIMARY KEY (`i`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



TRUNCATE TABLE test1;
delete from test1 where i = 1;
select sleep(2);
delete from test1 where i = 2;
insert into test2 values(10000000);
select sleep(2);
insert into test2 values(10000002);
select sleep(2);
insert into test2 values(10000003);
select sleep(2);
delete from test2 where i = 10000002;
select sleep(2);
INSERT INTO test1 SELECT distinct i from test2;
select sleep(2);
INSERT INTO test1 SELECT distinct i from test2;

DROP TABLE test1;
DROP TABLE test2;
