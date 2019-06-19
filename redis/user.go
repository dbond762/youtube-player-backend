package redis

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	player "github.com/dbond762/youtube-player-backend"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var (
	AuthError = errors.New("Authentication error")
)

type UserSession struct {
	Pool    *redis.Pool
	Service player.UserService
}

func NewUserService(addr string, service player.UserService) *UserSession {
	return &UserSession{
		Pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(addr)
			},
		},
		Service: service,
	}
}

func (s *UserSession) Login(login string) (*player.User, string, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	user, err := s.Service.UserByLogin(login)
	if err != nil {
		return nil, "", AuthError
	}

	token := uuid.New()
	key := fmt.Sprintf("token_%s", token)

	conn.Do("SET", key, user.ID)

	return user, token.String(), nil
}

func (s *UserSession) Logout(token string) {
	conn := s.Pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("token_%s", token)

	conn.Do("DEL", key)
}

func (s *UserSession) Authenticate(token string) (*player.User, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("token_%s", token)

	res, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return nil, AuthError
	} else if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(res)
	if err != nil {
		return nil, err
	}

	user, err := s.Service.UserByID(id)
	if err != nil {
		return nil, AuthError
	}

	return user, nil
}
