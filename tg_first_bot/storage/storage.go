package storage

import (
	"crypto/sha1"
	"encoding/gob"
	errors2 "errors"
	"example.com/errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	perm = 0644
)

func (b InternalBasePath) Save(p Page) error { // как задавать InternalBasePath?
	fName, _ := hash(p)
	if fName == "" {
		return errors2.New("fName is empty")
	}
	fPath := filepath.Join(b.BasePath, fName)
	file, err := os.OpenFile(fPath, os.O_CREATE, perm)
	if err != nil {
		return errors.WrapIfErr("failed create file", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("failed file closing %v", err)
		}
	}()
	return nil
}

func (b InternalBasePath) PickRandom() (*Page, error) { // почему не передаем chatID
	files, err := os.ReadDir(b.BasePath)
	if err != nil {
		return nil, errors.WrapIfErr("failed ReadDir", err)
	}
	if len(files) == 0 {
		return nil, errors2.New("files were not found in this directory")
	}

	n := len(files)
	file := files[rand.Intn(n)]
	fName := filepath.Join(b.BasePath, file.Name())

	f, err := os.Open(fName)
	if err != nil {
		return nil, errors.WrapIfErr("failed open file", err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("failed closing open file %s", fName)
		}
	}()

	p, _ := b.decodePage(f)
	return p, nil
}

func (b InternalBasePath) IsExist(p Page) (bool, error) {
	fName, _ := hash(p)
	if fName == "" {
		return false, errors2.New("fName is empty")
	}
	fPath := filepath.Join(b.BasePath, fName)
	fInfo, err := os.Stat(fPath)
	if err != nil && !os.IsNotExist(err) {
		return false, errors.WrapIfErr("failed to check path", err)
	}
	if fInfo.IsDir() == true {
		return true, nil
	}

	return false, nil
}

func (b InternalBasePath) Remove(p *Page) error {
	fName, _ := hash(*p)
	if fName == "" {
		return errors2.New("fName is empty")
	}
	fPath := filepath.Join(b.BasePath, fName)
	if err := os.Remove(fPath); err != nil {
		return errors.WrapIfErr("couldn't delete the file", err)
	}

	return nil
}

func (b InternalBasePath) decodePage(f io.Reader) (*Page, error) {
	var p Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("failed decode %s", p.TextPage)
	}
	return &p, nil
}

func hash(p Page) (string, error) {
	h := sha1.New()
	if err := gob.NewEncoder(h).Encode(p.TextPage); err != nil {
		return "", fmt.Errorf("unsuccessful text hashing %s", p.TextPage)
	}
	if err := gob.NewEncoder(h).Encode(p.ChatID); err != nil {
		return "", fmt.Errorf("unsuccessful username hashing %v", p.ChatID)
	}
	return string(h.Sum(nil)), nil
}
