package utils

// Article 文章结构体
type Article struct {
	Title          string        // 文章名称
	Status         string        // 文章状态（等待中，上传中，完成）
	MarkdownTool   MarkdownTool  // Markdown分析工具，分析的文章在这里
	Progress       float32       // 运行进度
	SelectAccounts []AccountInfo // 平台信息
}

// AccountInfo 平台信息
type AccountInfo struct {
	ID            uint    // 账户ID
	Username      string  // 平台名称
	PlatformKey   string  // 平台KEy
	PlatformAlias string  // 平台名称
	Status        string  // 状态
	StepCount     int     // 进度统计
	Progress      float32 // 运行进度(百分比)
	PublishURL    string  // 发布链接

}
