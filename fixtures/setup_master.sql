CREATE USER slave@'%' IDENTIFIED BY 'mypassword';
GRANT SELECT, PROCESS, FILE, SUPER, REPLICATION CLIENT, REPLICATION SLAVE, RELOAD ON *.* TO slave@'%';
Flush Privileges; 
use mysql;
CREATE TABLE `heartbeat` (
  `ts` int(11) DEFAULT NULL
);
insert into `heartbeat` (ts) values (1);