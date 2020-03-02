package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    h "lfs-postcode/api/handlers"
    "lfs-postcode/config"
    "lfs-postcode/messaging"
    "net/http"
    "os"
    "time"
)

const applicationName = "LFS PostCode Match"

func main() {

    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

    if config.Config.LogFormat == "Terminal" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
    }

    // Command line flag overrides the configuration file
    debug := flag.Bool("debug", false, "sets log level to debug")

    router := mux.NewRouter()

    flag.Parse()
    if *debug || config.Config.LogLevel == "Debug" {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
        router.Use(loggingMiddleware)
    } else {
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }

    log.Info().
        Str("startTime", time.Now().String()).
        Msg(fmt.Sprintf("%s: Starting up", applicationName))

    postCodeMatcher := h.NewPostCodeMatchHandler()
    router.HandleFunc("/postcode/{year}/{month}", postCodeMatcher.PostCodeMatchingHandler).Methods(http.MethodGet)

    listenAddress := config.Config.Service.ListenAddress

    writeTimeout, err := time.ParseDuration(config.Config.Service.WriteTimeout)
    if err != nil {
        log.Fatal().
            Err(err).
            Str("service", "LFS").
            Msgf("writeTimeout configuration error")
    }

    readTimeout, err := time.ParseDuration(config.Config.Service.ReadTimeout)
    if err != nil {
        log.Fatal().
            Err(err).
            Str("service", "LFS").
            Msgf("readTimeout configuration error")
    }

    // we'll allow anything for now. May need or want to restrict this to just the UI when we know its endpoint
    origins := []string{"*"}
    var cors = handlers.AllowedOrigins(origins)

    handlers.CORS(cors)(router)

    srv := &http.Server{
        Handler:      router,
        Addr:         listenAddress,
        WriteTimeout: writeTimeout,
        ReadTimeout:  readTimeout,
    }

    log.Info().
        Str("listenAddress", listenAddress).
        Str("writeTimeout", writeTimeout.String()).
        Str("readTimeout", readTimeout.String()).
        Msg(fmt.Sprintf("%s: Waiting for requests", applicationName))

    messaging.NewConnection()

    err = srv.ListenAndServe()
    log.Fatal().
        Err(err).
        Str("service", "LFS").
        Msgf("ListenAndServe failed")
}

func loggingMiddleware(next http.Handler) http.Handler {

    log.Info().Msg("Logging middleware registered")

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Debug().
            Str("URI:", r.RequestURI).
            Str("client", r.RemoteAddr).
            Msg("-> Received request")
        startTime := time.Now()
        next.ServeHTTP(w, r)

        log.Debug().
            Str("URI:", r.RequestURI).
            Str("elapsedTime", FmtDuration(startTime)).
            Msg("<- Request Completed")
    })
}

func FmtDuration(t time.Time) string {
    const (
        Decisecond = 100 * time.Millisecond
        Day        = 24 * time.Hour
    )
    ts := time.Since(t)
    sign := time.Duration(1)
    if ts < 0 {
        sign = -1
        ts = -ts
    }
    ts += +Decisecond / 2
    d := sign * (ts / Day)
    ts = ts % Day
    h := ts / time.Hour
    ts = ts % time.Hour
    m := ts / time.Minute
    ts = ts % time.Minute
    s := ts / time.Second
    ts = ts % time.Second
    f := ts / Decisecond
    return fmt.Sprintf("%dd:%02dh:%02dm:%02d.%d02s", d, h, m, s, f)
}
