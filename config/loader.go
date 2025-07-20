// Package config 配置加载器实现
// 提供配置文件的加载、解析、验证和监控功能
// @Author lanpang
// @Date 2024/12/19
// @Desc 配置文件加载逻辑，支持TOML格式，包含配置验证和错误处理
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"vhagar/errors"
	"vhagar/logger"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
	configPath   string            // 配置文件路径
	config       *AppConfig        // 当前配置
	watcher      *fsnotify.Watcher // 文件监控器
	callbacks    []ReloadCallback  // 重载回调函数列表
	mu           sync.RWMutex      // 读写锁
	lastModified time.Time         // 最后修改时间
	isWatching   bool              // 是否正在监控
}

// ReloadCallback 配置重载回调函数类型
type ReloadCallback func(oldConfig, newConfig *AppConfig) error

// LoaderOptions 加载器选项
type LoaderOptions struct {
	EnableWatch    bool          // 是否启用文件监控
	WatchInterval  time.Duration // 监控间隔
	ValidateConfig bool          // 是否验证配置
	BackupConfig   bool          // 是否备份配置
	BackupDir      string        // 备份目录
}

// DefaultLoaderOptions 返回默认的加载器选项
func DefaultLoaderOptions() LoaderOptions {
	return LoaderOptions{
		EnableWatch:    true,
		WatchInterval:  time.Second,
		ValidateConfig: true,
		BackupConfig:   false,
		BackupDir:      "./config_backup",
	}
}

// NewConfigLoader 创建新的配置加载器
func NewConfigLoader(configPath string, options ...LoaderOptions) (*ConfigLoader, error) {
	opts := DefaultLoaderOptions()
	if len(options) > 0 {
		opts = options[0]
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NewWithDetail(
				errors.ErrCodeConfigNotFound,
				"配置文件不存在",
				fmt.Sprintf("路径: %s", configPath),
			)
		}
		return nil, errors.Wrap(errors.ErrCodeConfigInvalid, "检查配置文件失败", err)
	}

	loader := &ConfigLoader{
		configPath: configPath,
		callbacks:  make([]ReloadCallback, 0),
	}

	// 初始加载配置
	config, err := loader.loadConfigFromFile(configPath, opts.ValidateConfig)
	if err != nil {
		return nil, err
	}

	loader.config = config
	loader.lastModified = loader.getFileModTime()

	// 备份配置文件
	if opts.BackupConfig {
		if err := loader.backupConfig(opts.BackupDir); err != nil {
			log := logger.GetLogger()
			log.Warnw("备份配置文件失败", "error", err)
		}
	}

	// 启用文件监控
	if opts.EnableWatch {
		if err := loader.startWatching(); err != nil {
			log := logger.GetLogger()
			log.Warnw("启动配置文件监控失败", "error", err)
		}
	}

	return loader, nil
}

// LoadConfig 加载配置文件（静态方法）
func LoadConfig(configPath string) (*AppConfig, error) {
	loader, err := NewConfigLoader(configPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		return nil, err
	}
	defer loader.Close()

	return loader.GetConfig(), nil
}

// LoadConfigWithValidation 加载并验证配置文件
func LoadConfigWithValidation(configPath string) (*AppConfig, error) {
	loader := &ConfigLoader{}
	return loader.loadConfigFromFile(configPath, true)
}

// loadConfigFromFile 从文件加载配置
func (cl *ConfigLoader) loadConfigFromFile(configPath string, validate bool) (*AppConfig, error) {
	log := logger.GetLogger()
	log.Infow("开始加载配置文件", "path", configPath)

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(configPath))
	if ext != ".toml" {
		return nil, errors.NewWithDetail(
			errors.ErrCodeConfigInvalid,
			"不支持的配置文件格式",
			fmt.Sprintf("仅支持 .toml 格式，当前: %s", ext),
		)
	}

	// 读取并解析配置文件
	config := &AppConfig{}
	metadata, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, errors.WrapWithDetail(
			errors.ErrCodeConfigInvalid,
			"解析配置文件失败",
			fmt.Sprintf("文件: %s", configPath),
			err,
		)
	}

	// 检查未定义的配置项
	if undecoded := metadata.Undecoded(); len(undecoded) > 0 {
		log.Warnw("发现未识别的配置项", "keys", undecoded)
	}

	// 设置默认值
	cl.setDefaults(config)

	// 验证配置
	if validate {
		if err := config.Validate(); err != nil {
			return nil, errors.Wrap(errors.ErrCodeConfigInvalid, "配置验证失败", err)
		}
	}

	log.Infow("配置文件加载成功",
		"path", configPath,
		"project", config.ProjectName,
		"version", VERSION,
	)

	return config, nil
}

// setDefaults 设置配置的默认值
func (cl *ConfigLoader) setDefaults(config *AppConfig) {
	// 设置全局配置默认值
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.ProjectName == "" {
		config.ProjectName = "vhagar"
	}
	if config.Interval == 0 {
		config.Interval = 5 * time.Minute
	}
	if config.Duration == 0 {
		config.Duration = time.Hour
	}

	// 设置监控配置默认值
	if config.Metric.Port == "" {
		config.Metric.Port = "8090"
	}

	// 设置AI配置默认值
	if config.Services.AI.Enable && len(config.Services.AI.Providers) > 0 {
		for name, provider := range config.Services.AI.Providers {
			if provider.Model == "" {
				// 根据提供商设置默认模型
				switch strings.ToLower(name) {
				case "openai":
					provider.Model = "gpt-3.5-turbo"
				case "openrouter":
					provider.Model = "deepseek/deepseek-chat"
				case "gemini":
					provider.Model = "gemini-pro"
				default:
					provider.Model = "default"
				}
				config.Services.AI.Providers[name] = provider
			}
		}
	}
}

