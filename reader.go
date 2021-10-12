package ldifparser

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/kgoins/backscanner"

	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/poscanner"

	"github.com/kgoins/ldifparser/entitybuilder"
	"github.com/kgoins/ldifparser/syntax"
)

type ReadSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}

// LdifReader constructs LDAP Entities from an ldif file.
type LdifReader struct {
	input ReadSeekerAt
	ReaderConf
}

// NewLdifReader returns a constructed LdifReader.
// The ReadSeekerAt input requirement ensures that the input implements both
// io.ReaderAt and io.ReadSeeker. A file or string buffer are good examples.
func NewLdifReader(input ReadSeekerAt, conf ...ReaderConf) LdifReader {
	var actualConf ReaderConf
	if len(conf) == 0 {
		actualConf = NewReaderConf()
	} else {
		actualConf = conf[0]
	}

	return LdifReader{
		input:      input,
		ReaderConf: actualConf,
	}
}

// SetAttributeFilter modifies the r to return only the ldap
// attributes present in the filter on entites that are parsed.
func (r *LdifReader) SetAttributeFilter(filter entitybuilder.AttributeFilter) {
	r.AttributeFilter = filter
}

// getEntityFromBlock constructs an entity from the lines starting
// at the scanner's current position. At the end of this call, the
// scanner will be positioned at the end of the entity.
func (r LdifReader) getEntityFromBlock(entityBlock *bufio.Scanner) (entity.Entity, error) {
	entityLines := []string{}

	for entityBlock.Scan() {
		line := entityBlock.Text()
		if syntax.IsEntitySeparator(line) {
			break
		}

		entityLines = append(entityLines, line)
	}

	return entitybuilder.BuildEntity(entityLines, r.AttributeFilter)
}

func (r *LdifReader) getScannerAtFirstEntityBlock() (*bufio.Scanner, error) {
	scanner := poscanner.NewPositionedScanner(r.input, r.ScannerBufferSize)

	pos := scanner.Position()
	for scanner.Scan() {
		line := scanner.Text()
		if !syntax.IsLdifAttributeLine(line) {
			pos = scanner.Position()
			continue
		}

		r.input.Seek(pos, 0)
		return bufio.NewScanner(r.input), nil
	}

	return nil, errors.New("unable to locate first entity block")
}

// getKeyAddrOffset returns -1 if the entity is not found
func (r LdifReader) getKeyAttrOffset(keyAttr entity.Attribute) (int64, error) {
	keyAttrStr := strings.ToLower(StringifyAttribute(keyAttr)[0])
	r.Logger.Info("searching with key: \"%s\"", keyAttrStr)

	scanner := poscanner.NewPositionedScanner(r.input)

	for scanner.Scan() {
		attrLine := strings.ToLower(scanner.Text())
		if keyAttrStr == attrLine {
			return scanner.Position(), nil
		}
		r.Logger.Debug("attrLine \"%s\" does not match key", attrLine)
	}

	return -1, errors.New("entity not found")
}

func (r LdifReader) getPrevEntityOffset(input io.ReaderAt, lineOffset int64) (int, error) {
	scanner := backscanner.New(input, int(lineOffset))
	for {
		line, pos, err := scanner.Line()
		if err != nil {
			return pos, err
		}

		if syntax.IsEntityTitle(line) {
			return pos, nil
		}
	}
}

// ReadEntity returns an empty Entity object if the object is not found,
// other wise it returns the entity object or an error if one is encountered.
func (r LdifReader) ReadEntity(keyAttrName string, keyAttrVal string) (e entity.Entity, err error) {
	keyAttr := entity.NewEntityAttribute(keyAttrName, keyAttrVal)

	keyAttrOffset, err := r.getKeyAttrOffset(keyAttr)
	if err != nil {
		return
	}
	r.Logger.Debug("key found at position: %d", keyAttrOffset)

	// Object not found
	if keyAttrOffset == -1 {
		return
	}

	entityOffset, err := r.getPrevEntityOffset(r.input, keyAttrOffset)
	if err != nil {
		return
	}
	r.Logger.Info("entity found at position: %d", entityOffset)

	_, err = r.input.Seek(int64(entityOffset), 0)
	if err != nil {
		return
	}

	entityScanner := bufio.NewScanner(r.input)
	buf := make([]byte, r.ScannerBufferSize)
	entityScanner.Buffer(buf, r.ScannerBufferSize)

	r.Logger.Info("parsing entity from block")
	return r.getEntityFromBlock(entityScanner)
}

// ReadEntities constructs an ldap entity per entry in the input ldif file.
func (r LdifReader) ReadEntities() ([]entity.Entity, error) {
	done := make(chan bool)
	results := r.ReadEntitiesChanneled(done)

	entities := []entity.Entity{}
	for e := range results {
		entities = append(entities, e)
	}

	return entities, nil
}

// ReadEntitiesChanneled constructs an ldap entity per entry in the input ldif file
// and returns the result via a channel. Any errors during processing will be logged
// and the entity that caused it will be skipped.
func (r LdifReader) ReadEntitiesChanneled(done <-chan bool) <-chan entity.Entity {
	results := make(chan entity.Entity)

	go func() {
		defer close(results)

		select {
		case <-done:
			return
		default:
		}

		r.Logger.Info("finding first entity block")
		scanner, err := r.getScannerAtFirstEntityBlock()
		if err != nil {
			r.Logger.Error(err.Error())
			return
		}

		hasNextEntity := true
		for hasNextEntity {
			r.Logger.Info("parsing entity")
			entity, parseErr := r.getEntityFromBlock(scanner)
			if parseErr != nil {
				r.Logger.Error(parseErr.Error())
				return
			}

			dn, dnFound := entity.GetDN()
			if !dnFound {
				r.Logger.Error("entity corrupted, unable to parse DN for entity")
				return
			}

			r.Logger.Info("appending matched entity: " + dn)
			results <- entity

			hasNextEntity = scanner.Scan()
		}
	}()

	return results
}
