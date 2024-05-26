package userleave

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/now"
	"github.com/worldkk1/employee-leave-go/internal/app/database"
	"github.com/worldkk1/employee-leave-go/internal/app/models"
)

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
