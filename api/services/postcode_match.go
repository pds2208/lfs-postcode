package services

import (
    "github.com/rs/zerolog/log"
    "lfs-postcode/db"
    "lfs-postcode/types"
)

type PostCodeMatchService struct {}

func (pm PostCodeMatchService) PostCodeDifferences(year int) ([]types.PostCodeDifferences, error){

    // Database connection
    dbase, err := db.GetDVPersistenceImpl()
    if err != nil {
        log.Error().Err(err)
        return nil, err
    }

    res, err := dbase.GetPostCodeDifferences(year)
    if err != nil {
        return nil, err
    }
    return res, nil
}
