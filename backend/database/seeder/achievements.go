package seeder

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"learningAssistant-backend/models"
)

type AchievementSeed struct {
	Code        string
	Name        string
	Description string
	Category    string
	Icon        string
	Condition   map[string]interface{}
}

var achievementSeeds = []AchievementSeed{
	// 任务达人 · 基础
	{Code: "task_starter_5", Name: "🎯 新手任务者", Description: "完成 5 个任务", Category: "task_master_basic", Icon: "mdi:target", Condition: map[string]interface{}{"type": "task_total_completed", "value": 5}},
	{Code: "task_diligent_20", Name: "📝 勤奋执行者", Description: "完成 20 个任务", Category: "task_master_basic", Icon: "mdi:note-edit", Condition: map[string]interface{}{"type": "task_total_completed", "value": 20}},
	{Code: "task_runner_50", Name: "🚀 连续推进者", Description: "完成 50 个任务", Category: "task_master_basic", Icon: "mdi:rocket-launch", Condition: map[string]interface{}{"type": "task_total_completed", "value": 50}},

	// 任务达人 · 进阶
	{Code: "task_star_100", Name: "🌟 百任务之星", Description: "完成 100 个任务", Category: "task_master_advanced", Icon: "mdi:star-circle", Condition: map[string]interface{}{"type": "task_total_completed", "value": 100}},
	{Code: "task_efficiency_200", Name: "💠 效率达人", Description: "完成 200 个任务", Category: "task_master_advanced", Icon: "mdi:diamond-stone", Condition: map[string]interface{}{"type": "task_total_completed", "value": 200}},
	{Code: "task_focus_300", Name: "🔥 专注执行者", Description: "完成 300 个任务", Category: "task_master_advanced", Icon: "mdi:fire", Condition: map[string]interface{}{"type": "task_total_completed", "value": 300}},

	// 任务达人 · 终极
	{Code: "task_iron_500", Name: "🏆 学习铁人", Description: "完成 500 个任务", Category: "task_master_ultimate", Icon: "mdi:trophy", Condition: map[string]interface{}{"type": "task_total_completed", "value": 500}},
	{Code: "task_galaxy_1000", Name: "🌌 任务银河级选手", Description: "完成 1000 个任务", Category: "task_master_ultimate", Icon: "mdi:galaxy", Condition: map[string]interface{}{"type": "task_total_completed", "value": 1000}},
	{Code: "task_king_2000", Name: "👑 终极任务王", Description: "完成 2000 个任务", Category: "task_master_ultimate", Icon: "mdi:crown", Condition: map[string]interface{}{"type": "task_total_completed", "value": 2000}},

	// 连续完成任务成就
	{Code: "streak_3", Name: "📅 坚持 3 天", Description: "连续完成全部任务 3 天", Category: "streak_basic", Icon: "mdi:calendar-check", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 3}},
	{Code: "streak_7", Name: "📆 坚持 7 天（第一周）", Description: "连续完成全部任务 7 天", Category: "streak_basic", Icon: "mdi:calendar-week", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 7}},
	{Code: "streak_10", Name: "🔒 坚持 10 天（特别成就）", Description: "连续完成全部任务 10 天", Category: "streak_basic", Icon: "mdi:lock-check", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 10}},
	{Code: "streak_15", Name: "🔄 半月坚持者", Description: "连续完成全部任务 15 天", Category: "streak_mid", Icon: "mdi:repeat-variant", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 15}},
	{Code: "streak_30", Name: "🌓 月度坚持者", Description: "连续完成全部任务 30 天", Category: "streak_mid", Icon: "mdi:calendar-range", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 30}},
	{Code: "streak_50", Name: "🌙 连续 50 天学习者", Description: "连续完成全部任务 50 天", Category: "streak_high", Icon: "mdi:weather-night", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 50}},
	{Code: "streak_100", Name: "🌞 持续 100 天精进者", Description: "连续完成全部任务 100 天", Category: "streak_high", Icon: "mdi:white-balance-sunny", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 100}},
	{Code: "streak_365", Name: "🔥 不间断的学习传奇", Description: "连续完成全部任务 365 天", Category: "streak_high", Icon: "mdi:fire-circle", Condition: map[string]interface{}{"type": "streak_task_completion", "value": 365}},

	// 自习室参与次数
	{Code: "studyroom_first_join", Name: "📚 初入学习室", Description: "第一次加入自习室", Category: "studyroom_participation", Icon: "mdi:book-open-variant", Condition: map[string]interface{}{"type": "studyroom_join_count", "value": 1}},
	{Code: "studyroom_partner_10", Name: "👥 学习伙伴", Description: "参与自习室 10 次", Category: "studyroom_participation", Icon: "mdi:account-multiple", Condition: map[string]interface{}{"type": "studyroom_join_count", "value": 10}},
	{Code: "studyroom_active_30", Name: "🔊 活跃室友", Description: "参与自习室 30 次", Category: "studyroom_participation", Icon: "mdi:bullhorn", Condition: map[string]interface{}{"type": "studyroom_join_count", "value": 30}},
	{Code: "studyroom_resident_60", Name: "🎧 专注空间居民", Description: "参与自习室 60 次", Category: "studyroom_participation", Icon: "mdi:headphones", Condition: map[string]interface{}{"type": "studyroom_join_count", "value": 60}},

	// 自习室时长
	{Code: "studyroom_focus_5h", Name: "⏳ 累计专注 5 小时", Description: "自习室累计时长 5 小时", Category: "studyroom_duration", Icon: "mdi:timer-sand", Condition: map[string]interface{}{"type": "studyroom_duration_hours", "value": 5}},
	{Code: "studyroom_focus_20h", Name: "🕒 累计专注 20 小时", Description: "自习室累计时长 20 小时", Category: "studyroom_duration", Icon: "mdi:clock-outline", Condition: map[string]interface{}{"type": "studyroom_duration_hours", "value": 20}},
	{Code: "studyroom_night_owl", Name: "🌙 夜猫学习者", Description: "晚上 22:00 后自习超过 2 小时", Category: "studyroom_duration", Icon: "mdi:weather-night", Condition: map[string]interface{}{"type": "studyroom_night_hours", "value": 2, "mode": "single_or_total"}},
	{Code: "studyroom_focus_100h", Name: "🎯 百小时修行者", Description: "自习室累计时长 100 小时", Category: "studyroom_duration", Icon: "mdi:bullseye-arrow", Condition: map[string]interface{}{"type": "studyroom_duration_hours", "value": 100}},

	// 自习室社交互动
	{Code: "studyroom_first_chat", Name: "💬 学习室发言者", Description: "第一次聊天互动", Category: "studyroom_social", Icon: "mdi:message-text", Condition: map[string]interface{}{"type": "studyroom_chat_count", "value": 1}},
	{Code: "studyroom_liked_10", Name: "👍 互助学习者", Description: "被他人点赞 10 次", Category: "studyroom_social", Icon: "mdi:thumb-up", Condition: map[string]interface{}{"type": "studyroom_likes_received", "value": 10}},
	{Code: "studyroom_support_20", Name: "❤️ 学习鼓励家", Description: "给别人点赞 20 次", Category: "studyroom_social", Icon: "mdi:heart", Condition: map[string]interface{}{"type": "studyroom_likes_given", "value": 20}},

	// 团队任务
	{Code: "team_joiner", Name: "👨‍👩‍👧 加入小队", Description: "第一次参与团队任务", Category: "team_tasks", Icon: "mdi:account-group", Condition: map[string]interface{}{"type": "team_task_completed", "value": 1}},
	{Code: "team_contributor_10", Name: "🔧 团队贡献者", Description: "完成团队任务 10 次", Category: "team_tasks", Icon: "mdi:handshake", Condition: map[string]interface{}{"type": "team_task_completed", "value": 10}},
	{Code: "team_core_30", Name: "🛠️ 核心执行者", Description: "完成团队任务 30 次", Category: "team_tasks", Icon: "mdi:shield-check", Condition: map[string]interface{}{"type": "team_task_completed", "value": 30}},
}

// SeedAchievements 创建或更新成就定义
func SeedAchievements(db *gorm.DB) error {
	for _, seed := range achievementSeeds {
		conditionBytes, err := json.Marshal(seed.Condition)
		if err != nil {
			return fmt.Errorf("marshal achievement condition %s: %w", seed.Code, err)
		}

		record := models.Achievement{
			Code:        seed.Code,
			Name:        seed.Name,
			Description: seed.Description,
			Category:    seed.Category,
			Icon:        seed.Icon,
			Condition:   string(conditionBytes),
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "description", "category", "icon", "condition"}),
		}).Create(&record).Error; err != nil {
			return fmt.Errorf("seed achievement %s: %w", seed.Code, err)
		}
	}
	return nil
}
