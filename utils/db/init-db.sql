USE hotelmagnamentdb;

CREATE TABLE IF NOT EXISTS `locations` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `city` VARCHAR(255) NOT NULL,
    `country` VARCHAR(255) NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `zip_code` VARCHAR(255) NOT NULL,
    `latitude` DOUBLE NULL,
    `longitude` DOUBLE NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `franchises` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `url` VARCHAR(255) NOT NULL,
    `location_id` INT NULL,
    `logo_url` VARCHAR(255) NULL,
    `is_website_live` TINYINT(1) NOT NULL DEFAULT 0,
    `created_date` DATE NULL,
    `expiry_date` DATE NULL,
    `registrar_name` VARCHAR(255) NULL,
    `contact_email` VARCHAR(255) NULL,
    `protocol` VARCHAR(50) NULL,
    `is_protocol_secure` TINYINT(1) NOT NULL DEFAULT 0,
    `ssl_grade` VARCHAR(2) NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_location` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `dns_records` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `franchise_id` INT(11) NOT NULL,
  `type` VARCHAR(50) NOT NULL,
  `value` TEXT NOT NULL,
  `ttl` INT(11) DEFAULT NULL,
  `priority` INT(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`franchise_id`) REFERENCES `franchises` (`id`) ON DELETE CASCADE,
  INDEX `franchise_id_idx` (`franchise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
