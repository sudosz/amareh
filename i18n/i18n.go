package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"slices"
)

//go:embed locales/*.json
var locales embed.FS

// Translator holds the translation bundle and the localizers.
type Translator struct {
	bundle             *goi18n.Bundle
	supportedLanguages []language.Tag
	current            *goi18n.Localizer
	defaultLocalizer   *goi18n.Localizer
	mutex             sync.RWMutex
}

// NewTranslator creates a new Translator, loads translation files, and sets English as the default.
func NewTranslator() (*Translator, error) {
	bundle := goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	t := &Translator{
		bundle: bundle,
	}

	if err := t.loadTranslations(); err != nil {
		return nil, fmt.Errorf("loading translations failed: %w", err)
	}

	t.defaultLocalizer = goi18n.NewLocalizer(bundle, "en")
	t.current = t.defaultLocalizer

	return t, nil
}

// loadTranslations reads all JSON translation files from the locales directory.
func (t *Translator) loadTranslations() error {
	files, err := locales.ReadDir("locales")
	if err != nil {
		return fmt.Errorf("reading locales directory failed: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || path.Ext(file.Name()) != ".json" {
			continue
		}

		lang := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		tag, err := language.Parse(lang)
		if err != nil {
			return fmt.Errorf("parsing language tag for %s failed: %w", file.Name(), err)
		}

		filePath := path.Join("locales", file.Name())
		if _, err := t.bundle.LoadMessageFileFS(locales, filePath); err != nil {
			return fmt.Errorf("loading file %s failed: %w", file.Name(), err)
		}

		t.supportedLanguages = append(t.supportedLanguages, tag)
	}

	return nil
}

// SetLanguage changes the active language if it is supported.
func (t *Translator) SetLanguage(lang string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return fmt.Errorf("invalid language %q: %w", lang, err)
	}

	if !t.isLanguageSupported(tag) {
		return fmt.Errorf("language %q not supported", lang)
	}

	t.mutex.Lock()
	t.current = goi18n.NewLocalizer(t.bundle, lang)
	t.mutex.Unlock()
	return nil
}

// isLanguageSupported checks if the provided language tag exists in the supported list.
func (t *Translator) isLanguageSupported(tag language.Tag) bool {
	return slices.Contains(t.supportedLanguages, tag)
}

// SupportedLang represents a supported language with its tag and display name.
type SupportedLang struct {
	Tag         language.Tag
	DisplayName string
}

// GetSupportedLanguages returns a slice of all supported languages.
func (t *Translator) GetSupportedLanguages() []SupportedLang {
	languages := make([]SupportedLang, 0, len(t.supportedLanguages))
	for _, tag := range t.supportedLanguages {
		languages = append(languages, SupportedLang{
			Tag:         tag,
			DisplayName: tag.String(),
		})
	}
	return languages
}

// T translates a message using the current language.
func (t *Translator) T(messageID string, args ...any) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.translate(t.current, messageID, args...)
}

// TDefault translates a message using the default language.
func (t *Translator) TDefault(messageID string, args ...any) string {
	return t.translate(t.defaultLocalizer, messageID, args...)
}

// translate is a helper function to perform message translation.
func (t *Translator) translate(localizer *goi18n.Localizer, messageID string, args ...any) string {
	cfg := &goi18n.LocalizeConfig{
		MessageID: messageID,
	}

	switch len(args) {
	case 1:
		cfg.TemplateData = args[0]
	case 2:
		cfg.TemplateData = args[0]
		cfg.PluralCount = args[1]
	}

	return localizer.MustLocalize(cfg)
}

var (
	defaultTranslator *Translator
	once             sync.Once
)

// init initializes the global translator safely.
func init() {
	once.Do(func() {
		var err error
		defaultTranslator, err = NewTranslator()
		if err != nil {
			panic(fmt.Sprintf("initializing translator failed: %v", err))
		}
	})
}

// T is a convenience function to translate a message using the current language.
func T(messageID string, args ...any) string {
	return defaultTranslator.T(messageID, args...)
}

// TDefault translates a message using the default language.
func TDefault(messageID string, args ...any) string {
	return defaultTranslator.TDefault(messageID, args...)
}

// SetLanguage globally changes the active translation language.
func SetLanguage(lang string) error {
	return defaultTranslator.SetLanguage(lang)
}

// GetSupportedLanguages returns all supported languages from the global translator.
func GetSupportedLanguages() []SupportedLang {
	return defaultTranslator.GetSupportedLanguages()
}
