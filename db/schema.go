package db

const USERSCHEMA = `
CREATE TABLE Users (
  id int,
  name varchar(80) UNIQUE PRIMARY KEY,
  hash varchar(255) NOT NULL,
  salt varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  type int NOT NULL
);

CREATE TABLE User_Types (
  code int PRIMARY KEY,
  name varchar(100) NOT NULL
);

INSERT INTO User_Types (code, name) VALUES (1, 'admin'),(2, 'employee'),(3, 'public');

ALTER TABLE Users ADD FOREIGN KEY (type) REFERENCES User_Types (code);
`

const TCSCHEMA = `
CREATE TABLE TC_Components (
  id int UNIQUE PRIMARY KEY AUTO_INCREMENT,
  name varchar(100) NOT NULL UNIQUE,
  text varchar(1000) NOT NULL
);

CREATE TABLE TC_Targeting (
  id int UNIQUE PRIMARY KEY AUTO_INCREMENT,
  name varchar(100) NOT NULL UNIQUE,
  text varchar(500) NOT NULL,
  cost int NOT NULL,
  size int NOT NULL
);

CREATE TABLE TC_Triggers (
  id int UNIQUE PRIMARY KEY AUTO_INCREMENT,
  name varchar(100) NOT NULL UNIQUE,
  text varchar(500) NOT NULL,
  cost int NOT NULL,
  area int NOT NULL
);

CREATE TABLE TC_Tiers (
  id int UNIQUE PRIMARY KEY AUTO_INCREMENT,
  tiers JSON NOT NULL
);

CREATE TABLE TC_Component_Tiers (
  component_id int NOT NULL,
  tier_id int NOT NULL,
  name varchar(100) NOT NULL
);

CREATE TABLE TC_Upgrades (
  id int UNIQUE PRIMARY KEY AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  text varchar(500) NOT NULL,
  cost int NOT NULL,
  type int NOT NULL,
  component_id int NOT NULL,
  max int DEFAULT 1
);

CREATE TABLE TC_Upgrade_Types (
  code int PRIMARY KEY,
  name varchar(100) NOT NULL
);

INSERT INTO TC_Upgrade_Types (code, name) VALUES (1, 'universal'),(2, 'trigger'),(3, 'target'),(4, 'component');

ALTER TABLE TC_Upgrades ADD FOREIGN KEY (type) REFERENCES TC_Upgrade_Types (code);

ALTER TABLE TC_Component_Tiers ADD FOREIGN KEY (component_id) REFERENCES TC_Components (id);

ALTER TABLE TC_Component_Tiers ADD FOREIGN KEY (tier_id) REFERENCES TC_Tiers (id);
`
