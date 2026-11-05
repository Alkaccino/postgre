CREATE TABLE IF NOT EXISTS assets(
    id SERIAL PRIMARY KEY,
    asset_name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    serial_number VARCHAR(20),
    model VARCHAR(255),
    manufacturer VARCHAR(255),
    assigned_user VARCHAR(255),
    asset_location VARCHAR(255),
    purchased_price NUMERIC(12,2),
    purchased_at VARCHAR(255),
    quantity INT,
    asset_status VARCHAR(255),
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS manufacturers(
    id SERIAL PRIMARY KEY,
    manufacturer_name VARCHAR(255) NOT NULL,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories(
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    department VARCHAR(255) NOT NULL,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);