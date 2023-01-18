package database_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	junodb "github.com/forbole/juno/v4/database"
	junodbcfg "github.com/forbole/juno/v4/database/config"
	"github.com/forbole/juno/v4/logging"

	"github.com/stretchr/testify/suite"

	"github.com/hjcore/gjuno/database"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database *database.Db
}

func (suite *DbTestSuite) SetupTest() {
	// Build the database
	encodingConfig := simapp.MakeTestEncodingConfig()
	databaseConfig := junodbcfg.NewDatabaseConfig(
		"postgres://gjuno:password@localhost:6432/gjuno?sslmode=disable&search_path=public",
		10,
		10,
		0,
		0,
	)

	db, err := database.Builder(junodb.NewContext(databaseConfig, &encodingConfig, logging.DefaultLogger()))
	suite.Require().NoError(err)

	gotabitDb, ok := (db).(*database.Db)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = gotabitDb.SQL.Exec(fmt.Sprintf(`DROP SCHEMA %s CASCADE;`, databaseConfig.GetSchema()))
	suite.Require().NoError(err)

	// Re-create the schema
	_, err = gotabitDb.SQL.Exec(fmt.Sprintf(`CREATE SCHEMA %s;`, databaseConfig.GetSchema()))
	suite.Require().NoError(err)

	dirPath := "schema"
	dir, err := ioutil.ReadDir(dirPath)
	for _, fileInfo := range dir {
		if !strings.HasSuffix(fileInfo.Name(), ".sql") {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			_, err := gotabitDb.SQL.Exec(commentsRegExp.ReplaceAllString(request, ""))
			suite.Require().NoError(err)
		}
	}

	suite.database = gotabitDb
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
