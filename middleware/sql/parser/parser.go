package parser

import (
	"github.com/percona/go-mysql/query"
	"github.com/pingcap/parser"
	"github.com/romberli/go-util/constant"
)

type Parser struct {
	Parser  *parser.Parser
	Visitor *Visitor
}

// NewParser returns a new *Parser
func NewParser(visitor *Visitor) *Parser {
	return &Parser{
		Parser:  parser.New(),
		Visitor: visitor,
	}
}

// NewParserWithDefault returns a new *Parser with default visitor
func NewParserWithDefault() *Parser {
	return &Parser{
		Parser:  parser.New(),
		Visitor: NewVisitorWithDefault(),
	}
}

// Parse parses sql and returns the result
func (p *Parser) Parse(sql string) (*Result, []error, error) {
	stmtNodes, warns, err := p.Parser.Parse(sql, constant.EmptyString, constant.EmptyString)
	if warns != nil || err != nil {
		return nil, warns, err
	}

	for _, stmtNode := range stmtNodes {
		stmtNode.Accept(p.Visitor)
	}

	return p.Visitor.result, nil, nil
}

// Split splits multiple sqls into a slice
func (p *Parser) Split(sqls string) ([]string, []error, error) {
	var sqlList []string

	stmtNodes, warns, err := p.Parser.Parse(sqls, constant.EmptyString, constant.EmptyString)
	if warns != nil || err != nil {
		return nil, warns, err
	}

	for _, stmtNode := range stmtNodes {
		sqlList = append(sqlList, stmtNode.Text())
	}

	return sqlList, nil, nil
}

// GetFingerprint returns fingerprint fo given sql
func (p *Parser) GetFingerprint(sql string) string {
	return query.Fingerprint(sql)
}

func (p *Parser) GetSQLID(sql string) string {
	return query.Id(p.GetFingerprint(sql))
}