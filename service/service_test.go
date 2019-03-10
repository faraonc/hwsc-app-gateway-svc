package service

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-app-gateway-svc/app"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/hwsc-org/hwsc-lib/hosts"
	"github.com/hwsc-org/hwsc-lib/logger"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"os"
	"testing"
)

var (
	// UserDB contains user database configs grabbed from env vars
	UserDB hosts.UserDBHost
)

func TestMain(m *testing.M) {
	logger.Info(consts.TestTag, "Initializing Test, this should ONLY print during unit tests")
	conf := config.NewConfig()
	src := env.NewSource(
		env.WithPrefix("hosts"),
	)
	if err := conf.Load(src); err != nil {
		logger.Fatal(consts.TestTag, "Failed to initialize configuration %v\n", err.Error())
	}
	// scan "hosts" prop "postgres" from environmental variables & copy values to UserDB struct
	if err := conf.Get("hosts", "postgres").Scan(&UserDB); err != nil {
		logger.Fatal(consts.TestTag, "Failed to get psql configuration", err.Error())
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s port=%s",
		UserDB.Host, UserDB.User, UserDB.Password, UserDB.Name, UserDB.SSLMode, UserDB.Port)

	fmt.Println(connectionString)
	postgresDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Fatal(consts.TestTag, "Failed to get psql connection", err.Error())
	}
	// create a postgres driver for migration
	driver, err := postgres.WithInstance(postgresDB, &postgres.Config{})
	if err != nil {
		logger.Fatal(consts.TestTag, "Failed to start postgres instance:", err.Error())
	}

	// create a migration instance
	migration, err := migrate.NewWithDatabaseInstance(
		"file://test_fixtures/psql",
		"postgres", driver,
	)
	if err != nil {
		logger.Fatal(consts.TestTag, "Failed to create a migration instance:", err.Error())
	}

	// run all migration up to the most active
	if err := migration.Up(); err != nil {
		logger.Fatal(consts.TestTag, "Failed to load active migration files:", err.Error())
	}

	// start the tests
	code := m.Run()
	os.Exit(code)
}

func TestGetStatus(t *testing.T) {
	cases := []struct {
		req         *pbsvc.AppGatewayServiceRequest
		serverState state
		expMsg      string
	}{
		{&pbsvc.AppGatewayServiceRequest{}, available, "OK"},
		{&pbsvc.AppGatewayServiceRequest{}, unavailable, "Unavailable"},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, _ := s.GetStatus(context.TODO(), c.req)
		assert.Equal(t, c.expMsg, res.GetMessage())
	}
}
