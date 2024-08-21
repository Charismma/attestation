DROP TABLE IF EXISTS comments_post;

CREATE TABLE comments_post(
    id SERIAL PRIMARY KEY,
	post_id BIGINT NOT NULL,
	parent_id BIGINT DEFAULT NULL,
	content TEXT NOT NULL,
	addTime BIGINT NOT NULL,
	visible bool,
	FOREIGN KEY (parent_id) REFERENCES comments_post(id)
);