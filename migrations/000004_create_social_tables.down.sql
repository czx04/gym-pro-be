DROP INDEX IF EXISTS idx_comments_post_created;
DROP INDEX IF EXISTS idx_comments_user;
DROP INDEX IF EXISTS idx_comments_post;
DROP TABLE IF EXISTS comments;

DROP INDEX IF EXISTS idx_likes_user;
DROP INDEX IF EXISTS idx_likes_post;
DROP TABLE IF EXISTS likes;

DROP INDEX IF EXISTS idx_posts_user_created;
DROP INDEX IF EXISTS idx_posts_content;
DROP INDEX IF EXISTS idx_posts_created;
DROP INDEX IF EXISTS idx_posts_user;
DROP TABLE IF EXISTS posts;

DROP INDEX IF EXISTS idx_follows_following;
DROP INDEX IF EXISTS idx_follows_follower;
DROP TABLE IF EXISTS follows;
