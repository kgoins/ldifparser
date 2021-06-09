package ldifparser_test

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/kgoins/ldifparser"
	"github.com/stretchr/testify/require"
)

var testFileName string = "test_users.ldif"
var numTestFileEntities int = 3

func getTestDataDir() string {
	_, mypath, _, _ := runtime.Caller(0)
	myparent := path.Dir(mypath)
	return path.Join(myparent, "testdata")
}

func TestReader_BuildEntity(t *testing.T) {
	r := require.New(t)
	testFilePath := filepath.Join(getTestDataDir(), testFileName)
	testFile, err := os.Open(testFilePath)
	r.NoError(err)
	defer testFile.Close()

	testAttr := "sAMAccountName"
	testName := "DISABLEDUSER"

	ldifReader := ldifparser.NewLdifReader(testFile)
	e, err := ldifReader.BuildEntity(testAttr, testName)
	r.NoError(err)

	r.False(e.IsEmpty())
	r.NotEmpty(e.GetDN())

	eName, found := e.GetSingleValuedAttribute(testAttr)
	r.True(found)
	r.Equal(eName, testName)
}

func TestReader_BuildEntities(t *testing.T) {
	r := require.New(t)

	testFilePath := filepath.Join(getTestDataDir(), testFileName)
	testFile, err := os.Open(testFilePath)

	r.NoError(err)
	defer testFile.Close()

	ldifReader := ldifparser.NewLdifReader(testFile)
	entities, err := ldifReader.BuildEntities()
	r.NoError(err)

	r.Equal(numTestFileEntities, len(entities))

	s := rand.NewSource(time.Now().Unix())
	rng := rand.New(s)

	randEntity := entities[rng.Intn(len(entities))]
	r.False(randEntity.IsEmpty())
}
