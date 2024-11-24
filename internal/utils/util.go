package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func ListFiles(root, ext string) ([]string, error) {
	var filenames []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ext {
			filenames = append(filenames, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return filenames, nil
}

func GetCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get current working directory")
	}
	return cwd
}

func ToJson(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func Merge2Maps(m1 map[string]string, m2 map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}

func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func RemoveExt(filename string, ext string) string {
	if strings.HasSuffix(filename, ext) {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}

func RemoveExtFromList(filenames []string, ext string) []string {
	var result []string
	for _, filename := range filenames {
		result = append(result, RemoveExt(filename, ext))
	}
	return result
}

func ReadFile(file string) string {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v\n", err)
	}
	return string(bytes)
}

func IsStructEmpty(s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface() {
			return false
		}
	}
	return true
}
