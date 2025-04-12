package model

import (
	"time"

	"gorm.io/gorm"
)

type DefectiveImage struct {
	ID        int       `gorm:"primaryKey"`
	ImageData []byte    `gorm:"type:bytea"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (DefectiveImage) TableName() string {
	return "defective_images"
}

func GetImages(db *gorm.DB) ([]DefectiveImage, error) {
	var images []DefectiveImage
	if err := db.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func CreateImage(db *gorm.DB, imageData []byte) (*DefectiveImage, error) {
	image := &DefectiveImage{
		ImageData: imageData,
	}
	if err := db.Create(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

func DeleteImage(db *gorm.DB, id int) error {
	return db.Delete(&DefectiveImage{}, id).Error
} 