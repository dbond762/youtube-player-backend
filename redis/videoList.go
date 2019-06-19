package redis

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	player "github.com/dbond762/youtube-player-backend"
	"github.com/gomodule/redigo/redis"
)

type VideoFinder struct {
	Pool   *redis.Pool
	Finder player.VideoFinder
}

const sleepTime = 250 * time.Millisecond

func NewVideoSearcher(addr string, finder player.VideoFinder) *VideoFinder {
	return &VideoFinder{
		Pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(addr)
			},
		},
		Finder: finder,
	}
}

func (f *VideoFinder) Search(query string) (*player.VideoList, error) {
	conn := f.Pool.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("EXISTS", "lock_query_"+query))
	if err != nil && err != redis.ErrNil {
		log.Printf("Redis: error on get query lock: %s", err)
		return nil, err
	}
	if reply == 1 {
		time.Sleep(sleepTime)
	}

	res, err := redis.Bytes(conn.Do("GET", "query_"+query))
	if err == redis.ErrNil {
		conn.Do("SET", "lock_query_"+query, 1)
		defer conn.Do("DEL", "lock_query_"+query)

		list, err := f.Finder.Search(query)
		if err != nil {
			log.Printf("Redis: error on search video: %s", err)
			return nil, err
		}

		var buf bytes.Buffer

		encoder := gob.NewEncoder(&buf)
		err = encoder.Encode(list)
		if err != nil {
			log.Printf("Redis: error on encoding: %s", err)
			return nil, err
		}

		const ttl = 24 * 60 * 60
		_, err = conn.Do("SET", "query_"+query, buf.Bytes(), "EX", ttl)
		if err != nil {
			log.Printf("Redis: error on set query: %s", err)
			return nil, err
		}

		log.Printf("Redis: Search from api")

		return list, nil
	} else if err != nil {
		log.Printf("Redis: error on get query: %s", err)
		return nil, err
	}

	var list player.VideoList

	decoder := gob.NewDecoder(bytes.NewBuffer(res))
	err = decoder.Decode(&list)
	if err != nil {
		log.Printf("Redis: error on decoding: %s", err)
		return nil, err
	}

	log.Printf("Redis: Search from cache")

	return &list, nil
}
