
CREATE TABLE products (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for each product
    name VARCHAR(255) NOT NULL,                -- Product name
    description TEXT,                           -- Product description
    category VARCHAR(100),                      -- Category of the product
    image_url VARCHAR(255),                     -- URL for the product image
    average_rating DECIMAL(3, 2) DEFAULT 0.00, -- Average rating of the product
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the product was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for when the product was last updated
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for each review
    product_id INT NOT NULL,                   -- Foreign key linking to the products table
    rating INT CHECK (rating >= 1 AND rating <= 5), -- Rating must be between 1 and 5
    helpful_count INT DEFAULT 0,               -- Count of how many users found the review helpful
    comment TEXT,                               -- Review comment
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the review was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Timestamp for when the review was last updated
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE -- Foreign key constraint
);
