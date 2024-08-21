DROP TABLE IF EXISTS comments_post;

CREATE TABLE comments_post(
    id SERIAL PRIMARY KEY,
	post_id BIGINT NOT NULL,
	parent_id BIGINT DEFAULT 0,
	content TEXT NOT NULL,
	addTime BIGINT NOT NULL,
	);