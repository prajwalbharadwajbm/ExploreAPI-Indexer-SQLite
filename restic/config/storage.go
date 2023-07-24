package config

type Storage interface {
	Load() error
}
