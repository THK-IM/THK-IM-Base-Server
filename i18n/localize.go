package i18n

import (
	"encoding/json"
	"os"
	"path"
	"sort"
	"strings"
)

type Localize interface {
	Get(key string, languageTag string) string
	SetDefaultLanguage(value string)
	GetSupportedLanguages() []string
}

type defaultLocalize struct {
	translations map[string]map[string]string
	defaultTag   string
}

func (d *defaultLocalize) SetDefaultLanguage(value string) {
	d.defaultTag = value
}

func (d *defaultLocalize) getTranslationMap(languageTag string) map[string]string {
	var translationMap map[string]string = nil
	tags := d.GetSupportedLanguages()
	for _, tag := range tags {
		if tag == languageTag {
			translationMap = d.translations[tag]
			break
		} else if strings.Contains(languageTag, tag) {
			translationMap = d.translations[tag]
		}
	}
	if translationMap == nil {
		return d.translations[d.defaultTag]
	}
	return translationMap
}

func (d *defaultLocalize) Get(key string, languageTag string) string {
	translationMap := d.getTranslationMap(languageTag)
	if translationMap[key] == "" {
		return key
	} else {
		return translationMap[key]
	}
}

func (d *defaultLocalize) GetSupportedLanguages() []string {
	tags := make([]string, 0)
	for k, _ := range d.translations {
		tags = append(tags, k)
	}
	sort.Slice(tags, func(i, j int) bool {
		if len(tags[i]) < len(tags[j]) {
			return true
		}
		return false
	})
	return tags
}

func NewLocalize(assetsPath string) Localize {
	languageDirPath := path.Join(assetsPath, "languages")
	languagesPath, err := os.ReadDir(languageDirPath)
	translations := make(map[string]map[string]string)
	if err == nil {
		for _, file := range languagesPath {
			if !file.IsDir() {
				languageTag := strings.Replace(file.Name(), ".json", "", -1)
				filePath := path.Join(languageDirPath, file.Name())
				bytes, errRead := os.ReadFile(filePath)
				if errRead != nil {
					panic(errRead)
				} else {
					jsonMap := make(map[string]string)
					errJson := json.Unmarshal(bytes, &jsonMap)
					if errJson != nil {
						panic(errJson)
					}
					translations[languageTag] = jsonMap
				}
			}
		}
	}
	return &defaultLocalize{translations: translations, defaultTag: "en"}
}
