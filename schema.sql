CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    title TEXT,
    photo_url TEXT NOT NULL,
    comment TEXT,
    metadata JSONB DEFAULT '{}'::jsonb
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

CREATE INDEX idx_photo_tags_photo_id ON photo_tags(photo_id);
CREATE INDEX idx_photo_tags_tag_title ON photo_tags(tag_title);
CREATE INDEX idx_photos_metadata ON photos USING GIN (metadata);
CREATE INDEX idx_tags_title ON tags(title);