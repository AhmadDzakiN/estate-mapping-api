-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 
CREATE TABLE IF NOT EXISTS estates (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    width INT NOT NULL,
    length INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trees (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    estate_id UUID NOT NULL,
    horizontal_position INT NOT NULL,
    vertical_position INT NOT NULL,
    height INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (estate_id) REFERENCES estates(id) ON DELETE CASCADE
);
