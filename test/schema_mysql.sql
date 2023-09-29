CREATE DATABASE dbtest;
USE dbtest;

-- SELECT TABLE
CREATE TABLE table_select_test (
	user_id INT PRIMARY KEY NOT NULL,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

INSERT INTO table_select_test(user_id, username, password, email, created_on, last_login)
    VALUES (1, 'clange', 'passwd', 'clange@wgf.fr', '2020-12-04 10:42:42', '2020-12-05 10:42:42');
INSERT INTO table_select_test(user_id, username, password, email, created_on, last_login)
    VALUES (2, 'tmazzotti', 'passwd', 'tmazzotti@wgf.fr', '2020-12-04 10:42:43', NULL);
INSERT INTO table_select_test(user_id, username, password, email, created_on, last_login)
    VALUES (3, 'frichard', 'passwd', 'frichard@wgf.fr', '2020-12-04 10:42:44', '2020-12-05 10:42:44');
INSERT INTO table_select_test(user_id, username, password, email, created_on, last_login)
    VALUES (4, 'acolin', 'passwd', 'acolin@wgf.fr', '2020-12-04 10:42:45', '2020-12-05 10:42:45');

CREATE TABLE table_select_test_void (
	user_id INT PRIMARY KEY NOT NULL,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

-- INSERT TABLE
CREATE TABLE table_insert_test (
	user_id INT PRIMARY KEY NOT NULL,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

INSERT INTO table_insert_test(user_id, username, password, email, created_on, last_login)
    VALUES (1, 'clange', 'passwd', 'clange@wgf.fr', '2020-12-04 10:42:42', '2020-12-05 10:42:42');
INSERT INTO table_insert_test(user_id, username, password, email, created_on, last_login)
    VALUES (2, 'tmazzotti', 'passwd', 'tmazzotti@wgf.fr', '2020-12-04 10:42:43', NULL);
INSERT INTO table_insert_test(user_id, username, password, email, created_on, last_login)
    VALUES (3, 'frichard', 'passwd', 'frichard@wgf.fr', '2020-12-04 10:42:44', '2020-12-05 10:42:44');
INSERT INTO table_insert_test(user_id, username, password, email, created_on, last_login)
    VALUES (4, 'acolin', 'passwd', 'acolin@wgf.fr', '2020-12-04 10:42:45', '2020-12-05 10:42:45');


-- UPDATE TABLE
CREATE TABLE table_update_test (
	user_id INT PRIMARY KEY NOT NULL,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

INSERT INTO table_update_test(user_id, username, password, email, created_on, last_login)
    VALUES (1, 'clange', 'passwd', 'clange@wgf.fr', '2020-12-04 10:42:42', '2020-12-05 10:42:42');
INSERT INTO table_update_test(user_id, username, password, email, created_on, last_login)
    VALUES (2, 'tmazzotti', 'passwd', 'tmazzotti@wgf.fr', '2020-12-04 10:42:43', NULL);
INSERT INTO table_update_test(user_id, username, password, email, created_on, last_login)
    VALUES (3, 'frichard', 'passwd', 'frichard@wgf.fr', '2020-12-04 10:42:44', '2020-12-05 10:42:44');
INSERT INTO table_update_test(user_id, username, password, email, created_on, last_login)
    VALUES (4, 'acolin', 'passwd', 'acolin@wgf.fr', '2020-12-04 10:42:45', '2020-12-05 10:42:45');


-- DELETE TABLE
CREATE TABLE table_delete_test (
	user_id INT PRIMARY KEY NOT NULL,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	password VARCHAR ( 50 ) NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

INSERT INTO table_delete_test(user_id, username, password, email, created_on, last_login)
    VALUES (1, 'clange', 'passwd', 'clange@wgf.fr', '2020-12-04 10:42:42', '2020-12-05 10:42:42');
INSERT INTO table_delete_test(user_id, username, password, email, created_on, last_login)
    VALUES (2, 'tmazzotti', 'passwd', 'tmazzotti@wgf.fr', '2020-12-04 10:42:43', NULL);
INSERT INTO table_delete_test(user_id, username, password, email, created_on, last_login)
    VALUES (3, 'frichard', 'passwd', 'frichard@wgf.fr', '2020-12-04 10:42:44', '2020-12-05 10:42:44');
INSERT INTO table_delete_test(user_id, username, password, email, created_on, last_login)
    VALUES (4, 'acolin', 'passwd', 'acolin@wgf.fr', '2020-12-04 10:42:45', '2020-12-05 10:42:45');
