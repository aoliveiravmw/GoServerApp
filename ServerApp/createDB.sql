CREATE DATABASE IF NOT EXISTS plant_watering;
USE plant_watering;
CREATE TABLE IF NOT EXISTS users (
    ID int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    PRIMARY KEY (ID));
CREATE TABLE IF NOT EXISTS plants (
    plant_id int NOT NULL AUTO_INCREMENT,
    plant_name varchar(255) NOT NULL,
    last_watered int(11) NOT NULL,
    watering_interval_hours int(11) NOT NULL,
    user_id int NOT NULL,
    PRIMARY KEY (plant_id),
    FOREIGN KEY (user_id) REFERENCES users(ID));
CREATE TABLE IF NOT EXISTS pictures (
    picture_id int NOT NULL AUTO_INCREMENT,
    picture blob,
    plant_id int NOT NULL,PRIMARY KEY (picture_id),
    FOREIGN KEY (plant_id) REFERENCES plants(plant_id));