CREATE INDEX tweet_user_id_idx ON tweet ((tweet.tweet->>'user_id'));
