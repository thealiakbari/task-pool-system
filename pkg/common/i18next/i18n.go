package i18next

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
)

var (
	languages   map[language.Tag]*i18n.Localizer
	defaultLang language.Tag
)

func NewLanguage(defaultLng language.Tag) error {
	languages = make(map[language.Tag]*i18n.Localizer)
	defaultLang = defaultLng
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if err := filepath.Walk("assets/locales", func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, "toml") {
			return nil
		}
		if _, err = bundle.LoadMessageFile(path); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.New("can not read locale files")
	}

	for _, tag := range bundle.LanguageTags() {
		languages[tag] = i18n.NewLocalizer(bundle, tag.String())
	}

	return nil
}

func ByLangWithData(lang language.Tag, id string, data interface{}) string {
	local, ok := languages[lang]
	if !ok {
		if local, ok = languages[defaultLang]; !ok {
			return id
		}
	}
	str, err := local.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		return id
	}

	return str
}

func ByLang(lang language.Tag, id string) string {
	return ByLangWithData(lang, id, nil)
}

func ByContextWithData(ctx context.Context, id string, data interface{}) string {
	lang, _ := GetLang(ctx)
	return ByLangWithData(lang, id, data)
}

func ByContext(ctx context.Context, id string) string {
	return ByContextWithData(ctx, id, nil)
}

func GetLang(ctx context.Context) (language.Tag, bool) {
	value, ok := GetValue(ctx, "lang")
	if !ok {
		return language.Und, false
	}
	tag, err := language.Parse(value)
	if err != nil {
		return language.Und, false
	}
	return tag, true
}

func GetValue(ctx context.Context, key string) (string, bool) {
	value, ok := ctx.Value(key).(string)
	if ok {
		return value, true
	}
	var md metadata.MD
	md, ok = metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	if len(md[key]) > 0 {
		return md[key][0], true
	}

	return "", false
}
