CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    title TEXT,
    photo_url TEXT NOT NULL,
    comment TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE TABLE posts (
    slug TEXT PRIMARY KEY,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tags (
    title TEXT PRIMARY KEY,
    comment TEXT
);

CREATE TABLE photo_tags (
    photo_id INT NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    tag_title TEXT NOT NULL REFERENCES tags(title) ON DELETE CASCADE,
    PRIMARY KEY (photo_id, tag_title)
);

CREATE TABLE post_tags (
    post_slug TEXT NOT NULL REFERENCES posts(slug) ON DELETE CASCADE,
    tag_title TEXT NOT NULL REFERENCES tags(title) ON DELETE CASCADE,
    PRIMARY KEY (post_slug, tag_title)
);


CREATE INDEX idx_photo_tags_photo_id ON photo_tags(photo_id);
CREATE INDEX idx_photo_tags_tag_title ON photo_tags(tag_title);
CREATE INDEX idx_photos_metadata ON photos USING GIN (metadata);
CREATE INDEX idx_tags_title ON tags(title);
CREATE INDEX idx_post_tags_post_slug ON post_tags(post_slug);
CREATE INDEX idx_post_tags_tag_title ON post_tags(tag_title);
CREATE INDEX idx_posts_published ON posts(published);
CREATE INDEX idx_posts_slug ON posts(slug);