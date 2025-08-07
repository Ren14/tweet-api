package timeline

import (
	"context"
	"fmt"
	"log"
)

const (
	// Defines a consistent key structure for user timelines in Redis.
	timelineKeyFormat = "timeline:%s"
)

// UpdateTimeline performs the "fan-out" operation. It finds all followers of the tweet's author
// and pushes the new tweet's ID onto each of their timeline lists in Redis.
// It's designed to be called asynchronously (e.g., in a goroutine).
func (s Service) UpdateTimeline(ctx context.Context, tweetAuthorID, tweetID string) {
	// 1. Get all followers from the database.
	followers, err := s.Storage.SelectFollowersByUserID(ctx, tweetAuthorID)
	if err != nil {
		log.Printf("ERROR: UpdateTimeline could not get followers for user %s: %v", tweetAuthorID, err)
		return
	}

	if len(followers) == 0 {
		log.Printf("INFO: User %s has no followers to update.", tweetAuthorID)
		return
	}

	// TODO add span with followers_count trace, for check performance

	// TODO [spike] add log if user has more than 10.000 followers.
	// avoid next logic if the user has too many followers. In this case, another approach is needed.
	// [spike] learn to celebrity user problem / pattern and how apply it.

	// 2. For each follower, push the new tweet ID to their timeline list in Redis.
	var updatedCount int
	for _, followerID := range followers {
		// TODO [spike] parallelize each follower-UpdateTimeline with go-routines. Use waitGroup to await finish results.

		// Construct the unique Redis key for this follower's timeline.
		timelineKey := fmt.Sprintf(timelineKeyFormat, followerID)

		// TODO [technical debt] firs check if timelineKey exist into Redis. If not exist, create with a TTL.
		// If exist, continue with LPush
		// NOTE: // With this implementation, each register on cache is alive for ever.

		// LPUSH adds the new tweet ID to the beginning of the list.
		if err := s.Cache.LPush(ctx, timelineKey, tweetID); err != nil {
			// Log the error but continue, so one failure doesn't stop the whole process.
			log.Printf("ERROR: Failed to push tweet %s to timeline for follower %s: %v", tweetID, followerID, err)
			continue
		}
		log.Printf("INFO: add timeline fan-out on cache for tweetID: %s followerID: %s tweetAuthorID: %s", tweetID, followerID, tweetAuthorID)
		updatedCount++

		// TODO [technical debt] evaluate the value of elements for the followerID. If is more than maxTweetsCached, trim for avoid
		// bigger cache. This method should be execute with async process.
	}

	log.Printf("INFO: Finished timeline fan-out. Successfully updated %d of %d follower timelines.", updatedCount, len(followers))
	// TODO add metric if updated is distinct to len(followers). Then see logs for troubleshooting
}
