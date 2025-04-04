package i18n

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestNewTranslator(t *testing.T) {
	tests := []struct {
		name    string
		lang    string
		wantErr bool
	}{
		{name: "valid language", lang: "en", wantErr: false},
		{name: "valid language farsi", lang: "fa", wantErr: false},
		{name: "invalid language", lang: "xx", wantErr: true},
		{name: "empty language", lang: "", wantErr: true},
		{name: "malformed tag", lang: "e n", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := NewTranslator(tt.lang)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tr)
			}
			if tr != nil {
				tr.Release()
			}
		})
	}
}

func TestTranslator_T(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		msgID    string
		args     []any
		expected string
	}{
		{
			name:     "simple translation",
			lang:     "en",
			msgID:    "welcome",
			expected: "Welcome!",
		},
		{
			name:     "template translation",
			lang:     "en",
			msgID:    "greeting",
			args:     []any{map[string]string{"name": "John"}},
			expected: "Hello, John!",
		},
		{
			name:     "farsi translation",
			lang:     "fa",
			msgID:    "welcome",
			expected: "خوش آمدید!",
		},
		{
			name:     "missing translation",
			lang:     "en",
			msgID:    "nonexistent",
			expected: "nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := NewTranslator(tt.lang)
			assert.NoError(t, err)
			defer tr.Release()

			got := tr.T(tt.msgID, tt.args...)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestConcurrentTranslations(t *testing.T) {
	const (
		numGoroutines = 100
		numIterations = 1000
	)

	langs := []string{"en", "fa"}
	errCh := make(chan error, numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				lang := langs[j%len(langs)]
				tr, err := NewTranslator(lang)
				if err != nil {
					errCh <- fmt.Errorf("goroutine %d: NewTranslator() error = %v", id, err)
					return
				}

				msg := tr.T("welcome")
				if msg == "" {
					errCh <- fmt.Errorf("goroutine %d: empty translation", id)
				}

				tr.Release()
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		assert.NoError(t, err)
	}
}

func TestTranslatorPool(t *testing.T) {
	t.Run("multiple translators", func(t *testing.T) {
		tr1, err := NewTranslator("en")
		assert.NoError(t, err)

		tr2, err := NewTranslator("fa")
		assert.NoError(t, err)

		msg1 := tr1.T("welcome")
		msg2 := tr2.T("welcome")

		assert.NotEqual(t, msg1, msg2, "Translations from different languages should not be equal")

		tr1.Release()
		tr2.Release()
	})

	t.Run("translator reuse", func(t *testing.T) {
		tr1, _ := NewTranslator("en")
		msg1 := tr1.T("welcome")
		tr1.Release()

		tr2, _ := NewTranslator("en")
		defer tr2.Release()

		msg2 := tr2.T("welcome")
		assert.Equal(t, msg1, msg2, "Reused translator should produce same translation")
	})
}

func TestDefaultTranslator(t *testing.T) {
	tests := []struct {
		name     string
		msgID    string
		args     []any
		expected string
	}{
		{
			name:     "simple message",
			msgID:    "welcome",
			expected: "Welcome!",
		},
		{
			name:     "template message",
			msgID:    "greeting",
			args:     []any{map[string]string{"name": "Test"}},
			expected: "Hello, Test!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := T(tt.msgID, tt.args...)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTWithLang(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		msgID    string
		args     []any
		expected string
	}{
		{name: "english", lang: "en", msgID: "welcome", expected: "Welcome!"},
		{name: "farsi", lang: "fa", msgID: "welcome", expected: "خوش آمدید!"},
		{name: "invalid lang", lang: "xx", msgID: "welcome", expected: "welcome"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TWithLang(tt.lang, tt.msgID, tt.args...)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestSupportedLanguages(t *testing.T) {
	expectedLangs := map[language.Tag]bool{
		language.English: true,
		language.Persian: true,
	}

	langs := SupportedLanguages()
	gotLangs := make(map[language.Tag]bool)
	for _, lang := range langs {
		gotLangs[lang] = true
	}

	assert.Equal(t, len(expectedLangs), len(langs), "Wrong number of supported languages")

	for lang := range expectedLangs {
		assert.True(t, gotLangs[lang], "Expected language %v not found in supported languages", lang)
	}
}

func TestFuzzyTranslations(t *testing.T) {
	tests := []struct {
		name    string
		lang    string
		msgID   string
		data    map[string]any
		wantErr bool
	}{
		{
			name:  "valid translation",
			lang:  "en",
			msgID: "welcome",
			data:  map[string]any{"name": "Test"},
		},
		{
			name:    "invalid language",
			lang:    "xx",
			msgID:   "welcome",
			wantErr: true,
		},
		{
			name:    "empty language",
			lang:    "",
			msgID:   "greeting",
			data:    map[string]any{"name": ""},
			wantErr: true,
		},
		{
			name:  "invalid UTF-8 in template data",
			lang:  "en",
			msgID: "greeting",
			data:  map[string]any{"name": string([]byte{0xff, 0xfe, 0xfd})},
		},
		{
			name:  "nonexistent message ID",
			lang:  "en",
			msgID: "nonexistent.key",
		},
		{
			name:  "invalid template key",
			lang:  "en",
			msgID: "greeting",
			data:  map[string]any{"invalid": "value"},
		},
		{
			name:  "null byte in message ID",
			lang:  "en",
			msgID: "welcome\x00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTWithLang(t, tt)
			testTranslatorT(t, tt)
		})
	}
}

func testTWithLang(t *testing.T, tt struct {
	name    string
	lang    string
	msgID   string
	data    map[string]any
	wantErr bool
}) {
	result := TWithLang(tt.lang, tt.msgID, tt.data)
	assert.NotEmpty(t, result, "TWithLang returned empty translation")
	if tt.wantErr {
		assert.Equal(t, tt.msgID, result, "Expected fallback to msgID")
	}
}

func testTranslatorT(t *testing.T, tt struct {
	name    string
	lang    string
	msgID   string
	data    map[string]any
	wantErr bool
}) {
	tr, err := NewTranslator(tt.lang)
	if tt.wantErr {
		assert.Error(t, err)
		return
	}
	assert.NoError(t, err)
	defer tr.Release()

	result := tr.T(tt.msgID, tt.data)
	assert.NotEmpty(t, result, "Translator.T returned empty translation")
	if tt.wantErr {
		assert.Equal(t, tt.msgID, result, "Expected fallback to msgID")
	}
}

func BenchmarkTranslation(b *testing.B) {
	tr, err := NewTranslator("en")
	assert.NoError(b, err)
	defer tr.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.T("welcome")
	}
}

func BenchmarkTranslationWithTemplate(b *testing.B) {
	tr, err := NewTranslator("en")
	assert.NoError(b, err)
	defer tr.Release()

	data := map[string]string{"name": "John"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.T("greeting", data)
	}
}
