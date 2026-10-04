package boot

// demoCategory describes a seeded category together with its sites.
type demoCategory struct {
	Name        string
	Icon        string
	Description string
	Sites       []demoSite
}

// demoSite describes a single seeded navigation entry.
type demoSite struct {
	Name        string
	Url         string
	Icon        string
	Description string
}

// demoCategories is only used on a completely empty installation.
var demoCategories = []demoCategory{
	{
		Name:        "开发工具",
		Icon:        "🛠️",
		Description: "日常编码、调试与协作平台",
		Sites: []demoSite{
			{Name: "GitHub", Url: "https://github.com", Icon: "https://github.com/favicon.ico", Description: "全球最大的代码托管与协作平台"},
			{Name: "Go 官方文档", Url: "https://go.dev/doc/", Icon: "https://go.dev/favicon.ico", Description: "Go 语言官方文档与规范"},
			{Name: "GoFrame", Url: "https://goframe.org", Icon: "https://goframe.org/favicon.ico", Description: "本导航页所使用的 Go 应用开发框架"},
			{Name: "MDN Web Docs", Url: "https://developer.mozilla.org", Icon: "https://developer.mozilla.org/favicon.ico", Description: "最权威的 Web 技术文档"},
			{Name: "Stack Overflow", Url: "https://stackoverflow.com", Icon: "https://stackoverflow.com/favicon.ico", Description: "程序员问答社区"},
			{Name: "Docker Hub", Url: "https://hub.docker.com", Icon: "https://hub.docker.com/favicon.ico", Description: "容器镜像仓库"},
		},
	},
	{
		Name:        "AI 工具",
		Icon:        "🤖",
		Description: "大模型对话、绘画与模型社区",
		Sites: []demoSite{
			{Name: "ChatGPT", Url: "https://chat.openai.com", Icon: "https://chat.openai.com/favicon.ico", Description: "OpenAI 对话助手"},
			{Name: "Claude", Url: "https://claude.ai", Icon: "https://claude.ai/favicon.ico", Description: "Anthropic 出品的 AI 助手"},
			{Name: "Hugging Face", Url: "https://huggingface.co", Icon: "https://huggingface.co/favicon.ico", Description: "开源模型与数据集社区"},
			{Name: "Perplexity", Url: "https://www.perplexity.ai", Icon: "https://www.perplexity.ai/favicon.ico", Description: "AI 驱动的答案搜索引擎"},
		},
	},
	{
		Name:        "设计资源",
		Icon:        "🎨",
		Description: "界面设计、图标与图片素材",
		Sites: []demoSite{
			{Name: "Figma", Url: "https://www.figma.com", Icon: "https://www.figma.com/favicon.ico", Description: "在线协作界面设计工具"},
			{Name: "Iconify", Url: "https://icon-sets.iconify.design", Icon: "https://iconify.design/favicon.ico", Description: "海量开源图标集合"},
			{Name: "Dribbble", Url: "https://dribbble.com", Icon: "https://dribbble.com/favicon.ico", Description: "设计师作品展示社区"},
			{Name: "Unsplash", Url: "https://unsplash.com", Icon: "https://unsplash.com/favicon.ico", Description: "免费高质量图片素材"},
			{Name: "Google Fonts", Url: "https://fonts.google.com", Icon: "https://fonts.google.com/favicon.ico", Description: "免费开源网页字体"},
		},
	},
	{
		Name:        "效率办公",
		Icon:        "⚡",
		Description: "笔记、文档与团队协作",
		Sites: []demoSite{
			{Name: "Notion", Url: "https://www.notion.so", Icon: "https://www.notion.so/favicon.ico", Description: "一体化笔记与知识库"},
			{Name: "语雀", Url: "https://www.yuque.com", Icon: "https://www.yuque.com/favicon.ico", Description: "专业的云端知识库"},
			{Name: "飞书", Url: "https://www.feishu.cn", Icon: "https://www.feishu.cn/favicon.ico", Description: "先进企业协作与管理平台"},
			{Name: "ProcessOn", Url: "https://www.processon.com", Icon: "https://www.processon.com/favicon.ico", Description: "在线思维导图与流程图"},
		},
	},
}
