package converters

import (
	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/domain/models/entity"
	"AggregationService/internal/pkg/utils"
	"time"
)

type SubscriptionConverter struct {
}

func New() *SubscriptionConverter {
	return &SubscriptionConverter{}
}

func (c *SubscriptionConverter) ToSubscriptionEntity(req *dto.CreateSubscriptionRequest) *entity.Subscription {
	startDate, _ := utils.ParseMonthYearToTime(req.StartDate)
	var endDate *time.Time
	if req.EndDate != nil {
		ed, _ := utils.ParseMonthYearToTime(*req.EndDate)
		endDate = &ed
	}
	return &entity.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

func (c *SubscriptionConverter) ToSubscriptionDTO(sub *entity.Subscription) *dto.SubscriptionResponse {
	return &dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   utils.TimeToMonthYear(sub.StartDate),
		EndDate: func() *string {
			if sub.EndDate == nil {
				return nil
			}
			s := utils.TimeToMonthYear(*sub.EndDate)
			return &s
		}(),
		CreatedAt: sub.CreatedAt,
		UpdatedAt: sub.UpdatedAt,
	}
}

func (c *SubscriptionConverter) ApplyUpdateToEntity(sub *entity.Subscription, req *dto.UpdateSubscriptionRequest) {
	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.Price = *req.Price
	}
	if req.EndDate != nil {
		ed, _ := utils.ParseMonthYearToTime(*req.EndDate)
		sub.EndDate = &ed
	}
}

func (c *SubscriptionConverter) ToSubscriptionDTOs(subs []*entity.Subscription) []*dto.SubscriptionResponse {
	result := make([]*dto.SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		result = append(result, c.ToSubscriptionDTO(sub))
	}
	return result
}
