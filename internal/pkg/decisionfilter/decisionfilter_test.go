package decisionfilter

import (
	"testing"

	"github.com/lapt0r/goose/internal/pkg/decisiontree"
)

func TestTokenize(t *testing.T) {
	var testString = "a b c"
	var result = tokenize(testString)

	if len(result) != 3 {
		t.Errorf("Expected 3 items but found %v", len(result))
	}

	if result[0] != "a" {
		t.Errorf("Expected first element to be 'a'  but was '%v'", result[0])
	}
	if result[1] != "b" {
		t.Errorf("Expected second element to be 'b'  but was '%v'", result[0])
	}
	if result[2] != "c" {
		t.Errorf("Expected third element to be 'c'  but was '%v'", result[0])
	}
}

func TestTokenizeWithSeparator(t *testing.T) {
	var testString = "a;b;c"
	var result = tokenizeWithSeparator(testString, ";")

	if len(result) != 3 {
		t.Errorf("Expected 3 items but found %v", len(result))
	}

	if result[0] != "a" {
		t.Errorf("Expected first element to be 'a'  but was '%v'", result[0])
	}
	if result[1] != "b" {
		t.Errorf("Expected second element to be 'b'  but was '%v'", result[0])
	}
	if result[2] != "c" {
		t.Errorf("Expected third element to be 'c'  but was '%v'", result[0])
	}
}

func TestContainsXMLTagHasTag(t *testing.T) {
	var testString = "<xml/>"
	var result = containsXMLTag(testString)
	if !result {
		t.Errorf("expected [%v] to contain an xml tag.", testString)
	}
}

func TestContainsXMLTagNoTag(t *testing.T) {
	var testString = "blah blah something xml blah/!@#$%"
	var result = containsXMLTag(testString)
	if result {
		t.Errorf("expected [%v] not to contain an xml tag.", testString)
	}
}

func TestXMLFilter(t *testing.T) {
	var testString = "<?xml version=\"1.0\" encoding=\"utf-8\" ?>"
	var result = filterXMLTags(tokenize(testString))
	if len(result) != 2 {
		t.Errorf("Expected 2 results but found %v", len(result))
	}
	if result[0] != "version=\"1.0\"" {
		t.Errorf("Expected first result to be [version=\"1.0\"]  but was [%v]", result[0])
	}
	if result[1] != "encoding=\"utf-8\"" {
		t.Errorf("Expected first result to be [encoding=\"utf-8\"]  but was [%v]", result[1])
	}
}

func TestGenerateAssignmentsRecursiveDoubleQuote(t *testing.T) {
	var testString = "foo = \"bar\""
	var tree = decisiontree.ConstructTree(tokenize(testString))
	var result = generateAssignmentsRecursive(tree)
	if len(result) != 1 {
		t.Errorf("Expected 1 result in collection but got %v", len(result))
	}
	var assignment = result[0]
	if assignment.Name != "foo" {
		t.Errorf("expected name [foo] but found [%v]", assignment.Name)
	}
	if assignment.Value != "\"bar\"" {
		t.Errorf("expected value [\"bar\"] but found [%v]", assignment.Value)
	}
}

func TestGenerateAssignmentsRecursiveSingleQuote(t *testing.T) {
	var testString = "foo = 'bar'"
	var tree = decisiontree.ConstructTree(tokenize(testString))
	var result = generateAssignmentsRecursive(tree)
	if len(result) != 1 {
		t.Errorf("Expected 1 result in collection but got %v", len(result))
	}
	var assignment = result[0]
	if assignment.Name != "foo" {
		t.Errorf("expected name [foo] but found [%v]", assignment.Name)
	}
	if assignment.Value != "'bar'" {
		t.Errorf("expected value ['bar'] but found [%v]", assignment.Value)
	}
}

func TestGenerateAssignmentsRecursiveNoQuote(t *testing.T) {
	var testString = "foo = bar"
	var tree = decisiontree.ConstructTree(tokenize(testString))
	var result = generateAssignmentsRecursive(tree)
	if len(result) != 1 {
		t.Errorf("Expected 1 result in collection but got %v", len(result))
	}
	var assignment = result[0]
	if assignment.Name != "foo" {
		t.Errorf("expected name [foo] but found [%v]", assignment.Name)
	}
	if assignment.Value != "bar" {
		t.Errorf("expected value [bar] but found [%v]", assignment.Value)
	}
}

func TestGenerateAssignmentsRecursiveMultiAssignment(t *testing.T) {
	var testString = "foo = bar biz : baz fizz := buzz"
	var tree = decisiontree.ConstructTree(tokenize(testString))
	var result = generateAssignmentsRecursive(tree)
	if len(result) != 3 {
		t.Errorf("Expected 3 results in collection but got %v", len(result))
	}
	var assignment1 = result[0]
	if assignment1.Name != "foo" {
		t.Errorf("expected name [foo] but found [%v]", assignment1.Name)
	}
	if assignment1.Value != "bar" {
		t.Errorf("expected value [bar] but found [%v]", assignment1.Value)
	}
	var assignment2 = result[1]
	if assignment2.Name != "biz" {
		t.Errorf("expected name [biz] but found [%v]", assignment2.Name)
	}
	if assignment2.Value != "baz" {
		t.Errorf("expected value [baz] but found [%v]", assignment2.Value)
	}
	var assignment3 = result[2]
	if assignment3.Name != "fizz" {
		t.Errorf("expected name [fizz] but found [%v]", assignment3.Name)
	}
	if assignment3.Value != "buzz" {
		t.Errorf("expected value [buzz] but found [%v]", assignment3.Value)
	}
}

func TestGenerateAssignmentsMultiAssignment(t *testing.T) {
	var testString = "foo = bar biz : baz"
	var result = generateAssignments(testString)
	if len(result) != 2 {
		t.Errorf("Expected 2 results in collection but got %v", len(result))
	}
	var assignment1 = result[0]
	if assignment1.Name != "foo" {
		t.Errorf("expected name [foo] but found [%v]", assignment1.Name)
	}
	if assignment1.Value != "bar" {
		t.Errorf("expected value [bar] but found [%v]", assignment1.Value)
	}
	var assignment2 = result[1]
	if assignment2.Name != "biz" {
		t.Errorf("expected name [biz] but found [%v]", assignment2.Name)
	}
	if assignment2.Value != "baz" {
		t.Errorf("expected value [baz] but found [%v]", assignment2.Value)
	}
}

func TestEvaluateRuleMatch(t *testing.T) {
	teststring := "password = foobar"
	result := evaluateRule(teststring)
	if result.IsEmpty() {
		t.Errorf("Expected 1 result but got an empty result")
	}
	if result.Match != teststring {
		t.Errorf("Expected match to be 'password' but was '%v'", result.Match)
	}
	if result.Confidence != 0.7 ||
		result.Rule != "DecisionTree" {
		t.Errorf("Expected confidence 0.5 and severity 1 but found confidence %v and severity %v", result.Confidence, result.Severity)
	}
}
