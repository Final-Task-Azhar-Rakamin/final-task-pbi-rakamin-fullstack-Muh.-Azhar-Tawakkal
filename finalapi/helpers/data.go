package helpers

import (
	"finalapi/database"
	"finalapi/models"
	"fmt"
	"os"
)

func GetId() (int, error) {
	var user models.User
	result := database.DB.Order("id desc").First(&user)
	latestID := int(user.Id) + 1
	return latestID, result.Error
}

func DeletePhoto(photoType string, filename string) {
	var filePath string
	if photoType == "profile" {
		if filename == "default.jpg" {
			return
		}
		filePath = "./assets/profiles/" + filename
	} else if photoType == "post" {
		filePath = "./assets/posts/" + filename
	}
	if err := os.Remove(filePath); err == nil {
		fmt.Println("Delete file success")
	} else {
		fmt.Println("Delete file failed, delete manually")
	}

}

func CheckFile(ext string) bool {
	imageExts := []string{"jpg", "jpeg", "png", "gif", "bmp"}
	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}
