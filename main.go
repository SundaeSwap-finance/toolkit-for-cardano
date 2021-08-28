package main

import (
	"embed"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SundaeSwap-finance/toolkit-for-cardano/internal/cardano"
	"github.com/SundaeSwap-finance/toolkit-for-cardano/internal/gql"
	"github.com/SundaeSwap-finance/toolkit-for-cardano/internal/gql/graphiql"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/savaki/zapctx"
	"github.com/segmentio/ksuid"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

//go:embed internal/built.txt
var built string

//go:embed internal/version.txt
var version string

//go:embed ui/dist/*
var dist embed.FS

var opts struct {
	Assets  string // Assets contains optional directory for static assets
	Debug   bool   // Debug mode for additional logging
	Dir     string // Dir to store data in
	Port    int    // Port to listen on
	Cardano struct {
		CLI              cli.StringSlice // Cardano cli invocation e.g. cardano-cli or ssh hostname cardano-cli
		SocketPath       string          // SocketPath holds ${CARDANO_NODE_SOCKET_PATH}
		TestnetMagic     string          // TestnetMagic
		TreasuryAddr     string          // TreasuryAddr is the address of the treasury wallet
		TreasuryAddrFile string          // TreasuryAddrFile is a file that holds the address of the treasury wallet
		TreasurySkeyFile string          // TreasurySkeyFile is a pointer to the skey file for the treasury wallet
	}
}

func main() {
	app := cli.NewApp()
	app.Usage = "launch toolkit-for-cardano server"
	app.Version = strings.TrimSpace(version)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "assets",
			Usage:       "optional path to static assets",
			EnvVars:     []string{"ASSETS"},
			Destination: &opts.Assets,
		},
		&cli.StringSliceFlag{
			Name:        "cardano-cli",
			Usage:       "command to invoke cardano-cli",
			Value:       cli.NewStringSlice("cardano-cli"),
			EnvVars:     []string{"CARDANO_CLI"},
			Destination: &opts.Cardano.CLI,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "debug mode",
			Destination: &opts.Debug,
		},
		&cli.StringFlag{
			Name:        "dir",
			Usage:       "path to data directory",
			Value:       os.ExpandEnv("data"),
			EnvVars:     []string{"DATA_DIR"},
			Destination: &opts.Dir,
		},
		&cli.IntFlag{
			Name:        "port",
			Usage:       "port to listen on",
			Value:       3000,
			EnvVars:     []string{"PORT"},
			Destination: &opts.Port,
		},
		&cli.StringFlag{
			Name:        "socket-path",
			Usage:       "socket path for cardano node e.g. node.sock",
			EnvVars:     []string{"CARDANO_NODE_SOCKET_PATH"},
			Required:    true,
			Destination: &opts.Cardano.SocketPath,
		},
		&cli.StringFlag{
			Name:        "testnet-magic",
			Usage:       "testnet-magic value",
			Value:       "8",
			EnvVars:     []string{"TESTNET_MAGIC"},
			Destination: &opts.Cardano.TestnetMagic,
		},
		&cli.StringFlag{
			Name:        "treasury-addr",
			Usage:       "address with lovelace to fund other addresses from",
			EnvVars:     []string{"TREASURY_ADDR"},
			Destination: &opts.Cardano.TreasuryAddr,
		},
		&cli.StringFlag{
			Name:        "treasury-addr-file",
			Usage:       "load the treasury address from the specified file",
			EnvVars:     []string{"TREASURY_ADDR_FILE"},
			Destination: &opts.Cardano.TreasuryAddrFile,
		},
		&cli.StringFlag{
			Name:        "treasury-skey-file",
			Usage:       "file containing treasury signing key",
			EnvVars:     []string{"TREASURY_SIGNING_KEY_FILE"},
			Required:    true,
			Destination: &opts.Cardano.TreasurySkeyFile,
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func action(_ *cli.Context) error {
	dir, err := filepath.Abs(opts.Dir)
	if err != nil {
		return fmt.Errorf("failed to start toolkit-for-cardano: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "tmp"), 0755); err != nil {
		return fmt.Errorf("failed to start toolkit-for-cardano: failed to create tmp dir: %w", err)
	}

	// allow the treasury addr to be either provided or read from file
	addr := opts.Cardano.TreasuryAddr
	if addr == "" {
		data, err := ioutil.ReadFile(opts.Cardano.TreasuryAddrFile)
		if err != nil {
			return fmt.Errorf("failed to start toolkit-for-cardano: unable to read treasury-addr-file: %w", err)
		}
		addr = strings.TrimSpace(string(data))
	}

	cardanoCLI := cardano.CLI{
		Cmd:              opts.Cardano.CLI.Value(),
		Dir:              dir,
		SocketPath:       opts.Cardano.SocketPath,
		TestnetMagic:     opts.Cardano.TestnetMagic,
		TreasuryAddr:     addr,
		TreasurySkeyFile: opts.Cardano.TreasurySkeyFile,
		Debug:            opts.Debug,
	}
	config := gql.Config{
		Built:   strings.TrimSpace(built),
		CLI:     &cardanoCLI,
		Version: strings.TrimSpace(version),
	}
	handler, err := gql.New(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Use(
		withLogger(logger),
		withCORS(),
	)
	router.Get("/graphql", graphiql.New("/graphql"))
	router.Post("/graphql", handler.ServeHTTP)
	if opts.Assets != "" {
		fs := http.FileServer(http.Dir(opts.Assets))
		router.NotFound(fs.ServeHTTP)
	} else {
		fs := withRoot("ui/dist", http.FS(dist))
		router.NotFound(http.FileServer(fs).ServeHTTP)
	}

	logger.Info("started server", zap.Int("port", opts.Port))
	defer logger.Info("stopped server")

	return http.ListenAndServe(fmt.Sprintf(":%v", opts.Port), router)
}

type fileSystemFunc func(name string) (http.File, error)

func (fn fileSystemFunc) Open(name string) (http.File, error) {
	return fn(name)
}

func withRoot(root string, fs http.FileSystem) fileSystemFunc {
	return func(name string) (http.File, error) {
		path := filepath.Join(root, name)
		return fs.Open(path)
	}
}

func withCORS() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})
}

func withLogger(logger *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := zapctx.NewContext(req.Context(),
				logger.With(zap.String("tid", ksuid.New().String())),
			)
			req = req.WithContext(ctx)
			handler.ServeHTTP(w, req)
		})
	}
}
