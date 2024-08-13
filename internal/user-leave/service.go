package userleave

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/now"
	"github.com/worldkk1/employee-leave-go/internal/app/database"
	"github.com/worldkk1/employee-leave-go/internal/app/models"
	"gorm.io/gorm"
)

type RequestLeaveInput struct {
	UserId        uuid.UUID `json:"userId" binding:"required"`
	LeaveTypesId  uuid.UUID `json:"leaveTypesId" binding:"required"`
	StartDate     string    `json:"startDate"`
	EndDate       string    `json:"endDate" binding:"required"`
	TotalLeaveDay int       `json:"totalLeaveDay" binding:"required"`
	Reason        string    `json:"reason"`
	AttachmentUrl string    `json:"attachmentUrl"`
}

type RequestLeaveResponse struct {
	UserId    uuid.UUID `json:"userId"`
	Remaining int       `json:"remaining"`
	Used      int       `json:"used"`
}

type LeaveDetailResponse struct {
	Id            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"userId"`
	LeaveTypesId  uuid.UUID `json:"leaveTypesId"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	TotalLeaveDay int       `json:"totalLeaveDay,omitempty"`
	Reason        *string   `json:"reason"`
	AttachmentURL *string   `json:"attachmentURL"`
	ApproveBy     *string   `json:"approveBy"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func RequestLeave(c *gin.Context) {
	var input RequestLeaveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		queryUserLeaves := tx.Model(models.UserLeave{})
		queryUserLeaves.Where("user_id = ?", input.UserId)
		queryUserLeaves.Where("leave_types_id = ? and remaining >= ?", input.LeaveTypesId, input.TotalLeaveDay)
		queryUserLeaves.Where("start_date <= ? and end_date >= ?", now, now)
		updatedUserLeave := queryUserLeaves.Updates(map[string]interface{}{
			"remaining": gorm.Expr("remaining - ?", input.TotalLeaveDay),
			"used":      gorm.Expr("used + ?", input.TotalLeaveDay),
		})
		if updatedUserLeave.Error != nil || updatedUserLeave.RowsAffected == 0 {
			return errors.New("cannot update user leave")
		}

		leave := models.UserLeaveRecord{
			UserId:        input.UserId,
			LeaveTypesId:  input.LeaveTypesId,
			StartDate:     startDate,
			EndDate:       endDate,
			TotalLeaveDay: input.TotalLeaveDay,
			Reason:        &input.Reason,
			AttachmentURL: &input.AttachmentUrl,
		}
		err = tx.Create(&leave).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userLeave models.UserLeave
	if err := database.DB.Where("user_id = ?", input.UserId).First(&userLeave).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": RequestLeaveResponse{
		UserId:    input.UserId,
		Remaining: userLeave.Remaining,
		Used:      userLeave.Used,
	}})
}

func GetLeaveDetail(c *gin.Context) {
	var leaveRecord models.UserLeaveRecord
	if err := database.DB.Where("id = ?", c.Param("leaveRecordId")).First(&leaveRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": LeaveDetailResponse{
		Id:            leaveRecord.Id,
		UserId:        leaveRecord.UserId,
		LeaveTypesId:  leaveRecord.LeaveTypesId,
		StartDate:     leaveRecord.StartDate,
		EndDate:       leaveRecord.EndDate,
		Reason:        leaveRecord.Reason,
		AttachmentURL: leaveRecord.AttachmentURL,
		ApproveBy:     leaveRecord.ApproveBy,
		CreatedAt:     leaveRecord.CreatedAt,
		UpdatedAt:     leaveRecord.UpdatedAt,
	}})
}

func AllocateLeaves(userId uuid.UUID) {
	startDate := now.BeginningOfYear()
	endDate := now.EndOfYear()

	var activeUserLeaves []models.UserLeave
	queryUserLeaves := database.DB.Where("user_id = ?", userId)
	queryUserLeaves.Where("start_date >= ? and end_date <= ?", startDate, endDate)
	queryUserLeaves.Find(&activeUserLeaves)

	var alreadyHaveLeaveTypeIds []uuid.UUID
	for _, userLeave := range activeUserLeaves {
		alreadyHaveLeaveTypeIds = append(alreadyHaveLeaveTypeIds, userLeave.LeaveTypesId)
	}

	var activeLeaves []models.LeaveType
	queryLeaveTypes := database.DB.Where("is_active = ?", true)
	if alreadyHaveLeaveTypeIds != nil {
		queryLeaveTypes.Where("id not in ?", alreadyHaveLeaveTypeIds)
	}
	queryLeaveTypes.Find(&activeLeaves)

	var createdUserLeave []models.UserLeave
	for _, leave := range activeLeaves {
		createdUserLeave = append(createdUserLeave, models.UserLeave{
			UserId:       userId,
			LeaveTypesId: leave.Id,
			Remaining:    calculateRatioLeaves(leave),
			Used:         0,
			StartDate:    now.BeginningOfDay(),
			EndDate:      endDate,
		})
	}

	if createdUserLeave != nil {
		database.DB.Create(createdUserLeave)
	}
}

func calculateRatioLeaves(leaveType models.LeaveType) int {
	totalLeaves := leaveType.TotalLeaves
	IsRatioAllocate := leaveType.IsRatioAllocate
	if !IsRatioAllocate {
		return totalLeaves
	}

	currentMonth := int(time.Now().Month())
	availableLeaveMonth := 12 - (currentMonth + 1)
	ratioLeaves := (totalLeaves * availableLeaveMonth) / 12

	return ratioLeaves
}
