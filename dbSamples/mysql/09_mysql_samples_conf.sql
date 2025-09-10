create user 'script_runner'@'%' identified by 'ILoveDatabases';
flush privileges;

grant select on `olympics`.* to 'script_runner'@'%';
flush privileges;