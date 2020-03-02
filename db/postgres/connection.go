package postgres

import (
	"github.com/rs/zerolog/log"
	"lfs-postcode/config"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

type PostgresConnection struct {
	DB sqlbuilder.Database
	config.MonthlyTablesConfiguration
}

func (p *PostgresConnection) Connect() error {

	var settings = postgresql.ConnectionURL{
		Database: config.Config.Database.Database,
		Host:     config.Config.Database.Server,
		User:     config.Config.Database.User,
		Password: config.Config.Database.Password,
	}

	p.LoadConfiguration()

	log.Debug().
		Str("databaseName", config.Config.Database.Database).
		Msg("Connecting to database")

	sess, err := postgresql.Open(settings)

	if err != nil {
		log.Error().
			Err(err).
			Str("databaseName", config.Config.Database.Database).
			Msg("Cannot connect to database")
		return err
	}

	log.Debug().
		Str("databaseName", config.Config.Database.Database).
		Msg("Connected to database")

	if config.Config.Database.Verbose {
		sess.SetLogging(true)
	}

	p.DB = sess

	poolSize := config.Config.Database.ConnectionPool.MaxPoolSize
	maxIdle := config.Config.Database.ConnectionPool.MaxIdleConnections
	maxLifetime := config.Config.Database.ConnectionPool.MaxLifetimeSeconds

	if maxLifetime > 0 {
		maxLifetime = maxLifetime * time.Second
		sess.SetConnMaxLifetime(maxLifetime)
	}

	log.Debug().
		Int("MaxPoolSize", poolSize).
		Int("MaxIdleConnections", maxIdle).
		Dur("MaxLifetime", maxLifetime*time.Second).
		Msg("Connection Attributes")

	sess.SetMaxOpenConns(poolSize)
	sess.SetMaxIdleConns(maxIdle)

	return nil
}

func (p PostgresConnection) Close() {
	if p.DB != nil {
		_ = p.DB.Close()
	}
}

func (s *PostgresConnection) LoadConfiguration() {
	s.PostCodeDifferences = config.Config.MonthlyTables.PostCodeDifferences
	if s.PostCodeDifferences == "" {
		panic("PostCodeDifferences table configuration not set")
	}
}