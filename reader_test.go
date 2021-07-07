package ldifparser_test

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/kgoins/ldapentity/entity/ad"
	"github.com/kgoins/ldifparser"
	"github.com/kgoins/ldifparser/entitybuilder"
	"github.com/stretchr/testify/require"
)

var testFileName string = "test_users.ldif"
var numTestFileEntities int = 3

func getTestDataDir() string {
	_, mypath, _, _ := runtime.Caller(0)
	myparent := path.Dir(mypath)
	return path.Join(myparent, "testdata")
}

func TestReader_ReadEntity(t *testing.T) {
	r := require.New(t)
	testFilePath := filepath.Join(getTestDataDir(), testFileName)
	testFile, err := os.Open(testFilePath)
	r.NoError(err)
	defer testFile.Close()

	testAttr := "sAMAccountName"
	testName := "DISABLEDUSER"

	ldifReader := ldifparser.NewLdifReader(testFile)
	e, err := ldifReader.ReadEntity(testAttr, testName)
	r.NoError(err)

	r.False(e.IsEmpty())
	r.NotEmpty(e.GetDN())

	eName, found := e.GetSingleValuedAttribute(testAttr)
	r.True(found)
	r.Equal(eName, testName)
}

func TestReader_ReadEntityWithAttrFilter(t *testing.T) {
	r := require.New(t)

	testFilePath := filepath.Join(getTestDataDir(), testFileName)
	testFile, err := os.Open(testFilePath)
	r.NoError(err)
	defer testFile.Close()

	testAttr := ad.ATTR_sAMAccountName
	testName := "DISABLEDUSER"

	conf := ldifparser.NewReaderConf()
	conf.AttributeFilter = entitybuilder.NewAttributeFilter(ad.ATTR_sAMAccountName)

	ldifReader := ldifparser.NewLdifReader(testFile, conf)
	e, err := ldifReader.ReadEntity(testAttr, testName)
	r.NoError(err)

	want := []string{
		ad.ATTR_DN,
		ad.ATTR_sAMAccountName,
	}

	attrNames := e.GetAllAttributeNames()
	sort.Strings(attrNames)

	r.Equal(want, attrNames)
}

func TestReader_ReadEntities(t *testing.T) {
	r := require.New(t)

	testFilePath := filepath.Join(getTestDataDir(), testFileName)
	testFile, err := os.Open(testFilePath)

	r.NoError(err)
	defer testFile.Close()

	ldifReader := ldifparser.NewLdifReader(testFile)
	entities, err := ldifReader.ReadEntities()
	r.NoError(err)

	r.Equal(numTestFileEntities, len(entities))

	s := rand.NewSource(time.Now().Unix())
	rng := rand.New(s)

	randEntity := entities[rng.Intn(len(entities))]
	r.False(randEntity.IsEmpty())
}

func TestReader_ReadEntitiesWithoutPrologue(t *testing.T) {
	r := require.New(t)

	testFilePath := filepath.Join(getTestDataDir(), "user_no_prologue.ldif")
	testFile, err := os.Open(testFilePath)

	r.NoError(err)
	defer testFile.Close()

	ldifReader := ldifparser.NewLdifReader(testFile)
	entities, err := ldifReader.ReadEntities()
	r.NoError(err)

	r.Equal(1, len(entities))
	r.NotEmpty(entities[0].GetDN())
}
