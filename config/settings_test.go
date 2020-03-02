package config_test

import (
    conf "lfs-postcode/config"
    "testing"
)

func TestConfig(t *testing.T) {

    server := conf.Config.Database.Server
    if server != "localhost" {
        t.Errorf("server = %s; want localhost", server)
    } else {
        t.Logf("Server %s\n", server)
    }

    user := conf.Config.Database.User
    if user != "lfs" {
        t.Errorf("user = %s; want lfs", user)
    } else {
        t.Logf("user %s\n", user)
    }

    password := conf.Config.Database.Password
    if password != "lfs" {
        t.Errorf("password = %s; want lfs", password)
    } else {
        t.Logf("password %s\n", password)
    }

    databaseName := conf.Config.Database.Database
    if databaseName != "lfs" {
        t.Errorf("database name = %s; want lfs", databaseName)
    } else {
        t.Logf("database name %s\n", databaseName)
    }

    maxPoolsize := conf.Config.Database.ConnectionPool.MaxPoolSize
    if maxPoolsize != 10 {
        t.Errorf("maxPoolsize = %d; want 10", maxPoolsize)
    } else {
        t.Logf("maxPoolsize %d\n", maxPoolsize)
    }

    postCodeDifferences := conf.Config.MonthlyTables.PostCodeDifferences
    if postCodeDifferences != "postcode_differences" {
        t.Errorf("postcode_differences = %s; want 'postcode_differences'", postCodeDifferences)
    } else {
        t.Logf("postcode_differences: %s\n", postCodeDifferences)
    }

}
