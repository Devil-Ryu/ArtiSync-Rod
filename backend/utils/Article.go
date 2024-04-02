package utils

// Article 文章结构体
type Article struct {
	Title           string         // 文章名称
	Status          string         // 文章状态（等待中，上传中，完成）
	MarkdownTool    MarkdownTool   // Markdown分析工具，分析的文章在这里
	Progress        float32        // 运行进度
	SelectPlatforms []string       // 选择上传到平台
	PlatformsInfo   []PlatformInfo // 平台信息
}

// PlatformInfo 平台信息
type PlatformInfo struct {
	Name       string  // 名称
	Status     string  // 状态
	StepCount  int     // 进度统计
	Progress   float32 // 运行进度(百分比)
	PublishURL string  // 发布链接

}
