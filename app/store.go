package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/thanhpk/randstr"
)

var decks Decks

const deckSize = 52
const deckIdLen = 10
const lockIdLen = 6
const maxLockRetry = 3

func DrawFrom(deckId string) (uint64, Card, error) {
	// check if decks exists
	exists, err := decks.has(deckId)
	if err != nil {
		return 0, Card{}, err
	}
	if !exists {
		return 0, Card{}, fmt.Errorf(fmt.Sprintf(DECK_DONT_EXIST, deckId))
	}
	// acquiring lock
	lockId, ok := acquireLock(deckId)
	if !ok {
		msg := fmt.Sprintf(MAX_RETRY_LOCK_ACQ, deckId)
		Logger.Error(msg)
		return 0, Card{}, errors.New(msg)
	}
	// release lock
	defer releaseLock(deckId, lockId)

	// drawing card
	deck, err := decks.get(deckId)
	if err != nil {
		Logger.Error(err)
		return 0, Card{}, err
	}
	newDeck, card, err := DrawCard(deck)
	if err != nil {
		return 0, Card{}, err
	}
	if newDeck == 0 {
		err = decks.del(deckId)
	} else {
		err = decks.set(deckId, newDeck)
	}
	return deck, card, err
}

func NewDeck() (string, error) {
	deckId := randstr.String(deckIdLen)
	for exists, err := decks.has(deckId); exists; {
		if err != nil {
			return "", err
		}
		fmt.Printf("deck id: %s already exists", deckId)
		deckId = randstr.String(deckIdLen)
	}

	err := decks.set(deckId, 1<<deckSize-1)
	return deckId, err
}

var locksClient *redis.Client

type Decks struct {
	store *redis.Client
}

func (d *Decks) get(deckId string) (uint64, error) {
	return d.store.Get(context.Background(), deckId).Uint64()
}

func (d *Decks) set(deckId string, deck uint64) error {
	return d.store.Set(context.Background(), deckId, deck, time.Minute*10).Err()
}

func (d *Decks) del(deckId string) error {
	return d.store.Del(context.Background(), deckId).Err()
}

func (d *Decks) has(deckId string) (bool, error) {
	exists, err := d.store.Exists(context.Background(), deckId).Result()
	return exists == 1, err
}

func init() {
	locksClient = redis.NewClient(&redis.Options{
		Addr:     "redis-locks-service:6379", // Replace with your Redis server address and port
		Password: "",                         // Set if Redis requires authentication
		DB:       0,                          // Specify the Redis database
	})
	decks.store = redis.NewClient(&redis.Options{
		Addr:     "redis-decks-service:6379", // Replace with your Redis server address and port
		Password: "",                         // Set if Redis requires authentication
		DB:       0,                          // Specify the Redis database
	})
}

// Acquire and Release lock
func acquireLock(deckId string) (string, bool) {
	lockId := randstr.String(lockIdLen)
	attempt := 0

	for attempt < maxLockRetry {
		result, err := locksClient.SetNX(context.Background(), getLockName(deckId), lockId, time.Second*2).Result()
		if err != nil {
			Logger.Errorf("Attempt %d: Failed to execute SetNX - %s\n", attempt, err)
		} else if !result {
			Logger.Errorf("Attempt %d: Key %s already exists\n", attempt, deckId)
		} else {
			return lockId, true
		}
		time.Sleep(time.Millisecond)
		attempt += 1
	}
	return lockId, false
}

func releaseLock(deckId, lockId string) error {
	key := getLockName(deckId)
	lockKey, err := locksClient.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return locksClient.Watch(context.Background(), func(tx *redis.Tx) error {
		if lockKey == lockId {
			return tx.Del(context.Background(), key).Err()
		}
		return fmt.Errorf(LOCK_HAS_CHANGED, deckId)
	}, key)

}

func getLockName(deckId string) string {
	return fmt.Sprintf("lock:%s", deckId)
}
