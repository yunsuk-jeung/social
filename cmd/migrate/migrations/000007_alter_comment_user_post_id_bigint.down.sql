ALTER TABLE comments
ALTER COLUMN post_id SET DEFAULT nextval('comments_post_id_seq'),
ALTER COLUMN user_id SET DEFAULT nextval('comments_user_id_seq');
