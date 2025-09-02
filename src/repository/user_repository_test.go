package repository

import (
	"auth-service/src/domain"
	"auth-service/src/test_artefacts/seeder"
	"auth-service/src/test_artefacts/stubs"
	"context"
	"errors"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestUserRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserRepository Integration Suite")
}

var (
	db       *pgxpool.Pool
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

var _ = BeforeSuite(func() {
	var err error
	pool, err = dockertest.NewPool("")
	Expect(err).NotTo(HaveOccurred())

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres", Tag: "15-alpine",
		Env: []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=test_db"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	Expect(err).NotTo(HaveOccurred())

	err = pool.Retry(func() error {
		dbURL := "postgres://postgres:secret@" + resource.GetHostPort("5432/tcp") + "/test_db?sslmode=disable"
		db, err = pgxpool.New(context.Background(), dbURL)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	})
	Expect(err).NotTo(HaveOccurred())

	migration, err := os.ReadFile("../../database/000001_create_users_table.up.sql")
	Expect(err).NotTo(HaveOccurred())
	_, err = db.Exec(context.Background(), string(migration))
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(pool.Purge(resource)).To(Succeed())
})

var _ = Describe("UserRepository", func() {
	var userRepo UserRepository
	var testSeeder *seeder.TestSeeder
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		userRepo = NewUser(db)
		testSeeder = seeder.NewTestSeeder(db)

		err := testSeeder.TruncateTables(ctx)
		Expect(err).NotTo(HaveOccurred())
	})

	// --- Cenários de Teste ---
	Describe("Creating a new user", func() {
		Context("when the email already exists", func() {
			It("should return an ErrEmailAlreadyExists error", func() {
				// Arrange
				user1 := stubs.NewUserStub().WithEmail("duplicate@example.com").Get()
				err := testSeeder.InsertUser(ctx, user1)
				Expect(err).NotTo(HaveOccurred())

				user2 := stubs.NewUserStub().WithEmail("duplicate@example.com").Get()

				// Act
				err = userRepo.Create(ctx, user2)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, domain.ErrEmailAlreadyExists)).To(BeTrue())
			})
		})
	})
})
