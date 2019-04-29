package service

import (
	"database/sql"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-app-gateway-svc/app"
	pbuser "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-user-svc/user"
	pblib "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-app-gateway-svc/consts"
	"github.com/hwsc-org/hwsc-lib/auth"
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
		logger.Error(consts.TestTag, "Failed to load active migration files:", err.Error())
		// hack to reset DB to default settings with no entries
		logger.Info(consts.TestTag, "Resetting migration")
		if downErr := migration.Down(); downErr != nil {
			logger.Fatal(consts.TestTag, "Failed to migrate down", err.Error())
		}
		if upErr := migration.Up(); upErr != nil {
			logger.Fatal(consts.TestTag, "Failed to load active migration files:", err.Error())
		}
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

func TestCreateUser(t *testing.T) {
	validEmail := randomdata.Email()

	cases := []struct {
		desc        string
		user        *pblib.User
		isExpErr    bool
		errStr      string
		serverState state
	}{
		{
			"Test for unavailable service",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: "TestOrg",
			},
			true,
			"rpc error: code = Unavailable desc = service is unavailable",
			unavailable,
		},
		{
			"Test for nil user",
			nil,
			true,
			"rpc error: code = InvalidArgument desc = nil request",
			available,
		},
		{
			"Test for valid user registration",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: "TestOrg",
			},
			false,
			"",
			available,
		},
		{
			"Test for duplicate user email",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "Abcd!123@",
				Organization: "TestOrg",
			},
			true,
			`rpc error: code = Internal desc = pq: duplicate key value violates unique constraint "accounts_email_key"`,
			available,
		},
		{
			"Test for empty string",
			&pblib.User{
				FirstName: "",
			},
			true,
			"rpc error: code = Internal desc = invalid User first name",
			available,
		},
		{
			"Test for missing password",
			&pblib.User{
				FirstName:    randomdata.FirstName(randomdata.Male),
				LastName:     randomdata.LastName(),
				Email:        validEmail,
				Password:     "",
				Organization: "TestOrg",
			},
			true,
			"rpc error: code = Internal desc = invalid User password",
			available,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		resp, err := s.CreateUser(context.TODO(), &pbsvc.AppGatewayServiceRequest{
			UserRequest: &pbuser.UserRequest{
				User: c.user,
			},
		})
		if c.isExpErr {
			assert.Nil(t, resp, c.desc)
			assert.EqualError(t, err, c.errStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, resp.GetUser().GetUuid(), c.desc)
			assert.Equal(t, auth.PermissionStringMap[auth.NoPermission], resp.GetUser().GetPermissionLevel(), c.desc)
		}
	}
}
