package converters

import (
	"AggregationService/internal/domain/models/dto"
	"AggregationService/internal/domain/models/entity"
)

type SubscriptionConverter struct {
}

func New() *SubscriptionConverter {
	return &SubscriptionConverter{}
}

func (c *SubscriptionConverter) ToSubscriptionEntity(req *dto.CreateSubscriptionRequest) *entity.Subscription {
	return &entity.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}

func (c *SubscriptionConverter) ToSubscriptionDTO(sub *entity.Subscription) *dto.SubscriptionResponse {
	return &dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
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
		sub.EndDate = req.EndDate
	}
}

func (c *SubscriptionConverter) ToSubscriptionDTOs(subs []*entity.Subscription) []*dto.SubscriptionResponse {
	result := make([]*dto.SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		result = append(result, c.ToSubscriptionDTO(sub))
	}
	return result
}
