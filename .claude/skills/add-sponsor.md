---
skill: add-sponsor
description: Add a new sponsor to CCX project following the established integration pattern
triggerWords: [add sponsor, new sponsor, integrate sponsor]
tags: [sponsor, integration, documentation]
---

# Add Sponsor Skill

自动化添加新赞助商的完整集成流程。

## 用法

```
/add-sponsor <sponsor-name> <description> <urls>
```

或者交互式调用：
```
/add-sponsor
```

## 执行流程

本 skill 会按照既定的赞助商集成模式，自动完成以下任务：

### 1. 收集赞助商信息

通过交互式问答收集：
- 赞助商名称（中英文）
- 简短介绍（中英文）
- 详细描述（中英文）
- Base URL
- 控制台 URL
- 推广链接（可选）
- 图标文件路径
- 赞助商顺序位置（在哪个赞助商之后插入）

### 2. 更新 README 文件

在两个 README 文件中按照赞助商顺序添加：
- `README.md` - 英文版本
- `README.zh-CN.md` - 中文版本

**位置**: 
- 查找现有赞助商列表
- 根据指定的顺序位置插入新赞助商

**格式**:
```markdown
### [赞助商名称](推广链接)

<div align="center">
  <img src="docs/sponsors/sponsor-id.jpg" alt="Sponsor Logo" width="120"/>
</div>

赞助商详细描述...
```

### 3. 保存图标文件

- 复制图标到 `docs/sponsors/<sponsor-id>.jpg`
- 复制图标到 `desktop/frontend/src/assets/<sponsor-id>.jpg`

### 4. 前端 Web UI 集成

#### 4.1 快速识别配置
**文件**: `frontend/src/utils/quickInputParser.ts`

在 `knownOpenAIUrls` 或 `knownClaudeUrls` 中按顺序添加 Base URL：
```typescript
const knownOpenAIUrls = new Set([
  // ... 现有 URL
  'https://previous-sponsor.com/v1',
  'https://new-sponsor.com/v1',  // 新增
  'https://next-sponsor.com/v1',
])
```

### 5. 桌面端 Agent 配置集成

#### 5.1 前端类型定义
**文件**: `desktop/frontend/src/types/index.ts`

添加到 `AgentProvider` 类型：
```typescript
export type AgentProvider = '...' | 'new-sponsor' | '...'
```

#### 5.2 Provider 表单
**文件**: `desktop/frontend/src/components/agent/ProviderForm.vue`

在下拉选项中按顺序添加：
```vue
<option value="previous-sponsor">{{ t('agent.provider.previousDirect') }}</option>
<option value="new-sponsor">{{ t('agent.provider.newSponsorDirect') }}</option>
<option value="next-sponsor">{{ t('agent.provider.nextDirect') }}</option>
```

并添加图标导入：
```typescript
import newSponsorIcon from '@/assets/new-sponsor.jpg'

const providerIcons: Record<string, string> = {
  'new-sponsor': newSponsorIcon,
}
```

#### 5.3 Agent 配置逻辑
**文件**: `desktop/frontend/src/composables/useAgentConfig.ts`

在以下位置按顺序添加 `newSponsor`:

1. `claudeProviderLabels` - 添加标签
2. `codexProviderLabels` 的 computed - 添加标签
3. `claudeProviderKeys` - 添加空字符串
4. `isClaudeProvider` - 添加判断条件
5. `isCodexThirdPartyWithMode` - 添加判断条件
6. `claudeTargetBaseUrl()` - 添加 case 语句
7. `codexTargetBaseUrl()` - 添加 case 语句
8. `openCodeTargetBaseUrl()` - 添加 case 语句

**示例**:
```typescript
const claudeProviderLabels: Record<AgentProvider | 'custom', string> = {
  // ...
  'previous-sponsor': 'Previous Sponsor',
  'new-sponsor': 'New Sponsor',
  'next-sponsor': 'Next Sponsor',
}

const claudeTargetBaseUrl = () => {
  switch (selectedClaudeProvider.value) {
    case 'new-sponsor':
      return 'https://new-sponsor.com/v1'
    // ...
  }
}
```

#### 5.4 外部链接配置
**文件**: `desktop/frontend/src/lib/external-link.ts`

添加控制台和推广链接：
```typescript
export const providerConsoleLinks: Record<string, string> = {
  // ...
  'new-sponsor': 'https://new-sponsor.com/dashboard',
}

export const providerPromotionLinks: Record<string, string> = {
  // ...
  'new-sponsor': 'https://new-sponsor.com/register?aff=ccx',
}
```

#### 5.5 国际化翻译
**文件**: 
- `desktop/frontend/src/locales/zh-CN.json`
- `desktop/frontend/src/locales/en.json`

添加翻译键：
```json
{
  "agent.provider.newSponsorDirect": "New Sponsor 直连"
}
```

### 6. 桌面端渠道中心集成

#### 6.1 渠道顺序配置
**文件**: `desktop/frontend/src/components/channel/ChannelTab.vue`

在 `presetOrder` 数组中按顺序添加：
```typescript
const presetOrder = [
  // ...
  'previous-sponsor',
  'new-sponsor',
  'next-sponsor',
  // ...
]
```

### 7. 后端 Go 代码集成

#### 7.1 配置服务
**文件**: `desktop/internal/configservice/service.go`

