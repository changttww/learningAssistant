package task

import (
	"fmt"
	"math"
	"time"

	"learningAssistant-backend/database"
	"learningAssistant-backend/models"
)

// DailyBarStat 封装单个时间粒度的任务统计
type DailyBarStat struct {
	Day       string `json:"day"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
	Rate      int    `json:"rate"`
}

// BarStats 返回指定范围的柱状图数据
type BarStats struct {
	Range     string         `json:"range"`
	StartDate string         `json:"start_date"`
	WeekStart string         `json:"week_start,omitempty"`
	Data      []DailyBarStat `json:"data"`
}

// MonthlyCompletionStat 聚合当月截止的任务完成率
type MonthlyCompletionStat struct {
	Month     string `json:"month"`
	TotalDue  int64  `json:"total_due"`
	Completed int64  `json:"completed"`
	Rate      int    `json:"rate"`
}

// GetBarStats 统计当前自然日/周/月/季度的任务到期完成情况
func GetBarStats(rangeKey string, userID uint64) (*BarStats, error) {
	switch rangeKey {
	case "day":
		return calcDailyStats(userID)
	case "month":
		return calcMonthlyStats(userID)
	case "quarter":
		return calcQuarterStats(userID)
	default:
		return calcWeeklyStats(userID)
	}
}

// GetWeeklyBarStats 为兼容保留的周统计方法
func GetWeeklyBarStats(userID uint64) (*BarStats, error) {
	return calcWeeklyStats(userID)
}

// GetRecentMonthlyCompletion 返回最近 N 个月（含当月）的截止任务完成率
func GetRecentMonthlyCompletion(userID uint64, months int) ([]MonthlyCompletionStat, error) {
	if months <= 0 {
		months = 3
	}

	now := time.Now().In(time.Local)
	stats := make([]MonthlyCompletionStat, 0, months)

	for i := months - 1; i >= 0; i-- {
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0)

		stat, err := queryMonthlyCompletion(userID, monthStart, monthEnd)
		if err != nil {
			return nil, err
		}
		stats = append(stats, *stat)
	}

	return stats, nil
}

func queryMonthlyCompletion(userID uint64, start, end time.Time) (*MonthlyCompletionStat, error) {
	var total, completed int64

	db := database.GetDB()
	if err := db.Model(&models.Task{}).
		Where("owner_user_id = ? AND deleted_at IS NULL AND due_at IS NOT NULL AND due_at >= ? AND due_at < ?", userID, start, end).
		Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Task{}).
		Where("owner_user_id = ? AND status = 2 AND completed_at IS NOT NULL AND deleted_at IS NULL AND due_at IS NOT NULL AND due_at >= ? AND due_at < ?", userID, start, end).
		Count(&completed).Error; err != nil {
		return nil, err
	}

	return &MonthlyCompletionStat{
		Month:     start.Format("2006-01"),
		TotalDue:  total,
		Completed: completed,
		Rate:      calcRate(completed, total),
	}, nil
}

func calcWeeklyStats(userID uint64) (*BarStats, error) {
	startOfWeek := getWeekStart(time.Now().In(time.Local))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var rows []struct {
		Offset    int    `gorm:"column:offset"`
		Day       string `gorm:"column:day"`
		Total     int64  `gorm:"column:total"`
		Completed int64  `gorm:"column:completed"`
	}

	sql := `WITH RECURSIVE week_days AS (
                SELECT 0 AS offset
                UNION ALL
                SELECT offset + 1 FROM week_days WHERE offset < 6
            ),
            due_stats AS (
                SELECT DATE(due_at) AS day_date, COUNT(*) AS total
                FROM tasks
                WHERE owner_user_id = ?
                  AND deleted_at IS NULL
                  AND due_at IS NOT NULL
                  AND due_at >= ?
                  AND due_at < ?
                GROUP BY DATE(due_at)
            ),
            completed_due AS (
                SELECT DATE(due_at) AS day_date, COUNT(*) AS completed
                FROM tasks
                WHERE owner_user_id = ?
                  AND status = 2
                  AND completed_at IS NOT NULL
                  AND deleted_at IS NULL
                  AND due_at IS NOT NULL
                  AND due_at >= ?
                  AND due_at < ?
                GROUP BY DATE(due_at)
            )
            SELECT
                w.offset,
                DATE_FORMAT(DATE_ADD(?, INTERVAL w.offset DAY), '%Y-%m-%d') AS day,
                COALESCE(ds.total, 0) AS total,
                COALESCE(cd.completed, 0) AS completed
            FROM week_days w
            LEFT JOIN due_stats ds ON ds.day_date = DATE_ADD(?, INTERVAL w.offset DAY)
            LEFT JOIN completed_due cd ON cd.day_date = DATE_ADD(?, INTERVAL w.offset DAY)
            ORDER BY w.offset;`

	if err := database.GetDB().Raw(sql,
		userID, startOfWeek, endOfWeek,
		userID, startOfWeek, endOfWeek,
		startOfWeek, startOfWeek, startOfWeek,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	dayLabels := []string{"周一", "周二", "周三", "周四", "周五", "周六", "周日"}
	data := make([]DailyBarStat, 0, len(dayLabels))

	for _, row := range rows {
		data = append(data, DailyBarStat{
			Day:       safeLabel(dayLabels, row.Offset, row.Day),
			Total:     row.Total,
			Completed: row.Completed,
			Rate:      calcRate(row.Completed, row.Total),
		})
	}

	start := startOfWeek.Format("2006-01-02")
	return &BarStats{
		Range:     "week",
		StartDate: start,
		WeekStart: start,
		Data:      data,
	}, nil
}

func calcDailyStats(userID uint64) (*BarStats, error) {
	now := time.Now().In(time.Local)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var row struct {
		Day       string `gorm:"column:day"`
		Total     int64  `gorm:"column:total"`
		Completed int64  `gorm:"column:completed"`
	}

	sql := `SELECT
				DATE_FORMAT(?, '%Y-%m-%d') AS day,
				(SELECT COUNT(*) FROM tasks WHERE owner_user_id = ? AND deleted_at IS NULL AND due_at IS NOT NULL AND due_at >= ? AND due_at < ?) AS total,
				(SELECT COUNT(*) FROM tasks WHERE owner_user_id = ? AND status = 2 AND completed_at IS NOT NULL AND deleted_at IS NULL AND due_at IS NOT NULL AND due_at >= ? AND due_at < ?) AS completed`

	if err := database.GetDB().Raw(sql,
		start,
		userID, start, end,
		userID, start, end,
	).Scan(&row).Error; err != nil {
		return nil, err
	}

	stat := DailyBarStat{
		Day:       "今日",
		Total:     row.Total,
		Completed: row.Completed,
		Rate:      calcRate(row.Completed, row.Total),
	}

	return &BarStats{
		Range:     "day",
		StartDate: start.Format("2006-01-02"),
		Data:      []DailyBarStat{stat},
	}, nil
}

func calcMonthlyStats(userID uint64) (*BarStats, error) {
	now := time.Now().In(time.Local)
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)
	daysInMonth := end.AddDate(0, 0, -1).Day()

	var rows []struct {
		Offset    int    `gorm:"column:offset"`
		Day       string `gorm:"column:day"`
		Total     int64  `gorm:"column:total"`
		Completed int64  `gorm:"column:completed"`
	}

	sql := `WITH RECURSIVE month_days AS (
                SELECT 0 AS offset
                UNION ALL
                SELECT offset + 1 FROM month_days WHERE offset + 1 < ?
            ),
            due_stats AS (
                SELECT DATE(due_at) AS day_date, COUNT(*) AS total
                FROM tasks
                WHERE owner_user_id = ?
                  AND deleted_at IS NULL
                  AND due_at IS NOT NULL
                  AND due_at >= ?
                  AND due_at < ?
                GROUP BY DATE(due_at)
            ),
            completed_due AS (
                SELECT DATE(due_at) AS day_date, COUNT(*) AS completed
                FROM tasks
                WHERE owner_user_id = ?
                  AND status = 2
                  AND completed_at IS NOT NULL
                  AND deleted_at IS NULL
                  AND due_at IS NOT NULL
                  AND due_at >= ?
                  AND due_at < ?
                GROUP BY DATE(due_at)
            )
            SELECT
                d.offset,
                DATE_FORMAT(DATE_ADD(?, INTERVAL d.offset DAY), '%Y-%m-%d') AS day,
                COALESCE(ds.total, 0) AS total,
                COALESCE(cd.completed, 0) AS completed
            FROM month_days d
            LEFT JOIN due_stats ds ON ds.day_date = DATE_ADD(?, INTERVAL d.offset DAY)
            LEFT JOIN completed_due cd ON cd.day_date = DATE_ADD(?, INTERVAL d.offset DAY)
            ORDER BY d.offset;`

	if err := database.GetDB().Raw(sql,
		daysInMonth,
		userID, start, end,
		userID, start, end,
		start, start, start,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make([]DailyBarStat, 0, daysInMonth)
	for _, row := range rows {
		label := fmt.Sprintf("%d日", row.Offset+1)
		data = append(data, DailyBarStat{
			Day:       label,
			Total:     row.Total,
			Completed: row.Completed,
			Rate:      calcRate(row.Completed, row.Total),
		})
	}

	return &BarStats{
		Range:     "month",
		StartDate: start.Format("2006-01-02"),
		Data:      data,
	}, nil
}

func calcQuarterStats(userID uint64) (*BarStats, error) {
	now := time.Now().In(time.Local)
	quarterStartMonth := ((int(now.Month())-1)/3)*3 + 1
	start := time.Date(now.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 3, 0)

	var rows []struct {
		Offset    int    `gorm:"column:offset"`
		Day       string `gorm:"column:day"`
		Total     int64  `gorm:"column:total"`
		Completed int64  `gorm:"column:completed"`
	}

	sql := `WITH RECURSIVE months AS (
                SELECT 0 AS offset
                UNION ALL
                SELECT offset + 1 FROM months WHERE offset + 1 < 3
            ),
            due_stats AS (
                SELECT DATE_FORMAT(due_at, '%Y-%m-01') AS month_date, COUNT(*) AS total
                FROM tasks
                WHERE owner_user_id = ?
                  AND deleted_at IS NULL
                  AND due_at IS NOT NULL
                  AND due_at >= ?
                  AND due_at < ?
                GROUP BY YEAR(due_at), MONTH(due_at)
            ),
            completed_due AS (
                SELECT DATE_FORMAT(due_at, '%Y-%m-01') AS month_date, COUNT(*) AS completed
                FROM tasks
                WHERE owner_user_id = ?
                  AND status = 2
                  AND completed_at IS NOT NULL
                  AND deleted_at IS NULL
                  AND due_at IS NOT NULL
                  AND due_at >= ?
                  AND due_at < ?
                GROUP BY YEAR(due_at), MONTH(due_at)
            )
            SELECT
                m.offset,
                DATE_FORMAT(DATE_ADD(?, INTERVAL m.offset MONTH), '%Y-%m-01') AS day,
                COALESCE(ds.total, 0) AS total,
                COALESCE(cd.completed, 0) AS completed
            FROM months m
            LEFT JOIN due_stats ds ON ds.month_date = DATE_ADD(?, INTERVAL m.offset MONTH)
            LEFT JOIN completed_due cd ON cd.month_date = DATE_ADD(?, INTERVAL m.offset MONTH)
            ORDER BY m.offset;`

	if err := database.GetDB().Raw(sql,
		userID, start, end,
		userID, start, end,
		start, start, start,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make([]DailyBarStat, 0, 3)
	for _, row := range rows {
		monthLabel := start.AddDate(0, row.Offset, 0).Month()
		label := fmt.Sprintf("%d月", monthLabel)
		data = append(data, DailyBarStat{
			Day:       label,
			Total:     row.Total,
			Completed: row.Completed,
			Rate:      calcRate(row.Completed, row.Total),
		})
	}

	return &BarStats{
		Range:     "quarter",
		StartDate: start.Format("2006-01-02"),
		Data:      data,
	}, nil
}

// getWeekStart 返回本周周一的零点时间
func getWeekStart(now time.Time) time.Time {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	delta := weekday - 1

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return startOfDay.AddDate(0, 0, -delta)
}

func calcRate(completed, total int64) int {
	if total <= 0 {
		return 0
	}
	rate := int(math.Round(float64(completed) / float64(total) * 100))
	if rate > 100 {
		return 100
	}
	if rate < 0 {
		return 0
	}
	return rate
}

func safeLabel(labels []string, offset int, fallback string) string {
	if offset >= 0 && offset < len(labels) {
		return labels[offset]
	}
	return fallback
}
