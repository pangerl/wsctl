package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/errors"
	"vhagar/logger"
)

// 和风天气 API 响应结构体（只保留常用字段）
type QWeatherNowResp struct {
	Code string `json:"code"`
	Now  struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Text      string `json:"text"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Vis       string `json:"vis"`
		Cloud     string `json:"cloud"`
	} `json:"now"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

type QWeatherCityLookupResp struct {
	Code     string `json:"code"`
	Location []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Adm1    string `json:"adm1"`
		Adm2    string `json:"adm2"`
		Country string `json:"country"`
	} `json:"location"`
}

// 校验 location 是否为 LocationID 或 经纬度
func isLocationIDOrLatLon(location string) bool {
	// LocationID: 全数字，通常为6位以上
	idRe := regexp.MustCompile(`^\d{6,}$`)
	// 经纬度: 116.41,39.92
	latlonRe := regexp.MustCompile(`^-?\d{1,3}\.\d{1,6},-?\d{1,3}\.\d{1,6}$`)
	return idRe.MatchString(location) || latlonRe.MatchString(location)
}

// 查询城市名称对应的 LocationID
func queryLocationID(apiHost, apiKey, city string, client *http.Client) (string, error) {
	urlStr := fmt.Sprintf("%s/geo/v2/city/lookup?location=%s", strings.TrimRight(apiHost, "/"), url.QueryEscape(city))

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", errors.WrapError(errors.ErrCodeNetworkFailed, "创建城市查询请求失败", err)
	}

	req.Header.Set("X-QW-Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.WrapError(errors.ErrCodeNetworkFailed, "城市查询请求失败", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WrapError(errors.ErrCodeNetworkFailed, "读取城市查询响应失败", err)
	}

	var lookup QWeatherCityLookupResp
	if err := json.Unmarshal(body, &lookup); err != nil {
		return "", errors.WrapError(errors.ErrCodeAIResponseInvalid, "城市查询响应解析失败", err)
	}

	if lookup.Code != "200" || len(lookup.Location) == 0 {
		return "", errors.NewErrorWithDetail(errors.ErrCodeNotFound, "城市查询失败",
			fmt.Sprintf("code=%s, body=%s", lookup.Code, string(body)))
	}

	return lookup.Location[0].ID, nil
}

// 查询实时天气
func queryQWeatherNow(apiHost, apiKey, location string, client *http.Client) (*QWeatherNowResp, error) {
	urlStr := fmt.Sprintf("%s/v7/weather/now?location=%s", strings.TrimRight(apiHost, "/"), url.QueryEscape(location))

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeNetworkFailed, "创建天气查询请求失败", err)
	}

	req.Header.Set("X-QW-Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeNetworkFailed, "天气查询请求失败", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WrapError(errors.ErrCodeNetworkFailed, "读取天气查询响应失败", err)
	}

	var now QWeatherNowResp
	if err := json.Unmarshal(body, &now); err != nil {
		return nil, errors.WrapError(errors.ErrCodeAIResponseInvalid, "天气响应解析失败", err)
	}

	if now.Code != "200" {
		return nil, errors.NewErrorWithDetail(errors.ErrCodeNotFound, "天气查询失败",
			fmt.Sprintf("code=%s, body=%s", now.Code, string(body)))
	}

	return &now, nil
}

// CallWeatherTool 天气查询工具入口
func CallWeatherTool(ctx context.Context, args map[string]any) (string, error) {
	// 参数验证
	location, ok := args["location"].(string)
	if !ok || location == "" {
		err := errors.NewError(errors.ErrCodeInvalidParam, "缺少或无效的 location 参数")
		errors.LogError(err, "天气工具调用")
		return "", err
	}

	// 配置验证
	if config.Config == nil {
		err := errors.NewError(errors.ErrCodeConfigNotFound, "系统配置未初始化")
		errors.LogError(err, "天气工具调用")
		return "", err
	}

	apiHost := config.Config.Services.Weather.ApiHost
	apiKey := config.Config.Services.Weather.ApiKey

	if apiHost == "" || apiKey == "" {
		err := errors.NewError(errors.ErrCodeConfigInvalid, "天气API主机或密钥未配置")
		errors.LogError(err, "天气工具调用")
		return "", err
	}

	logger.Logger.Infow("开始查询天气", "location", location)

	client := &http.Client{Timeout: 8 * time.Second}
	var locationID string
	var err error

	// 判断输入类型并获取LocationID
	if isLocationIDOrLatLon(location) {
		locationID = location
		logger.Logger.Infow("使用LocationID或经纬度", "location_id", locationID)
	} else {
		locationID, err = queryLocationID(apiHost, apiKey, location, client)
		if err != nil {
			errors.LogErrorWithFields(err, "城市ID查询", map[string]interface{}{
				"city": location,
			})
			return "", err
		}
		logger.Logger.Infow("城市ID查询成功", "city", location, "location_id", locationID)
	}

	// 查询天气
	weather, err := queryQWeatherNow(apiHost, apiKey, locationID, client)
	if err != nil {
		errors.LogErrorWithFields(err, "天气查询", map[string]interface{}{
			"location_id": locationID,
		})
		return "", err
	}

	// 构造响应
	resp := map[string]any{
		"location":    location,
		"location_id": locationID,
		"weather":     weather,
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		appErr := errors.WrapError(errors.ErrCodeInternalErr, "结果序列化失败", err)
		errors.LogError(appErr, "天气工具调用")
		return "", appErr
	}

	logger.Logger.Infow("天气查询完成", "location", location, "temp", weather.Now.Temp)
	return string(jsonBytes), nil
}
