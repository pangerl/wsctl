// Package models 业务模型包
// @Author lanpang
// @Date 2024/9/12 下午5:21:00
// @Desc 租户相关业务模型
package models

// Tenant 租户模型
// 包含企业列表和相关业务逻辑
type Tenant struct {
	Corps []*Corp `json:"corps"` // 企业列表
}

// Corp 企业模型
// 包含企业基本信息、用户统计、消息统计等业务数据
type Corp struct {
	CorpID               string `json:"corp_id"`                 // 企业ID
	ConvEnabled          bool   `json:"conv_enabled"`            // 是否启用会话功能
	CorpName             string `json:"corp_name"`               // 企业名称
	MessageNum           int64  `json:"message_num"`             // 消息总数
	YesterdayMessageNum  int64  `json:"yesterday_message_num"`   // 昨日消息数
	UserNum              int    `json:"user_num"`                // 用户总数
	CustomerNum          int64  `json:"customer_num"`            // 客户总数
	CustomerGroupNum     int    `json:"customer_group_num"`      // 客户群组数
	CustomerGroupUserNum int    `json:"customer_group_user_num"` // 客户群组用户数
	DAUNum               int    `json:"dau_num"`                 // 日活跃用户数
	WAUNum               int    `json:"wau_num"`                 // 周活跃用户数
	MAUNum               int    `json:"mau_num"`                 // 月活跃用户数
}

// GetCorpByID 根据企业ID获取企业信息
func (t *Tenant) GetCorpByID(corpID string) *Corp {
	for _, corp := range t.Corps {
		if corp.CorpID == corpID {
			return corp
		}
	}
	return nil
}

// GetActiveCorps 获取启用会话功能的企业列表
func (t *Tenant) GetActiveCorps() []*Corp {
	var activeCorps []*Corp
	for _, corp := range t.Corps {
		if corp.ConvEnabled {
			activeCorps = append(activeCorps, corp)
		}
	}
	return activeCorps
}

// GetTotalUsers 获取租户下所有企业的用户总数
func (t *Tenant) GetTotalUsers() int {
	total := 0
	for _, corp := range t.Corps {
		total += corp.UserNum
	}
	return total
}

// GetTotalMessages 获取租户下所有企业的消息总数
func (t *Tenant) GetTotalMessages() int64 {
	var total int64
	for _, corp := range t.Corps {
		total += corp.MessageNum
	}
	return total
}

// IsActive 检查企业是否活跃（启用会话功能）
func (c *Corp) IsActive() bool {
	return c.ConvEnabled
}

// GetMessageGrowth 获取消息增长数（今日相比昨日）
func (c *Corp) GetMessageGrowth() int64 {
	return c.MessageNum - c.YesterdayMessageNum
}

// GetMessageGrowthRate 获取消息增长率
func (c *Corp) GetMessageGrowthRate() float64 {
	if c.YesterdayMessageNum == 0 {
		return 0
	}
	return float64(c.GetMessageGrowth()) / float64(c.YesterdayMessageNum) * 100
}

// GetCustomerEngagement 获取客户参与度（客户群组用户数/客户总数）
func (c *Corp) GetCustomerEngagement() float64 {
	if c.CustomerNum == 0 {
		return 0
	}
	return float64(c.CustomerGroupUserNum) / float64(c.CustomerNum) * 100
}
