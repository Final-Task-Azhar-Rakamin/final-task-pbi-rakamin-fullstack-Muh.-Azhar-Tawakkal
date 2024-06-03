package usercontroller

import (
	"finalapi/database"
	"finalapi/helpers"
	"finalapi/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginUser models.LoginUser
	var userData models.User

	// Binding data ----------------------------------------------------------------
	if err := c.ShouldBindBodyWithJSON(&loginUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validate data ---------------------------------------------------------------
	if valid, err := govalidator.ValidateStruct(loginUser); !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Where("email = ?", loginUser.Email).First(&userData).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// Comparing user password and input password -------------------------------------
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginUser.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Username atau password salah"})
		return
	}

	// Make token ----------------------------------------------------------------------
	token, err := helpers.GenerateToken(userData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	cookieName := "jwt_token"
	cookieValue := token
	maxAge := 3600
	path := "/"
	domain := "localhost"
	secure := false
	httpOnly := true

	c.SetCookie(cookieName, cookieValue, maxAge, path, domain, secure, httpOnly)
	c.JSON(http.StatusOK, gin.H{"message": "Login success", "data": map[string]any{"email": userData.Email, "username": userData.Username}, "token": token})

}

func Register(c *gin.Context) {
	var user models.User
	var userList []models.User

	//Binding-------------------------------------------------------------------------
	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validation-----------------------------------------------------------------
	if valid, err := govalidator.ValidateStruct(user); !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if check := database.DB.Find(&models.User{}, "email = ?", user.Email).RowsAffected; check > 0 {
		c.AbortWithStatusJSON(http.StatusFound, gin.H{"message": "Email is exist"})
		return
	}

	// User Data Management-------------------------------------------------------------
	database.DB.Find(&userList)
	id, err := helpers.GetId()
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			id = 1
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	user.Id = int64(id)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	user.Password = string(hashedPassword)

	// File check ----------------------------------------------------------------------------------
	file, _ := c.FormFile("profile_picture")
	if file == nil {
		user.ProfilUrl = "default.jpg"
	} else {
		typeFile := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
		user.ProfilUrl = strconv.Itoa(int(user.Id)) + "_profile." + typeFile
		if valid := helpers.CheckFile(typeFile); !valid && typeFile != "" {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "File is not an image"})
			return
		}
		savePath := "./assets/profiles/" + user.ProfilUrl
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
			return
		}
	}

	database.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Register success", "data": user})
}

func Logout(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User logged out! Good bye!"})
}

func UpdateUser(c *gin.Context) {
	var updateUser models.UpdateUser
	var user models.User
	id := c.Param("userId")
	file, _ := c.FormFile("profile_picture")

	// Validate token ----------------------------------------------------------------------
	if valid := helpers.ValidateUser(c); !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access unauthorized"})
		return
	}

	//Binding-------------------------------------------------------------------------
	if err := c.ShouldBind(&updateUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check and validate data
	if file != nil {
		typeFile := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
		updateUser.ProfilUrl = id + "_profile." + typeFile
		if valid := helpers.CheckFile(typeFile); !valid && typeFile != "" {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"message": "File is not an image"})
			return
		}
	}
	if valid, err := govalidator.ValidateStruct(updateUser); !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updateUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		updateUser.Password = string(hashedPassword)
	}

	// Update Data
	updateUser.UpdatedAt = time.Now()
	if database.DB.Model(&user).Where("id = ?", id).Updates(&updateUser).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Data is not found"})
		return
	}
	// Handle profile image----------------------------------------------------------

	if file != nil {
		savePath := "./assets/profiles/" + updateUser.ProfilUrl
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update user success"})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("userId")

	// Validate token
	if valid := helpers.ValidateUser(c); !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Access unauthorized"})
		return
	}

	// Delete profile image
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
	helpers.DeletePhoto("profile", user.ProfilUrl)

	// Delete data
	if database.DB.Delete(&user, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Delete failed"})
		return
	}
	postDir := "./assets/posts/"
	prefix := id + "_posts"
	if err := helpers.DeleteUserPhoto(postDir, prefix); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Delete failed"})
		return
	}
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Delete data success"})
}
