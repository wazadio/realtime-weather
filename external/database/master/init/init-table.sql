CREATE TABLE weather(
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255),
    name VARCHAR(255), 
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    icon VARCHAR(255),
    icon_num INT,
    summary TEXT,
    temperature REAL,
    wind_speed REAL,
    wind_angel INT,
    wind_dir VARCHAR(255),
    precipitation_total REAL,
    precipitation_type VARCHAR(255),
    cloud_cover INT
);

GRANT SELECT ON ALL TABLES IN SCHEMA public TO replica_user;