// GetConfig 获取当前配置（线程安全）
func (cl *ConfigLoader) GetConfig() *AppConfig {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	// 返回配置的深拷贝以避免并发修改
	return cl.copyConfig(cl.config)
}

// copyConfig 创建配置的深拷贝
func (cl *ConfigLoader) copyConfig(config *AppConfig) *AppConfig {
	// 这里简化处理，实际项目中可能需要更完整的深拷贝
	newConfig := *config
	return &newConfig
}

// ReloadConfig 重新加载配置
func (cl *ConfigLoader) ReloadConfig() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	log := logger.GetLogger()
	log.Info("开始重新加载配置文件")

	// 加载新配置
	newConfig, err := cl.loadConfigFromFile(cl.configPath, true)
	if err != nil {
		log.Errorw("重新加载配置失败", "error", err)
		return err
	}

	// 保存旧配置用于回调
	oldConfig := cl.config

	// 执行重载回调
	for _, callback := range cl.callbacks {
		if err := callback(oldConfig, newConfig); err != nil {
			log.Errorw("配置重载回调执行失败", "error", err)
			return errors.Wrap(errors.ErrCodeInternalErr, "配置重载回调失败", err)
		}
	}

	// 更新配置
	cl.config = newConfig
	cl.lastModified = cl.getFileModTime()

	log.Info("配置文件重新加载成功")
	return nil
}

// AddReloadCallback 添加配置重载回调函数
func (cl *ConfigLoader) AddReloadCallback(callback ReloadCallback) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.callbacks = append(cl.callbacks, callback)
}

// startWatching 开始监控配置文件变化
func (cl *ConfigLoader) startWatching() error {
	if cl.isWatching {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(errors.ErrCodeInternalErr, "创建文件监控器失败", err)
	}

	cl.watcher = watcher
	cl.isWatching = true

	// 添加配置文件到监控列表
	if err := watcher.Add(cl.configPath); err != nil {
		watcher.Close()
		cl.isWatching = false
		return errors.Wrap(errors.ErrCodeInternalErr, "添加文件监控失败", err)
	}

	// 启动监控协程
	go cl.watchLoop()

	log := logger.GetLogger()
	log.Infow("配置文件监控已启动", "path", cl.configPath)
	return nil
}

// watchLoop 文件监控循环
func (cl *ConfigLoader) watchLoop() {
	log := logger.GetLogger()

	for {
		select {
		case event, ok := <-cl.watcher.Events:
			if !ok {
				return
			}

			// 只处理写入和重命名事件
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Rename == fsnotify.Rename {

				// 检查文件修改时间，避免重复处理
				currentModTime := cl.getFileModTime()
				if currentModTime.After(cl.lastModified) {
					log.Infow("检测到配置文件变化", "event", event.Op.String())

					// 延迟一点时间，确保文件写入完成
					time.Sleep(100 * time.Millisecond)

					if err := cl.ReloadConfig(); err != nil {
						log.Errorw("自动重载配置失败", "error", err)
					}
				}
			}

		case err, ok := <-cl.watcher.Errors:
			if !ok {
				return
			}
			log.Errorw("配置文件监控错误", "error", err)
		}
	}
}

// getFileModTime 获取文件修改时间
func (cl *ConfigLoader) getFileModTime() time.Time {
	if info, err := os.Stat(cl.configPath); err == nil {
		return info.ModTime()
	}
	return time.Time{}
}

// backupConfig 备份配置文件
func (cl *ConfigLoader) backupConfig(backupDir string) error {
	// 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return errors.Wrap(errors.ErrCodeInternalErr, "创建备份目录失败", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Base(cl.configPath)
	backupPath := filepath.Join(backupDir, fmt.Sprintf("%s.%s.backup", filename, timestamp))

	// 复制文件
	sourceData, err := os.ReadFile(cl.configPath)
	if err != nil {
		return errors.Wrap(errors.ErrCodeInternalErr, "读取配置文件失败", err)
	}

	if err := os.WriteFile(backupPath, sourceData, 0644); err != nil {
		return errors.Wrap(errors.ErrCodeInternalErr, "写入备份文件失败", err)
	}

	log := logger.GetLogger()
	log.Infow("配置文件备份成功", "backup_path", backupPath)
	return nil
}

// ValidateConfigFile 验证配置文件而不加载
func ValidateConfigFile(configPath string) error {
	loader := &ConfigLoader{}
	_, err := loader.loadConfigFromFile(configPath, true)
	return err
}

// GetConfigInfo 获取配置文件信息
func (cl *ConfigLoader) GetConfigInfo() map[string]interface{} {
	info := make(map[string]interface{})

	if stat, err := os.Stat(cl.configPath); err == nil {
		info["path"] = cl.configPath
		info["size"] = stat.Size()
		info["mod_time"] = stat.ModTime()
		info["is_watching"] = cl.isWatching
	}

	config := cl.GetConfig()
	info["project_name"] = config.ProjectName
	info["log_level"] = config.LogLevel
	info["version"] = VERSION

	return info
}

// Close 关闭配置加载器
func (cl *ConfigLoader) Close() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	if cl.watcher != nil {
		cl.isWatching = false
		return cl.watcher.Close()
	}

	return nil
}

// ExportConfig 导出配置到文件
func (cl *ConfigLoader) ExportConfig(outputPath string) error {
	config := cl.GetConfig()

	file, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(errors.ErrCodeInternalErr, "创建导出文件失败", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return errors.Wrap(errors.ErrCodeInternalErr, "编码配置失败", err)
	}

	log := logger.GetLogger()
	log.Infow("配置导出成功", "output_path", outputPath)
	return nil
}
