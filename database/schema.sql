DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    `pvt_id` SERIAL PRIMARY KEY,
    `user_id` UUID UNIQUE NOT NULL,
    `username` VARCHAR(50) UNIQUE NOT NULL,
    `display_name` VARCHAR(150) NOT NULL,
    `password` BYTEA NOT NULL,
    `password_salt` BYTEA NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `last_logged_in` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_modified_column()   
RETURNS TRIGGER AS $$
BEGIN
    IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.`created_at` := OLD.`created_at`;
        NEW.`updated_at` = now();
        RETURN NEW;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_customer_modtime
BEFORE UPDATE ON customer
FOR EACH ROW
EXECUTE PROCEDURE update_modified_column();