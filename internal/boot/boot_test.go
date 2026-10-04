package boot

import (
	"strings"
	"testing"
)

func TestSplitStatements(t *testing.T) {
	script := `
-- a leading comment that must be dropped
CREATE TABLE a (id INTEGER);

-- another comment
CREATE TABLE b (
    id INTEGER
);

CREATE INDEX idx_a ON a (id);
`
	statements := splitStatements(script)
	if len(statements) != 3 {
		t.Fatalf("got %d statements, want 3: %#v", len(statements), statements)
	}
	for _, statement := range statements {
		if strings.Contains(statement, "--") {
			t.Errorf("comment leaked into statement: %q", statement)
		}
		if strings.HasSuffix(statement, ";") {
			t.Errorf("statement still ends with a semicolon: %q", statement)
		}
	}
	if !strings.HasPrefix(statements[0], "CREATE TABLE a") {
		t.Errorf("unexpected first statement: %q", statements[0])
	}
}

func TestSplitStatementsEmpty(t *testing.T) {
	if got := splitStatements("\n-- only a comment\n"); len(got) != 0 {
		t.Fatalf("expected no statements, got %#v", got)
	}
}

func TestEmbeddedSchemaIsParsable(t *testing.T) {
	statements := splitStatements(schemaSQL)
	if len(statements) == 0 {
		t.Fatal("embedded schema produced no statements")
	}
	// Every table used by the DAO layer must be created by the schema.
	joined := strings.Join(statements, "\n")
	for _, table := range []string{"nav_category", "nav_site", "nav_setting", "nav_admin"} {
		if !strings.Contains(joined, table) {
			t.Errorf("schema does not create table %q", table)
		}
	}
}
