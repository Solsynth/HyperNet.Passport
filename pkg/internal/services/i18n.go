package services

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	htmpl "html/template"
	"os"
	"path/filepath"
	"strings"
)

const FallbackLanguage = "en-US"

var LocaleBundle *i18n.Bundle

func LoadLocalization() error {
	LocaleBundle = i18n.NewBundle(language.AmericanEnglish)
	LocaleBundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	var count int

	basePath := viper.GetString("locales_dir")
	if entries, err := os.ReadDir(basePath); err != nil {
		return fmt.Errorf("unable to read locales directory: %v", err)
	} else {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if _, err := LocaleBundle.LoadMessageFile(filepath.Join(basePath, entry.Name())); err != nil {
				return fmt.Errorf("unable to load localization file %s: %v", entry.Name(), err)
			} else {
				count++
			}
		}
	}

	log.Info().Int("locales", count).Msg("Loaded localization files...")

	return nil
}

func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(LocaleBundle, lang)
}

func GetLocalizedString(name string, lang string) string {
	localizer := GetLocalizer(lang)
	msg, err := localizer.LocalizeMessage(&i18n.Message{
		ID: name,
	})
	if err != nil {
		log.Warn().Err(err).Str("lang", lang).Str("name", name).Msg("Failed to localize string...")
		return name
	}
	return msg
}

func GetLocalizedTemplatePath(name string, lang string) string {
	basePath := viper.GetString("templates_dir")
	filePath := filepath.Join(basePath, lang, name)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// Fallback to English
		filePath = filepath.Join(basePath, FallbackLanguage, name)
		return filePath
	}

	return filePath
}

func GetLocalizedTemplateHTML(name string, lang string) *htmpl.Template {
	path := GetLocalizedTemplatePath(name, lang)
	tmpl, err := htmpl.ParseFiles(path)
	if err != nil {
		log.Warn().Err(err).Str("lang", lang).Str("name", name).Msg("Failed to load localized template...")
		return nil
	}

	return tmpl
}

func RenderLocalizedTemplateHTML(name string, lang string, data any) string {
	tmpl := GetLocalizedTemplateHTML(name, lang)
	if tmpl == nil {
		return ""
	}
	buf := new(strings.Builder)
	err := tmpl.Execute(buf, data)
	if err != nil {
		log.Warn().Err(err).Str("lang", lang).Str("name", name).Msg("Failed to render localized template...")
		return ""
	}
	return buf.String()
}
