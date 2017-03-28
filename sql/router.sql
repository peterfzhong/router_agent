CREATE database db_router;
USE db_router;

CREATE  TABLE  tb_module(
    module_id   VARCHAR(256) NOT NULL  PRIMARY KEY ,
    module_name VARCHAR(256) NOT NULL,
    module_desc VARCHAR(256) NOT NULL,
    create_person VARCHAR(256) NOT NULL,
    create_time datetime NOT  NULL
);

CREATE TABLE  tb_module_machine(
  	module_id VARCHAR(256) NOT NULL,
	  ip VARCHAR(256) NOT NULL,
	  weight int NOT NULL,
	  create_user VARCHAR(128) NOT NULL,
	  create_time datetime NOT NULL,

	  PRIMARY KEY (module_id, ip)
);

CREATE TABLE  tb_module_machine_interest(
  	module_id VARCHAR(256) NOT NULL,
	  ip VARCHAR(256) NOT NULL,

	  PRIMARY KEY (module_id, ip)
);


