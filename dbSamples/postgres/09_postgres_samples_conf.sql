create user script_runner password 'ILoveDatabases';

grant usage on schema olympics to script_runner;
grant select on all tables in schema olympics to script_runner;