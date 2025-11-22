package services

import (
	"context"
	"errors"
	"time"

	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
)

type SummaryPeriod string

const (
	PeriodDaily   SummaryPeriod = "daily"
	PeriodWeekly  SummaryPeriod = "weekly"
	PeriodMonthly SummaryPeriod = "monthly"
	PeriodYearly  SummaryPeriod = "yearly"
	PeriodCustom  SummaryPeriod = "custom"
)

type GroupBy string

const (
	GroupNone        GroupBy = "none"
	GroupCategory    GroupBy = "category"
	GroupSubcategory GroupBy = "subcategory"
)

type SummaryTotals struct {
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
	Savings  float64 `json:"savings"`
}

type CategorySummary struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Income       float64 `json:"income"`
	Expenses     float64 `json:"expenses"`
}

type SubcategorySummary struct {
	SubcategoryID   string  `json:"subcategoryId"`
	SubcategoryName string  `json:"subcategoryName"`
	Income          float64 `json:"income"`
	Expenses        float64 `json:"expenses"`
}

type SummaryReport struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	Totals     SummaryTotals        `json:"totals"`
	ByCategory []CategorySummary    `json:"byCategory,omitempty"`
	BySubcat   []SubcategorySummary `json:"bySubcategory,omitempty"`
}

type ReportService interface {
	GetSummary(ctx context.Context, userID string, period SummaryPeriod, start, end *time.Time, groupBy GroupBy) (*SummaryReport, error)
}

type reportService struct {
	txRepo  repositories.TransactionRepository
	catRepo repositories.CategoryRepository
}

func NewReportService(txRepo repositories.TransactionRepository, catRepo repositories.CategoryRepository) ReportService {
	return &reportService{
		txRepo:  txRepo,
		catRepo: catRepo,
	}
}

func (s *reportService) GetSummary(ctx context.Context, userID string, period SummaryPeriod, start, end *time.Time, groupBy GroupBy) (*SummaryReport, error) {
	var from, to time.Time
	now := time.Now().UTC()

	switch period {
	case PeriodDaily:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		to = from.Add(24*time.Hour - time.Nanosecond)
	case PeriodWeekly:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, time.UTC)
		to = from.Add(7*24*time.Hour - time.Nanosecond)
	case PeriodMonthly:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case PeriodYearly:
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(1, 0, 0).Add(-time.Nanosecond)
	case PeriodCustom:
		if start == nil || end == nil {
			return nil, errors.New("start and end required for custom period")
		}
		from = *start
		to = *end
	default:
		return nil, errors.New("invalid period")
	}

	txs, err := s.txRepo.ListByDateRange(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	report := &SummaryReport{}
	report.Period.Start = from
	report.Period.End = to

	catName := map[string]string{}

	var totalIncome, totalExpenses float64
	for _, tx := range txs {
		if tx.Type == models.TransactionTypeIncome {
			totalIncome += tx.Amount
		} else if tx.Type == models.TransactionTypeExpense {
			totalExpenses += tx.Amount
		}
	}
	report.Totals = SummaryTotals{
		Income:   totalIncome,
		Expenses: totalExpenses,
		Savings:  totalIncome - totalExpenses,
	}

	switch groupBy {
	case GroupCategory:
		agg := map[string]*CategorySummary{}
		for _, tx := range txs {
			if tx.CategoryID == nil {
				continue
			}
			id := *tx.CategoryID
			if _, ok := agg[id]; !ok {
				// Lazy load category name
				if _, ok := catName[id]; !ok {
					c, _ := s.catRepo.FindByID(ctx, id, userID)
					if c != nil {
						catName[id] = c.Name
					}
				}
				agg[id] = &CategorySummary{
					CategoryID:   id,
					CategoryName: catName[id],
				}
			}
			if tx.Type == models.TransactionTypeIncome {
				agg[id].Income += tx.Amount
			} else {
				agg[id].Expenses += tx.Amount
			}
		}
		for _, v := range agg {
			report.ByCategory = append(report.ByCategory, *v)
		}
	case GroupSubcategory:
		agg := map[string]*SubcategorySummary{}
		for _, tx := range txs {
			if tx.SubcategoryID == nil {
				continue
			}
			id := *tx.SubcategoryID
			if _, ok := agg[id]; !ok {
				if _, ok := catName[id]; !ok {
					c, _ := s.catRepo.FindByID(ctx, id, userID)
					if c != nil {
						catName[id] = c.Name
					}
				}
				agg[id] = &SubcategorySummary{
					SubcategoryID:   id,
					SubcategoryName: catName[id],
				}
			}
			if tx.Type == models.TransactionTypeIncome {
				agg[id].Income += tx.Amount
			} else {
				agg[id].Expenses += tx.Amount
			}
		}
		for _, v := range agg {
			report.BySubcat = append(report.BySubcat, *v)
		}
	case GroupNone:
		// nothing extra
	default:
		return nil, errors.New("invalid groupBy")
	}

	return report, nil
}
