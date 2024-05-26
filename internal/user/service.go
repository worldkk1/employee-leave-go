package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/worldkk1/employee-leave-go/internal/app/database"
	"github.com/worldkk1/employee-leave-go/internal/app/models"
)

type CreateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type EditUserInput struct {
	Name string `json:"name"`
}

type LeaveDetail struct {
	LeaveType string `json:"leaveType"`
	Remaining int    `json:"remaining"`
	Used      int    `json:"used"`
}

type UserDetailResponse struct {
	models.User
	LeaveDetail []LeaveDetail `json:"leaveDetail"`
}

func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: "",
	}
	database.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func EditUser(c *gin.Context) {
	var input EditUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	database.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUserDetail(c *gin.Context) {
	userId := c.Param("id")
	var user models.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var userLeaves []models.UserLeave
	database.DB.Model(&models.UserLeave{}).Preload("LeaveType").Where("user_id = ?", userId).Find(&userLeaves)

	var leaveDetail []LeaveDetail
	for _, leave := range userLeaves {
		leaveDetail = append(leaveDetail, LeaveDetail{
			LeaveType: leave.LeaveType.Name,
			Remaining: leave.Remaining,
			Used:      leave.Used,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": UserDetailResponse{
		user,
		leaveDetail,
	}})
}
