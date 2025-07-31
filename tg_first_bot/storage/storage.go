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
	"strconv"
)

const (
	perm = 0755
)

func (b InternalBasePath) Save(p Page) error {
	fPath := filepath.Join(b.BasePath, strconv.Itoa(p.ChatID))
	if err := os.MkdirAll(fPath, perm); err != nil {
		return errors.WrapIfErr("failed MkdirAll", err)
	}
	fName, _ := hash(p)
	if fName == "" {
		return errors2.New("fName is empty")
	}
	fullPath := filepath.Join(fPath, fName)
	file, err := os.Create(fullPath)
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

func (b InternalBasePath) PickRandom(ChatID int) (*Page, error) {
	dirPath := filepath.Join(b.BasePath, strconv.Itoa(ChatID))
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("failed ReadDir %x", err)
		return nil, nil
	}
	if len(files) == 0 {
		return nil, nil
	}

	n := len(files)
	file := files[rand.Intn(n)]
	fPath := filepath.Join(dirPath, file.Name())

	f, err := os.ReadFile(fPath)
	if err != nil {
		log.Printf("failed ReadFile %x", err)
		return nil, nil
	}
	log.Printf("Метод PickRandom, готовим к отправке TextHref: %s", string(f))
	p := Page{
		ChatID:   ChatID,
		TextPage: string(f),
	}

	return &p, nil
}

func (b InternalBasePath) IsExist(p *Page) (bool, error) {
	dirPath := filepath.Join(b.BasePath, strconv.Itoa(p.ChatID))
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("failed read at directory %s:, %w", dirPath, err)
		return false, nil
	}
	if len(files) == 0 {
		log.Println("files were not found in this directory")
		return false, nil
	}

	for _, file := range files {
		fPath := filepath.Join(dirPath, file.Name())
		contentFile, err := os.ReadFile(fPath)
		if string(contentFile) == "" {
			return false, errors2.New("contentFile is empty")
		}
		if err != nil {
			return false, fmt.Errorf("failed os.ReadFile %w", err)
		}
		if string(contentFile) == p.TextPage {
			return true, nil
		}
		continue
	}

	return false, nil
}

func (b InternalBasePath) Remove(p *Page) error {
	dirPath := filepath.Join(b.BasePath, strconv.Itoa(p.ChatID))
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return errors.WrapIfErr("failed ReadDir", err)
	}
	if len(files) == 0 {
		log.Println("files were not found in this directory")
		return nil
	}

	for _, file := range files {
		fPath := filepath.Join(dirPath, file.Name())
		contentFile, err := os.ReadFile(fPath)
		if string(contentFile) == "" {
			return errors2.New("contentFile is empty")
		}
		if err != nil {
			return fmt.Errorf("failed os.ReadFile %w", err)
		}
		if string(contentFile) == p.TextPage {
			if err = os.Remove(fPath); err != nil {
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
