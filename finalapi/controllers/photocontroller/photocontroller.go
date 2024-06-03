package photocontroller

import (
	"finalapi/database"
	"finalapi/helpers"
	"finalapi/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var photos []models.Photo
	id, _ := c.Get("id")

	database.DB.Where("user_id = ?", id).Find(&photos)
	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

func Show(c *gin.Context) {
	var photo models.Photo
	var checkPhoto models.Photo
	id := c.Param("id")

	// Validate User
	database.DB.Model(&checkPhoto).Where("id = ?", id).First(&checkPhoto)
	if valid := helpers.ValidateOwner(c, checkPhoto); !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access unauthorized"})
		return
	}

	if err := database.DB.Where("id = ?", id).First(&photo).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"photo": photo})
}

func Create(c *gin.Context) {
	var photo models.Photo
	var user models.User
	id, _ := c.Get("id")
	file, _ := c.FormFile("photo_file")

	// Binding data
	if err := c.ShouldBind(&photo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Manage photo data
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	photo.Id = uuid.New().String()
	photo.UserID = user.Id
	photo.User = user

	if file != nil {
		typeFile := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
		photo.PhotoUrl = strconv.Itoa(int(photo.UserID)) + "_posts_" + photo.Id + "." + typeFile
		if valid := helpers.CheckFile(typeFile); !valid && typeFile != "" {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "File is not an image"})
			return
		}
	}

	// Validate data
	if valid, err := govalidator.ValidateStruct(photo); !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manage file
	if file != nil {
		savePath := "./assets/posts/" + photo.PhotoUrl
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
			return
		}
	}

	database.DB.Create(&photo)
	c.JSON(http.StatusOK, gin.H{"message": "Input data success", "photo Id": photo.Id})
}

func Update(c *gin.Context) {
	var UpdatePhoto models.UpdatePhoto
	var photo models.Photo
	id := c.Param("id")

	// Validate User
	database.DB.Where("id = ?", id).First(&photo)
	if valid := helpers.ValidateOwner(c, photo); !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access unauthorized"})
		return
	}

	// Binding data
	if err := c.ShouldBindBodyWithJSON(&UpdatePhoto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Update data
	results := database.DB.Model(&photo).Where("id = ?", id).Updates(&UpdatePhoto)
	if results.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Update not affect anything"})
		return
	} else if results.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update data success"})
}

func Delete(c *gin.Context) {
	var photo models.Photo
	var checkPhoto models.Photo
	id := c.Param("id")

	// Validate User
	database.DB.Where("id = ?", id).First(&checkPhoto)
	if valid := helpers.ValidateOwner(c, checkPhoto); !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access unauthorized"})
		return
	}

	helpers.DeletePhoto("post", checkPhoto.PhotoUrl)

	if database.DB.Delete(&photo, "id = ?", id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete data success"})

}
