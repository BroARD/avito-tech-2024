CREATE TABLE IF NOT EXISTS houses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    developer VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    last_flat_added TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS flats (
    id SERIAL PRIMARY KEY,
    house_id INT NOT NULL,
    number INT NOT NULL,
    price INT NOT NULL,
    rooms INT NOT NULL,
    status VARCHAR(50) DEFAULT 'created' NOT NULL,

    CONSTRAINT fk_flat_house FOREIGN KEY (house_id) REFERENCES houses(id) ON DELETE CASCADE,

    CONSTRAINT unique_house_flat_number UNIQUE (house_id, number)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- 1. Функция, которая обновит дату у дома
CREATE OR REPLACE FUNCTION update_house_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE houses SET last_flat_added = NOW() WHERE id = NEW.house_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Вешаем триггер на таблицу flats
CREATE TRIGGER tr_update_house_timestamp
AFTER INSERT ON flats
FOR EACH ROW EXECUTE FUNCTION update_house_timestamp();

-- 3. Поле индексов для поиска всех квартир со статусом
CREATE INDEX IF NOT EXISTS idx_flats_user_search 
ON flats(house_id, status) 
INCLUDE (number, price, rooms);

