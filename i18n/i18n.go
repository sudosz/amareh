// Package i18n provides internationalization and localization functionality.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var locales embed.FS

var (
	bundle            *goi18n.Bundle
	defaultLang       = language.English
	defaultTranslator *Translator
	translatorPool    sync.Pool
	supportedLangs    = make(map[language.Tag]bool)
	initOnce          sync.Once
)

func init() {
	initOnce.Do(func() {
		bundle = goi18n.NewBundle(defaultLang)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		// Load all locale files
		entries, err := locales.ReadDir("locales")
		if err != nil {
			panic(fmt.Errorf("failed to read locales directory: %w", err))
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			localeData, err := locales.ReadFile("locales/" + entry.Name())
			if err != nil {
				panic(fmt.Errorf("failed to read locale file %s: %w", entry.Name(), err))
			}

			if _, err := bundle.ParseMessageFileBytes(localeData, entry.Name()); err != nil {
				panic(fmt.Errorf("failed to parse locale file %s: %w", entry.Name(), err))
			}

			// Extract language tag from filename and add to supported languages
			if tag, err := language.Parse(entry.Name()[:2]); err == nil {
				supportedLangs[tag] = true
			}
		}

		// Initialize translator pool
		translatorPool = sync.Pool{
			New: func() interface{} {
				return &Translator{}
			},
		}

		// Initialize default translator
		defaultTranslator, err = NewTranslator(defaultLang.String())
		if err != nil {
			panic(fmt.Errorf("failed to create default translator: %w", err))
		}
	})
}

// Translator is a struct that contains a localizer and a mutex
type Translator struct {
	localizer *goi18n.Localizer
	mu        sync.RWMutex
}

// NewTranslator creates a new translator for the given language
func NewTranslator(lang string) (*Translator, error) {
	tag, err := language.Parse(lang)
	if err != nil {
		return nil, fmt.Errorf("invalid language tag %q: %w", lang, err)
	}

	if !supportedLangs[tag] {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	t := translatorPool.Get().(*Translator)
	t.localizer = goi18n.NewLocalizer(bundle, lang)
	return t, nil
}

// T returns the translation for the given message ID and arguments
func (t *Translator) T(msgID string, args ...any) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	cfg := &goi18n.LocalizeConfig{
		MessageID:      msgID,
		DefaultMessage: &goi18n.Message{ID: msgID},
	}

	switch len(args) {
	case 1:
		cfg.TemplateData = args[0]
	case 2:
		cfg.TemplateData = args[0]
		cfg.PluralCount = args[1]
	}

	msg, err := t.localizer.Localize(cfg)
	if err != nil {
		return msgID // Fallback to message ID if translation fails
	}
	return msg
}

// Release releases the translator back to the pool
func (t *Translator) Release() {
	t.mu.Lock()
	t.localizer = nil
	t.mu.Unlock()
	translatorPool.Put(t)
}

// T is a convenience function that returns the translation for the default language
func T(msgID string, args ...any) string {
	return defaultTranslator.T(msgID, args...)
}

// TWithLang is a convenience function that returns the translation for a specific language
func TWithLang(lang string, msgID string, args ...any) string {
	t, err := NewTranslator(lang)
	if err != nil {
		return msgID // Fallback to message ID if translator creation fails
	}
	defer t.Release()
	return t.T(msgID, args...)
}

// SupportedLanguages returns a slice of all supported languages
func SupportedLanguages() []language.Tag {
	langs := make([]language.Tag, 0, len(supportedLangs))
	for lang := range supportedLangs {
		langs = append(langs, lang)
	}
	return langs
}