1. 添加 Provider 常量：
```go
const (
    // ...
    ProviderPreviousSponsor = "previous-sponsor"
    ProviderNewSponsor      = "new-sponsor"
    ProviderNextSponsor     = "next-sponsor"
)
```

2. 添加 Base URL 常量：
```go
const (
    // ...
    newSponsorBaseURL = "https://new-sponsor.com/v1"
)
```

3. 在所有 switch/case 和列表中按顺序添加处理逻辑

#### 7.2 渠道预设
**文件**: `desktop/internal/channelpreset/preset.go`

1. 添加 Provider 常量：
```go
const (
    // ...
    ProviderNewSponsor = "new-sponsor"
)
```

2. 添加到 `providerConsoleURLs`：
```go
var providerConsoleURLs = map[string]string{
    // ...
    ProviderNewSponsor: "https://new-sponsor.com/dashboard",
}
```

3. 添加 Preset 配置：
```go
{
    ID:                  ProviderNewSponsor,
    Order:               45,
    Label:               "New Sponsor",
    Description:         "赞助商详细描述...",
    DirectAgent:         true,
    NativeMessages:      true,
    ChatCompatible:      true,
    ResponsesCompatible: true,
    Plans: []ProviderPlan{
        {ID: "openai-chat", Label: "OpenAI-compatible", BaseURL: "https://new-sponsor.com/v1", Description: "OpenAI Chat 兼容入口", Recommended: true},
    },
    Targets: []ChannelTarget{
        {Type: TargetMessages, Label: "Messages 原生透传", Description: "Claude Code 直连或 CCX messages 渠道", Recommended: true},
        {Type: TargetResponses, Label: "Codex Responses", Description: "OpenAI Responses 协议，供 Codex 使用"},
        {Type: TargetChat, Label: "Chat 渠道透传", Description: "OpenAI Chat 协议，供 Chat 客户端使用"},
    },
    DefaultTarget: TargetMessages,
}
```

4. 在配置映射中添加空配置：
```go
TargetMessages: {
    ProviderNewSponsor: {},
},
TargetChat: {
    ProviderNewSponsor: {},
},
TargetResponses: {
    ProviderNewSponsor: {
        CodexToolCompat:       boolRef(false),
        StripCodexClientTools: boolRef(false),
    },
}
```

### 8. Git 提交流程

创建功能分支并提交：
```bash
git checkout -b feat/add-<sponsor-id>-sponsor
git add <所有修改的文件>
git commit -m "feat(sponsors): add <SponsorName> sponsor integration"
git push -u origin feat/add-<sponsor-id>-sponsor
```

提交信息格式：
```
feat(sponsors): add <SponsorName> sponsor integration

- Add <SponsorName> to README files (en & zh-CN)
- Add <SponsorName> to channel quick input parser
- Add <SponsorName> to desktop agent configuration
- Add <SponsorName> to desktop channel center presets
- Add <SponsorName> icon to assets
- Follow sponsor order: <Previous> → <SponsorName> → <Next>
```

## 检查清单

执行完成后，确认以下所有项都已完成：

### 文档
- [ ] README.md 已更新
- [ ] README.zh-CN.md 已更新
- [ ] 图标已保存到 docs/sponsors/
- [ ] 图标已保存到 desktop/frontend/src/assets/

### 前端 Web UI
- [ ] quickInputParser.ts 已添加 URL

### 桌面端前端
- [ ] types/index.ts 已添加类型
- [ ] ProviderForm.vue 已添加选项和图标
- [ ] useAgentConfig.ts 已完整集成（8处）
- [ ] external-link.ts 已添加链接
- [ ] locales/zh-CN.json 已添加翻译
- [ ] locales/en.json 已添加翻译
- [ ] ChannelTab.vue 已添加排序

### 桌面端后端
- [ ] configservice/service.go 已完整集成
- [ ] channelpreset/preset.go 已完整集成

### Git
- [ ] 已创建功能分支
- [ ] 已提交所有更改
- [ ] 提交信息格式正确
- [ ] 已推送到远程仓库

## 注意事项

1. **严格遵守赞助商顺序**：所有文件中的顺序必须一致
2. **Targets 顺序固定**：必须按照 **Messages → Responses → Chat** 的顺序排列，不能颠倒
3. **图标格式**：建议使用 JPG 格式，尺寸约 120x120 像素
4. **URL 格式**：确保所有 URL 格式统一，去除尾部斜杠
5. **国际化**：中英文翻译都要提供
6. **Base URL 选择**：根据赞助商主要支持的协议选择合适的默认 Target
7. **Order 值**：按 10 的倍数递增，便于后续插入

## 常见问题

**Q: 如何确定赞助商顺序？**
A: 参考现有赞助商的 Order 值，在合适的位置插入。通常按重要性和合作时间排序。

**Q: 如果赞助商同时支持多种协议怎么办？**
A: 在 Plans 中添加多个选项，并设置 Recommended 标记主推方案。

**Q: 如何测试集成是否成功？**
A: 
1. 启动桌面应用，检查渠道中心是否显示新卡片
2. 在 Agent 配置中检查是否能选择新的直连选项
3. 在 Web UI 粘贴 URL 检查是否自动识别

## 参考实例

参考 Unity2.ai 的完整集成（commit: 2d962b7f）：
- 完整的文件修改列表
- 所有配置点的具体实现
- 提交信息格式示例
