package storage

import (
	"crypto/sha1"
	"encoding/gob"
	errors2 "errors"
	"example.com/errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	perm = 0644
)

func (b InternalBasePath) Save(p Page) error {
	fName, _ := hash(p)
	if fName == "" {
		return errors2.New("fName is empty")
	}
	if err := os.MkdirAll(b.BasePath, perm); err != nil {
		return errors.WrapIfErr("failed MkdirAll", err)
	}
	fPath := filepath.Join(b.BasePath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return errors.WrapIfErr("failed create file", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("failed file closing %v", err)
		}
	}()
	if _, err := file.WriteString(p.TextPage); err != nil {
		return fmt.Errorf("failed WriteString() in file %x", err)
	}
	return nil
}

func (b InternalBasePath) PickRandom() (string, error) {
	files, err := os.ReadDir(b.BasePath)
	if err != nil {
		return "", errors.WrapIfErr("failed ReadDir", err)
	}
	if len(files) == 0 {
		return "", nil
	}

	n := len(files)
	file := files[rand.Intn(n)]
	fPath := filepath.Join(b.BasePath, file.Name())

	f, err := os.ReadFile(fPath)
	if err != nil {
		return "", fmt.Errorf("failed ReadFile %x", err)
	}
	log.Printf("Метод PickRandom, готовим к отправке TextHref: %s", string(f))

	return string(f), nil
}

func (b InternalBasePath) IsExist(p Page) (bool, error) {
	fName, _ := hash(p)
	if fName == "" {
		return false, errors2.New("fName is empty")
	}
	fPath := filepath.Join(b.BasePath, fName)
	_, err := os.Stat(fPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.WrapIfErr("failed to check path", err)
	}
	return true, nil
}

func (b InternalBasePath) Remove(contentPage string) error {
	files, err := os.ReadDir(b.BasePath)
	if err != nil {
		return errors.WrapIfErr("failed ReadDir", err)
	}
	if len(files) == 0 {
		log.Println("files were not found in this directory")
		return nil
	}

	for _, file := range files {
		fPath := filepath.Join(b.BasePath, file.Name())
		contentFile, err := os.ReadFile(fPath)
		if string(contentFile) == "" {
			return errors2.New("fName is empty")
		}
		if err != nil {
			return fmt.Errorf("failed os.ReadFile %x", err)
		}
		if string(contentFile) == contentPage {
			if err := os.Remove(fPath); err != nil {
				return errors.WrapIfErr("couldn't delete the file", err)
			}
		}
	}
	return nil
}

//func (b InternalBasePath) decodePage(f io.Reader) (*Page, error) {
//	var p Page
//	if err := gob.NewDecoder(f).Decode(&p); err != nil {
//		return nil, fmt.Errorf("failed decode %s", p.TextPage)
//	}
//	return &p, nil
//}

func hash(p Page) (string, error) {
	h := sha1.New()
	if err := gob.NewEncoder(h).Encode(p.TextPage); err != nil {
		return "", fmt.Errorf("unsuccessful text hashing %s", p.TextPage)
	}
	if err := gob.NewEncoder(h).Encode(p.ChatID); err != nil {
		return "", fmt.Errorf("unsuccessful username hashing %v", p.ChatID)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
