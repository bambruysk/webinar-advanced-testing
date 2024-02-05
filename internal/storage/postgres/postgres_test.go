package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"webinar-testing/pkg/models"
)

type testStorager interface {
	Storager
	clean(ctx context.Context) error
}

type PostgresTestSuite struct {
	suite.Suite
	testStorager

	tc  *tcpostgres.PostgresContainer
	cfg *Config
}

type testLogger struct {
	t *testing.T
}

func (t testLogger) Info(args ...any) {
	t.t.Log(args...)
}

func (ts *PostgresTestSuite) SetupSuite() {
	cfg := &Config{
		ConnectTimeout:   5 * time.Second,
		QueryTimeout:     5 * time.Second,
		Username:         "postgres",
		Password:         "password",
		DBName:           "postgres",
		MigrationVersion: 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgc, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:latest"),
		tcpostgres.WithDatabase(cfg.DBName),
		tcpostgres.WithUsername(cfg.Username),
		tcpostgres.WithPassword(cfg.Password),
		tcpostgres.WithInitScripts(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	require.NoError(ts.T(), err)

	cfg.Host, err = pgc.Host(ctx)
	require.NoError(ts.T(), err)

	port, err := pgc.MappedPort(ctx, "5432")
	require.NoError(ts.T(), err)

	cfg.Port = uint16(port.Int())

	ts.tc = pgc
	ts.cfg = cfg

	db, err := New(cfg)
	require.NoError(ts.T(), err)

	ts.testStorager = db

	ts.T().Logf("stared postgres at %s:%d", cfg.Host, cfg.Port)
}

func (ts *PostgresTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) TestDummy() {}

func (s *storage) clean(ctx context.Context) error {
	newCtx, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := s.pool.Exec(newCtx, "DELETE FROM orders")
	return err
}

func (ts *PostgresTestSuite) SetupTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostgresTestSuite) TearDownTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostgresTestSuite) TestAdd() {
	order := models.Order{
		UserID: "user1",
		Goods: map[models.GoodID]int{
			"good1": 3,
		},
	}

	ts.NoError(ts.Add(context.Background(), order))

	goods, err := ts.ListByUser(context.Background(), order.UserID)
	ts.NoError(err)

	ts.Equal(order.Goods, goods.Goods)
}

func (ts *PostgresTestSuite) TestAddToExist() {
	userID := models.UserID("user1")
	goodID := models.GoodID("good1")
	order1 := models.Order{
		UserID: userID,
		Goods: map[models.GoodID]int{
			goodID: 3,
		},
	}

	order2 := models.Order{
		UserID: userID,
		Goods: map[models.GoodID]int{
			goodID: 5,
		},
	}

	ts.NoError(ts.Add(context.Background(), order1))
	ts.NoError(ts.Add(context.Background(), order2))

	goods, err := ts.ListByUser(context.Background(), userID)
	ts.NoError(err)

	ts.Equal(order1.Goods[goodID]+order2.Goods[goodID], goods.Goods[goodID])
}

func (ts *PostgresTestSuite) TestList() {
	order1 := models.Order{
		UserID: "user1",
		Goods: map[models.GoodID]int{
			"good1": 1,
			"good2": 2,
			"good3": 3,
			"good4": 4,
		},
	}

	order2 := models.Order{
		UserID: "user2",
		Goods: map[models.GoodID]int{
			"good1": 5,
			"good2": 6,
			"good7": 7,
			"good8": 8,
		},
	}

	ts.NoError(ts.Add(context.Background(), order1))
	ts.NoError(ts.Add(context.Background(), order2))

	goods, err := ts.ListByUser(context.Background(), order1.UserID)
	ts.NoError(err)

	ts.Equal(order1.Goods, goods.Goods)
}

func (ts *PostgresTestSuite) TestListNotExist() {
	goods, err := ts.ListByUser(context.Background(), "not existed user")
	ts.NoError(err)

	ts.Equal(map[models.GoodID]int{}, goods.Goods)
}

func (ts *PostgresTestSuite) TestDelete() {
	userID := models.UserID("user1")
	orderAdd := models.Order{
		UserID: userID,
		Goods: map[models.GoodID]int{
			"good1": 11,
			"good2": 12,
			"good3": 13,
			"good4": 14,
		},
	}

	orderDelete := models.Order{
		UserID: userID,
		Goods: map[models.GoodID]int{
			"good1": 1,
			"good2": 3,
			"good3": 5,
		},
	}

	ts.NoError(ts.Add(context.Background(), orderAdd))
	ts.NoError(ts.Delete(context.Background(), orderDelete))

	goods, err := ts.ListByUser(context.Background(), userID)
	ts.NoError(err)

	ts.Equal(map[models.GoodID]int{
		"good1": 10,
		"good2": 9,
		"good3": 8,
		"good4": 14,
	}, goods.Goods)
}
