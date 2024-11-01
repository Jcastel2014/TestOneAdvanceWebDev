
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    image_url TEXT
)
CREATE TABLE products (
    id SERIAL PRIMARY KEY,                  
    name TEXT NOT NULL,             
    description TEXT,                        
    category VARCHAR(100),
    image_id INT,    
    average_rating DECIMAL(3, 2), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE  
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,                     
    product_id INT NOT NULL,                 
    rating INT,
    helpful_count INT DEFAULT(0),             
    comment TEXT,                               
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE 
);

ALTER TABLE products
ADD COLUMN price DECIMAL(10, 2);
