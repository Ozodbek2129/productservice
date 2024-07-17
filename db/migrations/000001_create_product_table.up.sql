CREATE TABLE if not exists product_categories (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT not NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE if not exists products (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category_id UUID REFERENCES product_categories(id),
    artisan_id UUID not null,--REFERENCES users(id)
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE if not exists orders (
    id UUID PRIMARY KEY,
    user_id UUID not null,--REFERENCES users(id)
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    shipping_address JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE if not EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    product_id UUID REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE if not exists ratings (
    id UUID PRIMARY KEY,
    product_id UUID REFERENCES products(id),
    user_id UUID,--REFERENCES users(id)
    rating DECIMAL(2, 1) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE if not EXISTS payments (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    card_number VARCHAR(50) not NULL,
    card_expiry VARCHAR(7) NOT NULL,
    transaction_id VARCHAR(100),
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);