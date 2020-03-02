package config

type configuration struct {
    LogFormat     string
    LogLevel      string
    TmpDirectory  string
    TestDirectory string
    MonthlyTables MonthlyTablesConfiguration
    Database      DatabaseConfiguration
    Service       ServiceConfiguration
}
