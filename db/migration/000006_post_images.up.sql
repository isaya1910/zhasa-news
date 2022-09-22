CREATE TABLE post_images(
                            id SERIAL PRIMARY KEY,
                            image_url TEXT UNIQUE NOT NULL,
                            post_id INT REFERENCES posts(id) ON DELETE CASCADE NOT NULL
);