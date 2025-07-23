package models

import (
	"testing"
)

// TestTenant_GetCorpByID 测试根据企业ID获取企业信息
func TestTenant_GetCorpByID(t *testing.T) {
	// 创建测试数据
	tenant := &Tenant{
		Corps: []*Corp{
			{
				CorpID:      "corp1",
				CorpName:    "企业1",
				ConvEnabled: true,
			},
			{
				CorpID:      "corp2",
				CorpName:    "企业2",
				ConvEnabled: false,
			},
		},
	}

	// 测试存在的企业ID
	corp := tenant.GetCorpByID("corp1")
	if corp == nil {
		t.Error("GetCorpByID() returned nil for existing corp")
	} else if corp.CorpID != "corp1" {
		t.Errorf("GetCorpByID() returned wrong corp, got %s, want %s", corp.CorpID, "corp1")
	}

	// 测试不存在的企业ID
	corp = tenant.GetCorpByID("corp3")
	if corp != nil {
		t.Errorf("GetCorpByID() returned non-nil for non-existing corp: %v", corp)
	}
}

// TestTenant_GetActiveCorps 测试获取启用会话功能的企业列表
func TestTenant_GetActiveCorps(t *testing.T) {
	// 创建测试数据
	tenant := &Tenant{
		Corps: []*Corp{
			{
				CorpID:      "corp1",
				CorpName:    "企业1",
				ConvEnabled: true,
			},
			{
				CorpID:      "corp2",
				CorpName:    "企业2",
				ConvEnabled: false,
			},
			{
				CorpID:      "corp3",
				CorpName:    "企业3",
				ConvEnabled: true,
			},
		},
	}

	// 获取活跃企业
	activeCorps := tenant.GetActiveCorps()

	// 验证结果
	if len(activeCorps) != 2 {
		t.Errorf("GetActiveCorps() returned %d corps, want %d", len(activeCorps), 2)
	}

	// 验证返回的企业都是启用会话功能的
	for _, corp := range activeCorps {
		if !corp.ConvEnabled {
			t.Errorf("GetActiveCorps() returned inactive corp: %s", corp.CorpID)
		}
	}
}

// TestTenant_GetTotalUsers 测试获取租户下所有企业的用户总数
func TestTenant_GetTotalUsers(t *testing.T) {
	// 创建测试数据
	tenant := &Tenant{
		Corps: []*Corp{
			{
				CorpID:  "corp1",
				UserNum: 100,
			},
			{
				CorpID:  "corp2",
				UserNum: 200,
			},
			{
				CorpID:  "corp3",
				UserNum: 300,
			},
		},
	}

	// 获取用户总数
	totalUsers := tenant.GetTotalUsers()

	// 验证结果
	expectedTotal := 100 + 200 + 300
	if totalUsers != expectedTotal {
		t.Errorf("GetTotalUsers() = %d, want %d", totalUsers, expectedTotal)
	}
}

// TestTenant_GetTotalMessages 测试获取租户下所有企业的消息总数
func TestTenant_GetTotalMessages(t *testing.T) {
	// 创建测试数据
	tenant := &Tenant{
		Corps: []*Corp{
			{
				CorpID:     "corp1",
				MessageNum: 1000,
			},
			{
				CorpID:     "corp2",
				MessageNum: 2000,
			},
			{
				CorpID:     "corp3",
				MessageNum: 3000,
			},
		},
	}

	// 获取消息总数
	totalMessages := tenant.GetTotalMessages()

	// 验证结果
	expectedTotal := int64(1000 + 2000 + 3000)
	if totalMessages != expectedTotal {
		t.Errorf("GetTotalMessages() = %d, want %d", totalMessages, expectedTotal)
	}
}

// TestCorp_IsActive 测试检查企业是否活跃
func TestCorp_IsActive(t *testing.T) {
	// 创建测试数据
	activeCorp := &Corp{
		CorpID:      "corp1",
		ConvEnabled: true,
	}
	inactiveCorp := &Corp{
		CorpID:      "corp2",
		ConvEnabled: false,
	}

	// 测试活跃企业
	if !activeCorp.IsActive() {
		t.Errorf("IsActive() = false for active corp")
	}

	// 测试非活跃企业
	if inactiveCorp.IsActive() {
		t.Errorf("IsActive() = true for inactive corp")
	}
}

// TestCorp_GetMessageGrowth 测试获取消息增长数
func TestCorp_GetMessageGrowth(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name                string
		messageNum          int64
		yesterdayMessageNum int64
		expectedGrowth      int64
	}{
		{
			name:                "正增长",
			messageNum:          1000,
			yesterdayMessageNum: 800,
			expectedGrowth:      200,
		},
		{
			name:                "零增长",
			messageNum:          1000,
			yesterdayMessageNum: 1000,
			expectedGrowth:      0,
		},
		{
			name:                "负增长",
			messageNum:          800,
			yesterdayMessageNum: 1000,
			expectedGrowth:      -200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corp := &Corp{
				MessageNum:          tt.messageNum,
				YesterdayMessageNum: tt.yesterdayMessageNum,
			}

			growth := corp.GetMessageGrowth()
			if growth != tt.expectedGrowth {
				t.Errorf("GetMessageGrowth() = %d, want %d", growth, tt.expectedGrowth)
			}
		})
	}
}

