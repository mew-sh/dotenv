package dotenv

import (
	"strings"
	"testing"
)

var (
	smallEnvContent = `KEY1=value1
KEY2=value2
KEY3=value3`

	mediumEnvContent = strings.Repeat(`KEY=value
QUOTED_KEY="quoted value"
EXPANDED_KEY=${KEY}_expanded
`, 100)

	largeEnvContent = strings.Repeat(`KEY=value
QUOTED_KEY="quoted value with spaces and special chars !@#$%^&*()"
EXPANDED_KEY=${KEY}_expanded
MULTILINE_KEY="line1\nline2\nline3"
COMMENT_KEY=value # this is a comment
export EXPORTED_KEY=exported_value
`, 1000)

	complexEnvContent = `# Complex .env file for benchmarking
DATABASE_URL="postgres://user:pass@localhost:5432/db"
REDIS_URL=redis://localhost:6379
API_KEY=super-secret-api-key-123456789
DEBUG=true
PORT=8080

# Nested variable expansion
BASE_URL=https://api.example.com
API_V1=${BASE_URL}/v1
API_V2=${BASE_URL}/v2
USERS_ENDPOINT=${API_V1}/users
POSTS_ENDPOINT=${API_V1}/posts

# Quoted values with escaping
MESSAGE="Hello \"World\"\nWelcome to our app!"
MULTILINE="Line 1\nLine 2\nLine 3"
TAB_SEPARATED="Col1\tCol2\tCol3"

# Feature flags
ENABLE_FEATURE_A=yes
ENABLE_FEATURE_B=no
ENABLE_FEATURE_C=true
ENABLE_FEATURE_D=false

# Export syntax
export NODE_ENV=development
export PATH="/usr/local/bin:${PATH}"

# YAML style
smtp_host: smtp.example.com
smtp_port: 587
smtp_user: user@example.com`
)

func BenchmarkParseSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(smallEnvContent))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMedium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(mediumEnvContent))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(largeEnvContent))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(complexEnvContent))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(mediumEnvContent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	env := map[string]string{
		"KEY1":    "value1",
		"KEY2":    "value with spaces",
		"KEY3":    "value\nwith\nnewlines",
		"KEY4":    "value\"with\"quotes",
		"KEY5":    "",
		"NUMERIC": "12345",
		"BOOLEAN": "true",
		"URL":     "https://example.com/path?param=value",
		"PATH":    "/usr/local/bin:/usr/bin:/bin",
		"COMPLEX": "value with\ttabs and\nspecial \"chars\"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(env)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVariableExpansion(b *testing.B) {
	content := `BASE=hello
LEVEL1=${BASE}_world
LEVEL2=${LEVEL1}_again
LEVEL3=${LEVEL2}_and_again
FINAL=${LEVEL3}_final`

	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNoVariableExpansion(b *testing.B) {
	content := `BASE=hello
LEVEL1=${BASE}_world
LEVEL2=${LEVEL1}_again
LEVEL3=${LEVEL2}_and_again
FINAL=${LEVEL3}_final`

	parser := NewParserWithOptions(false)

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(strings.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEscapeProcessing(b *testing.B) {
	content := `ESCAPED="value with\nnewlines\tand\ttabs\"and quotes\""
ANOTHER="more\nescaping\rhere"`

	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCommentRemoval(b *testing.B) {
	content := `# This is a comment
KEY1=value1 # inline comment
KEY2="quoted value" # another comment
# Another full line comment
KEY3=value3`

	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark memory allocations
func BenchmarkParseAllocs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Parse(strings.NewReader(mediumEnvContent))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalAllocs(b *testing.B) {
	env := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(env)
		if err != nil {
			b.Fatal(err)
		}
	}
}
