ALTER TABLE comments
DROP
CONSTRAINT comments_post_id_fkey,
ADD CONSTRAINT  comments_post_id_fkey
    FOREIGN KEY (post_id)
    REFERENCES posts(id);
