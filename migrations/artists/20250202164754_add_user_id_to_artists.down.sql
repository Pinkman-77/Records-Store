ALTER TABLE artists DROP CONSTRAINT IF EXISTS fk_artists_user;
ALTER TABLE artists DROP COLUMN IF EXISTS user_id;
