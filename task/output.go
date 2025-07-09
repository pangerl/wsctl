package task

import (
	"io"
	"os"
	"sync"
)

var (
	outputWriter io.Writer
	outputFile   *os.File
	once         sync.Once
)

const outputFileName = "task_output.log" // 固定文件名

// GetOutputWriter 返回全局唯一的 io.Writer，写入终端和文件
func GetOutputWriter() io.Writer {
	once.Do(func() {
		file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// 打开失败只写终端
			outputWriter = os.Stdout
			return
		}
		outputFile = file
		outputWriter = io.MultiWriter(os.Stdout, file)
	})
	return outputWriter
}

// CloseOutputFile 关闭文件，建议在 main 退出时调用
func CloseOutputFile() {
	if outputFile != nil {
		_ = outputFile.Close()
	}
}

// ClearOutputFile 清空日志文件内容
func ClearOutputFile() error {
	// 关闭当前文件句柄
	if outputFile != nil {
		_ = outputFile.Close()
		outputFile = nil
		outputWriter = nil
	}
	// 以截断方式重新打开
	file, err := os.OpenFile(outputFileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
