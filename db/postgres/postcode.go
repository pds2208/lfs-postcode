package postgres

import (
    "lfs-postcode/types"
)

func (p PostgresConnection) GetPostCodeDifferences(year int) ([]types.PostCodeDifferences, error) {
    var differences []types.PostCodeDifferences

    all := p.DB.Collection(p.PostCodeDifferences)

    err := all.Find().All(&differences)
    if err != nil {
        return nil, err
    }

    return differences, nil
}