// TestCorp_GetMessageGrowthRate 测试获取消息增长率
func TestCorp_GetMessageGrowthRate(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name                string
		messageNum          int64
		yesterdayMessageNum int64
		expectedRate        float64
	}{
		{
			name:                "正增长率",
			messageNum:          1000,
			yesterdayMessageNum: 800,
			expectedRate:        25.0, // (1000-800)/800*100 = 25%
		},
		{
			name:                "零增长率",
			messageNum:          1000,
			yesterdayMessageNum: 1000,
			expectedRate:        0.0,
		},
		{
			name:                "负增长率",
			messageNum:          800,
			yesterdayMessageNum: 1000,
			expectedRate:        -20.0, // (800-1000)/1000*100 = -20%
		},
		{
			name:                "昨日消息为零",
			messageNum:          1000,
			yesterdayMessageNum: 0,
			expectedRate:        0.0, // 避免除以零
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corp := &Corp{
				MessageNum:          tt.messageNum,
				YesterdayMessageNum: tt.yesterdayMessageNum,
			}

			rate := corp.GetMessageGrowthRate()
			if rate != tt.expectedRate {
				t.Errorf("GetMessageGrowthRate() = %f, want %f", rate, tt.expectedRate)
			}
		})
	}
}

// TestCorp_GetCustomerEngagement 测试获取客户参与度
func TestCorp_GetCustomerEngagement(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name                 string
		customerNum          int64
		customerGroupUserNum int
		expectedEngagement   float64
	}{
		{
			name:                 "正常参与度",
			customerNum:          1000,
			customerGroupUserNum: 500,
			expectedEngagement:   50.0, // 500/1000*100 = 50%
		},
		{
			name:                 "零参与度",
			customerNum:          1000,
			customerGroupUserNum: 0,
			expectedEngagement:   0.0,
		},
		{
			name:                 "满参与度",
			customerNum:          1000,
			customerGroupUserNum: 1000,
			expectedEngagement:   100.0,
		},
		{
			name:                 "客户数为零",
			customerNum:          0,
			customerGroupUserNum: 500,
			expectedEngagement:   0.0, // 避免除以零
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corp := &Corp{
				CustomerNum:          tt.customerNum,
				CustomerGroupUserNum: tt.customerGroupUserNum,
			}

			engagement := corp.GetCustomerEngagement()
			if engagement != tt.expectedEngagement {
				t.Errorf("GetCustomerEngagement() = %f, want %f", engagement, tt.expectedEngagement)
			}
		})
	}
}

// TestTenant_EmptyCorps 测试空企业列表的情况
func TestTenant_EmptyCorps(t *testing.T) {
	// 创建空租户
	tenant := &Tenant{
		Corps: []*Corp{},
	}

	// 测试各种方法
	if corp := tenant.GetCorpByID("corp1"); corp != nil {
		t.Errorf("GetCorpByID() returned non-nil for empty corps: %v", corp)
	}

	if activeCorps := tenant.GetActiveCorps(); len(activeCorps) != 0 {
		t.Errorf("GetActiveCorps() returned %d corps, want 0", len(activeCorps))
	}

	if totalUsers := tenant.GetTotalUsers(); totalUsers != 0 {
		t.Errorf("GetTotalUsers() = %d, want 0", totalUsers)
	}

	if totalMessages := tenant.GetTotalMessages(); totalMessages != 0 {
		t.Errorf("GetTotalMessages() = %d, want 0", totalMessages)
	}
}

// TestTenant_NilCorps 测试 nil 企业列表的情况
func TestTenant_NilCorps(t *testing.T) {
	// 创建 nil 企业列表的租户
	tenant := &Tenant{
		Corps: nil,
	}

	// 测试各种方法
	if corp := tenant.GetCorpByID("corp1"); corp != nil {
		t.Errorf("GetCorpByID() returned non-nil for nil corps: %v", corp)
	}

	if activeCorps := tenant.GetActiveCorps(); len(activeCorps) != 0 {
		t.Errorf("GetActiveCorps() returned %d corps, want 0", len(activeCorps))
	}

	if totalUsers := tenant.GetTotalUsers(); totalUsers != 0 {
		t.Errorf("GetTotalUsers() = %d, want 0", totalUsers)
	}

	if totalMessages := tenant.GetTotalMessages(); totalMessages != 0 {
		t.Errorf("GetTotalMessages() = %d, want 0", totalMessages)
	}
}
