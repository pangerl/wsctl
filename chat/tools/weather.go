package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"vhagar/config"
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
	// LocationID: 全数字，通常为9位
	idRe := regexp.MustCompile(`^\d{6,}$`)
	// 经纬度: 116.41,39.92
	latlonRe := regexp.MustCompile(`^-?\d{1,3}\.\d{1,6},-?\d{1,3}\.\d{1,6}$`)
	return idRe.MatchString(location) || latlonRe.MatchString(location)
}

// 查询城市名称对应的 LocationID
func queryLocationID(apiHost, apiKey, city string, client *http.Client) (string, error) {
	urlStr := fmt.Sprintf("%s/geo/v2/city/lookup?location=%s", strings.TrimRight(apiHost, "/"), url.QueryEscape(city))
	// log.Printf("[INFO] 查询城市LocationID: %s", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	// 新认证方式
	req.Header.Set("X-QW-Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("城市查询请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var lookup QWeatherCityLookupResp
	if err := json.Unmarshal(body, &lookup); err != nil {
		return "", fmt.Errorf("城市查询响应解析失败: %w", err)
	}
	if lookup.Code != "200" || len(lookup.Location) == 0 {
		return "", fmt.Errorf("城市查询失败，code=%s, body=%s", lookup.Code, string(body))
	}
	return lookup.Location[0].ID, nil
}

// 查询实时天气
func queryQWeatherNow(apiHost, apiKey, location string, client *http.Client) (*QWeatherNowResp, error) {
	urlStr := fmt.Sprintf("%s/v7/weather/now?location=%s", strings.TrimRight(apiHost, "/"), url.QueryEscape(location))
	// log.Printf("[INFO] 查询天气: %s", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	// 新认证方式
	req.Header.Set("X-QW-Api-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("天气查询请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var now QWeatherNowResp
	if err := json.Unmarshal(body, &now); err != nil {
		return nil, fmt.Errorf("天气响应解析失败: %w", err)
	}
	if now.Code != "200" {
		return nil, fmt.Errorf("天气查询失败，code=%s, body=%s", now.Code, string(body))
	}
	return &now, nil
}

func CallWeatherTool(ctx context.Context, args map[string]any) (string, error) {
	location, ok := args["location"].(string)
	if !ok || location == "" {
		return "", errors.New("缺少或无效的 location 参数")
	}

	apiHost := config.Config.Weather.ApiHost
	apiKey := config.Config.Weather.ApiKey

	if apiHost == "" || apiKey == "" {
		return "", errors.New("天气API主机或密钥未配置")
	}
	client := &http.Client{Timeout: 8 * time.Second}
	var locationID string
	var err error
	if isLocationIDOrLatLon(location) {
		locationID = location
	} else {
		locationID, err = queryLocationID(apiHost, apiKey, location, client)
		if err != nil {
			log.Printf("[ERROR] 城市查ID失败: %v", err)
			return fmt.Sprintf("城市查ID失败: %v", err), nil
		}
	}

	weather, err := queryQWeatherNow(apiHost, apiKey, locationID, client)
	if err != nil {
		log.Printf("[ERROR] 天气查询失败: %v", err)
		return fmt.Sprintf("天气查询失败: %v", err), nil
	}
	resp := map[string]any{
		"location":    location,
		"location_id": locationID,
		"weather":     weather,
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[ERROR] 结果序列化失败: %v", err)
		return fmt.Sprintf("结果序列化失败: %v", err), nil
	}
	// log.Printf("[TOOL] get_forecast 返回: %s", string(jsonBytes))
	return string(jsonBytes), nil
}
