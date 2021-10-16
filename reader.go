package ldifparser

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/ansel1/merry/v2"
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
func (r LdifReader) getEntityFromBlock(entityBlock Scanner) (entity.Entity, error) {
	entityLines := []string{}

	for entityBlock.Scan() {
		line := entityBlock.Text()
		if syntax.IsEntitySeparator(line) {
			break
		}

		entityLines = append(entityLines, line)
	}

	if entityBlock.Err() != nil {
		return entity.Entity{}, entityBlock.Err()
	}

	return entitybuilder.BuildEntity(entityLines, r.AttributeFilter)
}

type Scanner interface {
	Scan() bool
	Err() error
	Text() string
	Position() int64
}

func (r LdifReader) newScanner(readSrc io.Reader) Scanner {
	return poscanner.NewPositionedScanner(readSrc, r.ScannerBufferSize)
}

func (r *LdifReader) getScannerAtFirstEntityBlock() (Scanner, error) {
	scanner := r.newScanner(r.input)

	pos := scanner.Position()
	for scanner.Scan() {
		line := scanner.Text()
		if !syntax.IsLdifAttributeLine(line) {
			pos = scanner.Position()
			continue
		}

		r.input.Seek(pos, 0)
		return r.newScanner(r.input), nil
	}

	if scanner.Err() != nil {
		err := merry.Wrap(scanner.Err(), merry.AppendMessagef(
			"error at position [%d]", scanner.Position(),
		))
		return nil, err
	}

	return nil, errors.New("unable to locate first entity block")
}

// getKeyAddrOffset returns -1 if the entity is not found
func (r LdifReader) getKeyAttrOffset(keyAttr entity.Attribute) (int64, error) {
	keyAttrStr := strings.ToLower(StringifyAttribute(keyAttr)[0])
	r.Logger.Info("searching with key: \"%s\"", keyAttrStr)

	scanner := poscanner.NewPositionedScanner(r.input, r.ScannerBufferSize)
	pos := int64(-1)

	for scanner.Scan() {
		attrLine := strings.ToLower(scanner.Text())
		if keyAttrStr == attrLine {
			pos = scanner.Position()
			break
		}
		r.Logger.Debug("attrLine \"%s\" does not match key", attrLine)
	}

	return pos, scanner.Err()
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

	r.Logger.Info("parsing entity from block")
	entityScanner := r.newScanner(r.input)
	return r.getEntityFromBlock(entityScanner)
}

type EntityResp struct {
	Entity entity.Entity
	Error  error
}

// ReadEntities constructs an ldap entity per entry in the input ldif file.
func (r LdifReader) ReadEntities() []EntityResp {
	interrupt := make(chan bool)
	defer close(interrupt)

	results := r.ReadEntitiesChanneled(interrupt)
	entities := []EntityResp{}

	for resp := range results {
		entities = append(entities, resp)
	}

	return entities
}

func (r LdifReader) readSingleEntity(scanner Scanner) (e entity.Entity, err error) {
	r.Logger.Info("parsing entity")
	e, err = r.getEntityFromBlock(scanner)
	if err != nil {
		return
	}

	dn, dnFound := e.GetDN()
	if !dnFound {
		err = errors.New("entity corrupted, unable to parse DN for entity")
		return
	}

	r.Logger.Info("appending matched entity: " + dn)
	return
}

// ReadEntitiesChanneled constructs an ldap entity per entry in the input ldif file
// and returns the result via a channel. Any errors during processing will be packaged
// with the entity causing them and returned over the channel. Overflowing the scan buffer
// (line too long) will corrupt the scanner, causing a panic.
func (r LdifReader) ReadEntitiesChanneled(interrupt <-chan bool) <-chan EntityResp {
	results := make(chan EntityResp)

	go func() {
		defer close(results)

		select {
		case <-interrupt:
			return
		default:
		}

		r.Logger.Info("finding first entity block")
		scanner, err := r.getScannerAtFirstEntityBlock()
		if err != nil {
			resp := EntityResp{Error: err}
			results <- resp
			return
		}

		hasNextEntity := true
		for hasNextEntity {
			e, err := r.readSingleEntity(scanner)

			resp := EntityResp{e, err}
			results <- resp

			if err != nil && err == bufio.ErrTooLong {
				err = merry.Wrap(err, merry.WithMessagef(
					"panic caused by line at position: %d", scanner.Position(),
				))
				panic(err)
			}

			hasNextEntity = scanner.Scan()

			if err != nil && !r.ContinueOnErr {
				return
			}
		}
	}()

	return results
}
