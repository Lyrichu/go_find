package finder

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
)

// Finder 定义文件查找器结构
// ANSI 颜色转义序列
const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

type Finder struct {
	pattern     *regexp.Regexp
	fileType    string
	workers     int // 并发工作协程数
	colorOutput bool // 是否启用颜色输出
}

// NewFinder 创建新的 Finder 实例
func NewFinder(pattern string, fileType string) *Finder {
	var reg *regexp.Regexp
	if pattern != "" {
		// 将通配符 * 转换为正则表达式的 .*
		regexPattern := ""
		for i := 0; i < len(pattern); i++ {
			if pattern[i] == '*' {
				regexPattern += ".*"
			} else {
				// 转义其他正则表达式的特殊字符
				regexPattern += regexp.QuoteMeta(string(pattern[i]))
			}
		}
		// 添加开始和结束标记以确保完整匹配
		reg = regexp.MustCompile("^" + regexPattern + "$")
	}

	// 检查是否支持颜色输出
	colorOutput := true
	if os.Getenv("TERM") == "dumb" || os.Getenv("NO_COLOR") != "" {
		colorOutput = false
	}

	return &Finder{
		pattern:     reg,
		fileType:    fileType,
		workers:     runtime.NumCPU(), // 默认使用CPU核心数作为工作协程数
		colorOutput: colorOutput,
	}
}

// Find 执行文件查找
// FormatResult 格式化并高亮显示匹配结果
func (f *Finder) FormatResult(path string) string {
	if f.pattern == nil {
		return path
	}

	base := filepath.Base(path)
	dir := filepath.Dir(path)

	// 高亮显示匹配的文件名
	highlightedBase := base
	if f.colorOutput {
		highlightedBase = f.pattern.ReplaceAllStringFunc(base, func(match string) string {
			return colorRed + match + colorReset
		})

		// 如果是目录，使用绿色显示
		if f.fileType == "d" {
			highlightedBase = colorGreen + highlightedBase + colorReset
		}
	}

	// 组合完整路径
	if dir == "." {
		return highlightedBase
	}
	return filepath.Join(dir, highlightedBase)
}

// Find 执行文件查找
func (f *Finder) Find(root string) ([]string, error) {
	// 创建任务和结果通道
	paths := make(chan string)
	results := make(chan string)
	done := make(chan error)
	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < f.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				info, err := os.Stat(path)
				if err != nil {
					continue
				}

				// 根据类型过滤
				if f.fileType == "f" && info.IsDir() {
					continue
				}
				if f.fileType == "d" && !info.IsDir() {
					continue
				}

				// 应用正则匹配
				if f.pattern != nil {
					if !f.pattern.MatchString(info.Name()) {
						continue
					}
				}

				// 发送匹配结果
				results <- path
			}
		}()
	}

	// 启动遍历协程
	go func() {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			paths <- path
			return nil
		})

		// 关闭paths通道并发送错误结果
		close(paths)
		done <- err
	}()

	// 收集结果
	var matchedResults []string
	go func() {
		wg.Wait()
		close(results)
	}()

	// 等待所有结果
	for {
		select {
		case err := <-done:
			if err != nil {
				return nil, err
			}
		case result, ok := <-results:
			if !ok {
				return matchedResults, nil
			}
			matchedResults = append(matchedResults, result)
		}
	}
}
