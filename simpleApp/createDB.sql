CREATE DATABASE IF NOT EXISTS simpleapp;
USE simpleapp; 
CREATE TABLE IF NOT EXISTS pictures (
    picture_id int NOT NULL AUTO_INCREMENT,
    picture MEDIUMBLOB,
    PRIMARY KEY (picture_id));