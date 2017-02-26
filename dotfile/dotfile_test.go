package dotfile

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func parseAndCompare(t *testing.T, rawEnvLine string, expectedKey string, expectedValue string) {
	var configItem ConfigItem
	parseLine(&configItem, rawEnvLine)
	if configItem.Key != expectedKey || configItem.Value != expectedValue {
		t.Errorf("Expected '%v' to parse as '%v' => '%v', got '%v' => '%v' instead", rawEnvLine, expectedKey, expectedValue, configItem.Key, configItem.Value)
	}
}

func TestReadEnvFile(t *testing.T) {
	filename, _ := filepath.Abs("../fixtrues/plain.env")
	if kvMap, err := ReadEnvFile(filename); err == nil {
		assert.Equal(t, kvMap["DB"].Value, "db")
		assert.Equal(t, kvMap["DB"].Key, "DB")
		assert.Equal(t, kvMap["DB"].Comment, "dbCommentLine")

		assert.Equal(t, kvMap["JWT_SECRET"].Value, "jwt_secret")
		assert.Equal(t, kvMap["JWT_SECRET"].Key, "JWT_SECRET")
		assert.Equal(t, kvMap["JWT_SECRET"].Comment, "")

		assert.Equal(t, kvMap["JWT_TIMEOUT"].Value, 1024)

		assert.Equal(t, kvMap["APP_KEY"].Value, "app_key")
		assert.Equal(t, kvMap["APP_KEY"].Key, "APP_KEY")
		assert.Equal(t, kvMap["APP_KEY"].Comment, "appkeyCommentLine")

		assert.Equal(t, kvMap["DOMAIN"].Value, "")
	} else {
		t.Error(err.Error())
	}
}

func TestParsing(t *testing.T) {
	// unquoted values
	parseAndCompare(t, "FOO=bar", "FOO", "bar")

	// parses values with spaces around equal sign
	parseAndCompare(t, "FOO =bar", "FOO", "bar")
	parseAndCompare(t, "FOO= bar", "FOO", "bar")

	// parses double quoted values
	parseAndCompare(t, "FOO=\"bar\"", "FOO", "bar")

	// parses single quoted values
	parseAndCompare(t, "FOO='bar'", "FOO", "bar")

	// parses escaped double quotes
	parseAndCompare(t, "FOO=escaped\\\"bar\"", "FOO", "escaped\"bar")

	// parses yaml style options
	parseAndCompare(t, "OPTION_A: 1", "OPTION_A", "1")

	// parses export keyword
	parseAndCompare(t, "export OPTION_A=2", "OPTION_A", "2")
	parseAndCompare(t, "export OPTION_B='\\n'", "OPTION_B", "\n")

	// it 'expands newlines in quoted strings' do
	// expect(env('FOO="bar\nbaz"')).to eql('FOO' => "bar\nbaz")
	parseAndCompare(t, "FOO=\"bar\\nbaz\"", "FOO", "bar\nbaz")

	// it 'parses varibales with "." in the name' do
	// expect(env('FOO.BAR=foobar')).to eql('FOO.BAR' => 'foobar')
	parseAndCompare(t, "FOO.BAR=foobar", "FOO.BAR", "foobar")

	// it 'parses varibales with several "=" in the value' do
	// expect(env('FOO=foobar=')).to eql('FOO' => 'foobar=')
	parseAndCompare(t, "FOO=foobar=", "FOO", "foobar=")

	// it 'strips unquoted values' do
	// expect(env('foo=bar ')).to eql('foo' => 'bar') # not 'bar '
	parseAndCompare(t, "FOO=bar ", "FOO", "bar")

	// it 'ignores inline comments' do
	// expect(env("foo=bar # this is foo")).to eql('foo' => 'bar')
	parseAndCompare(t, "FOO=bar # this is foo", "FOO", "bar")

	// it 'allows # in quoted value' do
	// expect(env('foo="bar#baz" # comment')).to eql('foo' => 'bar#baz')
	parseAndCompare(t, "FOO=\"bar#baz\" # comment", "FOO", "bar#baz")
	parseAndCompare(t, "FOO='bar#baz' # comment", "FOO", "bar#baz")
	parseAndCompare(t, "FOO=\"bar#baz#bang\" # comment", "FOO", "bar#baz#bang")

	// it 'parses # in quoted values' do
	// expect(env('foo="ba#r"')).to eql('foo' => 'ba#r')
	// expect(env("foo='ba#r'")).to eql('foo' => 'ba#r')
	parseAndCompare(t, "FOO=\"ba#r\"", "FOO", "ba#r")
	parseAndCompare(t, "FOO='ba#r'", "FOO", "ba#r")

	// it 'throws an error if line format is incorrect' do
	// expect{env('lol$wut')}.to raise_error(Dotenv::FormatError)
	badlyFormattedLine := "lol$wut"
	var configItem ConfigItem
	err := parseLine(&configItem, badlyFormattedLine)
	if err == nil {
		t.Errorf("Expected \"%v\" to return error, but it didn't", badlyFormattedLine)
	}
}

func TestIsCommentLine(t *testing.T) {
	if isCommentLine("") {
		t.Error("Line with nothing so not a comment line")
	}

	if !isCommentLine("# ato a") {
		t.Error("Line should be a comment line")
	}

	if !isCommentLine("\t# ato a") {
		t.Error(`\t Line should be a comment line`)
	}

	if isCommentLine("\n") {
		t.Error(`\n should not be a comment line`)
	}

}

func TestIsIgnoredLine(t *testing.T) {
	if !isIgnoredLine("\n") {
		t.Error("Line with nothing but line break wasn't ignored")
	}

	if !isIgnoredLine("\t\t ") {
		t.Error("Line full of whitespace wasn't ignored")
	}

	// it 'ignores comment lines' do
	// expect(env("\n\n\n # HERE GOES FOO \nfoo=bar")).to eql('foo' => 'bar')
	if isIgnoredLine("# comment") {
		t.Error("Comment should not be ignored")
	}

	if isIgnoredLine("\t#comment") {
		t.Error("Indented comment wasn't ignored")
	}

	// make sure we're not getting false positives
	if isIgnoredLine("export OPTION_B='\\n'") {
		t.Error("ignoring a perfectly valid line to parse")
	}

}
