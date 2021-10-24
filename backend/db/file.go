package db

import (
	"fmt"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type File struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"non null" json:"name"`
	Path      string `gorm:"non null" json:"path"`
	Owner     int    `gorm:"non null" json:"userID"`
	ChatID    int64  `gorm:"non null" json:"chatID"`
	MessageID int    `gorm:"non null" json:"messageID"`
	URL       string `gorm:"not null" json:"-"`
}

func (f *File) isUnique() (bool, error) {
	var c int64
	err := GetDB().Model(&File{}).Where("path = ? AND owner = ?", f.Path, f.Owner).Count(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}

	return c == 1 || c == 0, nil
}

func GetFileByID(id uint64) (*File, error) {
	f := File{}
	return &f, GetDB().Take(&f, id).Error
}

func GetAllFilesOwnedBy(userID int) ([]File, error) {
	f := []File{}
	return f, GetDB().Where("owner = ?", userID).Find(&f).Error
}

func PutFile(f *File) error {
	invalidPath := fmt.Errorf("Invalid path")
	alreadyExists := fmt.Errorf("The file already exists")
	if f.Path == "" {
		return invalidPath
	}

	f.Path = strings.ReplaceAll(f.Path, "\\", "/")
	if !strings.HasPrefix(f.Path, "/") {
		f.Path = "/" + f.Path
	}

	{
		new := ""
		for {
			new = strings.ReplaceAll(f.Path, "//", "/")
			if f.Path == new {
				break
			}
			f.Path = new
		}
	}

	if strings.Contains(f.Path, "/..") || strings.Contains(f.Path, "/.") {
		return invalidPath
	}

	f.Name = filepath.Base(f.Path)

	unique, err := f.isUnique()
	if err != nil {
		return err
	}

	if !unique {
		return alreadyExists
	}

	return GetDB().Create(f).Error
}
