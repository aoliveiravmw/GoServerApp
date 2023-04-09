CREATE DATABASE IF NOT EXISTS plant_watering;
USE plant_watering;
CREATE TABLE IF NOT EXISTS users (
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    alias varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    PRIMARY KEY (alias)
);
CREATE TABLE IF NOT EXISTS plants (
  plant_name varchar(255) NOT NULL,
  plant_alias varchar(255) NOT NULL UNIQUE,
  last_watered int(11) NOT NULL,
  watering_interval_hours int(11) NOT NULL,
  user_alias varchar(255) NOT NULL,
  PRIMARY KEY (plant_alias),
  FOREIGN KEY (user_alias) REFERENCES users(alias)
);
CREATE TABLE IF NOT EXISTS pictures (
    picture_id int NOT NULL AUTO_INCREMENT,
    picture blob,
    plant_alias varchar(255) NOT NULL UNIQUE,
    PRIMARY KEY (picture_id),
    FOREIGN KEY (plant_alias) REFERENCES plants(plant_alias));