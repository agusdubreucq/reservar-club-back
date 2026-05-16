package services

import (
	"context"
	"time"

	"github.com/adubr/reservar-club-back/internal/domain"
	"github.com/adubr/reservar-club-back/internal/repository"
	"github.com/google/uuid"
)

type PricingService struct {
	pricingRepo repository.PricingRuleRepository
}

func NewPricingService(pricingRepo repository.PricingRuleRepository) *PricingService {
	return &PricingService{
		pricingRepo: pricingRepo,
	}
}

func (s *PricingService) CalculatePrice(ctx context.Context, sportID string, startDate, endDate time.Time) (float64, error) {
	rules, err := s.pricingRepo.ListBySportID(ctx, sportID)
	if err != nil {
		return 0, err
	}

	totalPrice := 0.0
	current := startDate

	for current.Before(endDate) {
		nextSlot := current.Add(30 * time.Minute)
		if nextSlot.After(endDate) {
			nextSlot = endDate
		}

		price := s.getPriceForSlot(current, rules)
		durationMin := nextSlot.Sub(current).Minutes()
		slotPrice := (durationMin / 60.0) * price

		totalPrice += slotPrice
		current = nextSlot
	}

	return totalPrice, nil
}

func (s *PricingService) getPriceForSlot(slotTime time.Time, rules []*domain.PricingRule) float64 {
	dayOfWeek := int(slotTime.Weekday())
	hour := slotTime.Hour()

	var bestRule *domain.PricingRule

	for _, rule := range rules {
		if !s.ruleAppliesToTime(rule, dayOfWeek, hour) {
			continue
		}

		if bestRule == nil {
			bestRule = rule
		} else if rule.DayOfWeek != nil && bestRule.DayOfWeek == nil {
			bestRule = rule
		}
	}

	if bestRule != nil {
		return bestRule.PricePerHour
	}

	return 0
}

func (s *PricingService) ruleAppliesToTime(rule *domain.PricingRule, dayOfWeek, hour int) bool {
	if rule.DayOfWeek != nil && *rule.DayOfWeek != dayOfWeek {
		return false
	}

	if hour < rule.StartHour || hour >= rule.EndHour {
		return false
	}

	return true
}

func (s *PricingService) CreatePricingRule(ctx context.Context, sportID string, dayOfWeek *int, startHour, endHour int, pricePerHour float64) (*domain.PricingRule, error) {
	rule := domain.NewPricingRule(uuid.New().String(), sportID, dayOfWeek, startHour, endHour, pricePerHour)
	if err := s.pricingRepo.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *PricingService) ListPricingRules(ctx context.Context, sportID string) ([]*domain.PricingRule, error) {
	return s.pricingRepo.ListBySportID(ctx, sportID)
}

func (s *PricingService) DeletePricingRule(ctx context.Context, id string) error {
	return s.pricingRepo.Delete(ctx, id)
}
