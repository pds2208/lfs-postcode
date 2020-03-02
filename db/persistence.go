package db

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "lfs-postcode/config"
    "lfs-postcode/db/postgres"
    "lfs-postcode/types"
    "sync"
)

var cachedConnection DVPersistence
var connectionMux = &sync.Mutex{}

func GetDVPersistenceImpl() (DVPersistence, error) {
    connectionMux.Lock()
    defer connectionMux.Unlock()

    if cachedConnection != nil {
        log.Info().
            Str("databaseName", config.Config.Database.Database).
            Msg("Returning cached database connection")
        return cachedConnection, nil
    }

    cachedConnection = &postgres.PostgresConnection{}

    if err := cachedConnection.Connect(); err != nil {
        log.Info().
            Err(err).
            Str("databaseName", config.Config.Database.Database).
            Msg("Cannot connect to database")
        cachedConnection = nil
        return nil, fmt.Errorf("cannot connect to database")
    }

    return cachedConnection, nil
}

type DVPersistence interface {
    Connect() error
    Close()

    GetPostCodeDifferences(year int) ([]types.PostCodeDifferences, error)
}
