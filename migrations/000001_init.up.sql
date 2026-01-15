-- 1. Включаем поддержку UUID, если её нет
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Создаем таблицу
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(15, 2) NOT NULL DEFAULT 0.00, -- 15 цифр всего, 2 после запятой
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Добавим индекс на имя, так как поиск по нему будет частым
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);