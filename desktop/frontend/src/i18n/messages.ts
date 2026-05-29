export type SupportedLocale = 'en' | 'zh-CN'

export const defaultLocale: SupportedLocale = 'en'

export const languageOptions: { locale: SupportedLocale; label: string }[] = [
  { locale: 'en', label: 'English' },
  { locale: 'zh-CN', label: '中文' },
]

export type MessageKey =
  | 'common.gatewayLabel'
  | 'common.online'
  | 'common.connecting'
  | 'common.offline'
  | 'common.refreshWebUI'
  | 'common.version'
  | 'common.gatewayPort'
  | 'common.daemonPid'
  | 'common.autoStart'
  | 'common.autoStartOn'
  | 'common.autoStartOff'
  | 'common.serviceHealthy'
  | 'common.serviceStarting'
  | 'common.serviceDisconnected'
  | 'common.settings'
  | 'common.save'
  | 'common.cancel'
  | 'common.retry'
  | 'nav.status'
  | 'nav.statusDesc'
  | 'nav.agent'
  | 'nav.agentDesc'
  | 'nav.channels'
  | 'nav.channelsDesc'
  | 'nav.env'
  | 'nav.envDesc'
  | 'nav.web'
  | 'nav.webDesc'
  | 'tab.statusTitle'
  | 'tab.agentTitle'
  | 'tab.channelsTitle'
  | 'tab.envTitle'
  | 'tab.webTitle'
  | 'sidebar.language'
  | 'sidebar.languageEnglish'
  | 'sidebar.languageChinese'
  | 'setup.loading'
  | 'setup.title'
  | 'setup.description'
  | 'setup.regenerate'
  | 'setup.regenerateTitle'
  | 'setup.hide'
  | 'setup.show'
  | 'setup.copied'
  | 'setup.copyKey'
  | 'setup.configPath'
  | 'setup.copyPath'
  | 'setup.hint'
  | 'setup.saving'
  | 'setup.submit'
  | 'webui.notRunning'
  | 'webui.openInBrowser'
  | 'metrics.gatewayPort'
  | 'metrics.uptime'
  | 'metrics.channels'
  | 'metrics.version'
  | 'actions.start'
  | 'actions.stop'
  | 'actions.restart'
  | 'actions.openWebUI'
  | 'actions.openBrowser'
  | 'actions.refreshStatus'
  | 'details.title'
  | 'details.binary'
  | 'details.binaryMissing'
  | 'details.dataDir'
  | 'details.dataDirMissing'
  | 'details.healthStatus'
  | 'details.revealDir'
  | 'details.openDir'
  | 'channel.headerEyebrow'
  | 'channel.title'
  | 'channel.description'
  | 'channel.hasKey'
  | 'channel.promo'
  | 'channel.console'
  | 'channel.target'
  | 'channel.keySavedPlaceholder'
  | 'channel.keyInputPlaceholder'
  | 'channel.missingKey'
  | 'channel.reuseKey'
  | 'channel.name'
  | 'channel.nameHint'
  | 'channel.presetWrites'
  | 'channel.capabilityHint'
  | 'channel.addToCCX'
  | 'channel.badgeDirectAgent'
  | 'channel.badgeNativeMessages'
  | 'env.title'
  | 'env.pathDetecting'
  | 'env.refresh'
  | 'env.openingEditor'
  | 'env.openWithEditor'
  | 'env.openInEditor'
  | 'env.selectEditor'
  | 'env.save'
  | 'env.saving'
  | 'env.saved'
  | 'env.failed'
  | 'env.hide'
  | 'env.show'
  | 'env.copied'
  | 'env.copy'
  | 'env.fieldRequired'
  | 'env.fieldDisallow'
  | 'env.fieldNumber'
  | 'env.fieldInteger'
  | 'env.fieldMin'
  | 'env.fieldMax'
  | 'env.loadFailed'
  | 'env.saveSuccess'
  | 'agent.statusDetecting'
  | 'agent.statusConfigured'
  | 'agent.statusPortMismatch'
  | 'agent.statusUnconfigured'
  | 'agent.localGateway'
  | 'agent.custom'
  | 'agent.currentProvider'
  | 'agent.currentUrl'
  | 'agent.targetUrl'
  | 'agent.notSet'
  | 'agent.configPath'
  | 'agent.authPath'
  | 'agent.openFileInEditor'
  | 'agent.applyConfig'
  | 'agent.restoreConfig'
  | 'agent.openConsole'
  | 'agent.codexPlaceholderSaved'
  | 'agent.codexPlaceholderRequired'
  | 'agent.codexPlaceholderWriteOnly'
  | 'agent.diffPreviewApply'
  | 'agent.diffPreviewRestore'
  | 'agent.diffConfirmApply'
  | 'agent.diffConfirmRestore'
  | 'agent.diffActionCreate'
  | 'agent.diffActionDelete'
  | 'agent.diffActionModify'
  | 'agent.diffComputing'
  | 'agent.diffNoChanges'
  | 'agent.diffExpandContext'
  | 'agent.diffCollapseContext'
  | 'agent.diffCancel'
  | 'agent.provider.localGateway'
  | 'agent.provider.deepseekDirect'
  | 'agent.provider.mimoDirect'
  | 'agent.provider.compshareDirect'
  | 'agent.provider.kimiDirect'
  | 'agent.provider.glmDirect'
  | 'agent.provider.minimaxDirect'
  | 'agent.provider.dashscopeDirect'
  | 'agent.provider.opencodeZenDirect'
  | 'agent.provider.opencodeGoDirect'
  | 'agent.provider.openaiDirect'
  | 'agent.promo'
  | 'agent.planPayAsYouGo'
  | 'agent.planChina'
  | 'agent.planSingapore'
  | 'agent.planEurope'
  | 'agent.planSubscription'
  | 'agent.billingModeMiMo'
  | 'agent.billingModeDashScope'
  | 'agent.placeholderSaved'
  | 'agent.placeholderMimo'
  | 'agent.placeholderDashScope'
  | 'agent.placeholderRequired'
  | 'env.groupAccess'
  | 'env.groupAccessDesc'
  | 'env.fieldProxyAccessKey'
  | 'env.placeholderProxyAccessKey'
  | 'env.fieldAdminAccessKey'
  | 'env.placeholderAdminAccessKey'
  | 'env.descAdminAccessKey'
  | 'env.groupServer'
  | 'env.groupServerDesc'
  | 'env.fieldPort'
  | 'env.descPort'
  | 'env.fieldEnv'
  | 'env.descEnv'
  | 'env.groupWebUI'
  | 'env.groupWebUIDesc'
  | 'env.fieldEnableWebUI'
  | 'env.descEnableWebUI'
  | 'env.fieldAppUILanguage'
  | 'env.groupLogs'
  | 'env.groupLogsDesc'
  | 'env.fieldLogLevel'
  | 'env.fieldEnableRequestLogs'
  | 'env.fieldEnableResponseLogs'
  | 'env.descEnableResponseLogs'
  | 'env.fieldQuietPollingLogs'
  | 'env.fieldRawLogOutput'
  | 'env.fieldSseDebugLevel'
  | 'env.fieldRewriteResponseModel'
  | 'env.groupPerformance'
  | 'env.groupPerformanceDesc'
  | 'env.fieldRequestTimeout'
  | 'env.fieldServerReadTimeout'
  | 'env.fieldMaxRequestBodySize'
  | 'env.fieldResponseHeaderTimeout'
  | 'env.groupCors'
  | 'env.groupCorsDesc'
  | 'env.fieldEnableCors'
  | 'env.fieldCorsOrigin'
  | 'env.groupCircuitBreaker'
  | 'env.groupCircuitBreakerDesc'
  | 'env.fieldMetricsWindowSize'
  | 'env.fieldMetricsFailureThreshold'
  | 'env.groupMetricsPersistence'
  | 'env.groupMetricsPersistenceDesc'
  | 'env.fieldMetricsPersistenceEnabled'
  | 'env.fieldMetricsRetentionDays'
  | 'logs.searchPlaceholder'
  | 'logs.autoScroll'
  | 'logs.copied'
  | 'logs.copyAll'
  | 'logs.clear'
  | 'logs.noSearchResults'
  | 'logs.empty'
  | 'diagnostic.binaryTitle'
  | 'diagnostic.binarySuggestionBuild'
  | 'diagnostic.binarySuggestionCheckDataDir'
  | 'diagnostic.binarySuggestionDownload'
  | 'diagnostic.portTitle'
  | 'diagnostic.portSuggestionInstance'
  | 'diagnostic.portSuggestionEnv'
  | 'diagnostic.portSuggestionInspect'
  | 'diagnostic.healthTitle'
  | 'diagnostic.healthSuggestionLogs'
  | 'diagnostic.healthSuggestionEnv'
  | 'diagnostic.healthSuggestionChannels'
  | 'diagnostic.healthSuggestionRestart'
  | 'diagnostic.permissionTitle'
  | 'diagnostic.permissionSuggestionDataDir'
  | 'diagnostic.permissionSuggestionExecutable'
  | 'diagnostic.permissionSuggestionWindows'
  | 'diagnostic.genericTitle'
  | 'diagnostic.genericSuggestionLogs'
  | 'diagnostic.genericSuggestionRestart'
  | 'setup.errorEmptyKey'
  | 'env.saveSuccessHint'
  | 'env.openedInEditor'
  | 'sidebar.versionHintStore'
  | 'sidebar.versionHintTray'
  | 'sidebar.updateAvailable'
  | 'sidebar.updateAvailableHint'
  | 'sidebar.theme'
  | 'sidebar.themeAuto'
  | 'sidebar.themeLight'
  | 'sidebar.themeDark'

export type Messages = Record<MessageKey, string>

export const messages: Record<SupportedLocale, Messages> = {
  en: {
    'common.gatewayLabel': 'CCX CORE',
    'common.online': 'GATEWAY ONLINE',
    'common.connecting': 'CONNECTING...',
    'common.offline': 'GATEWAY OFFLINE',
    'common.refreshWebUI': 'Refresh Web UI',
    'common.version': 'Version',
    'common.gatewayPort': 'Gateway port',
    'common.daemonPid': 'Daemon PID',
    'common.autoStart': 'Autostart',
    'common.autoStartOn': 'Enabled',
    'common.autoStartOff': 'Disabled',
    'common.serviceHealthy': 'Service healthy',
    'common.serviceStarting': 'Gateway starting',
    'common.serviceDisconnected': 'Service disconnected',
    'common.settings': 'Settings',
    'common.save': 'Save',
    'common.cancel': 'Cancel',
    'common.retry': 'Retry',
    'nav.status': 'Status',
    'nav.statusDesc': 'Live status and logs',
    'nav.agent': 'Agent',
    'nav.agentDesc': 'Local agent configuration',
    'nav.channels': 'Channels',
    'nav.channelsDesc': 'Add upstream channels',
    'nav.env': 'Environment',
    'nav.envDesc': 'Edit gateway env file',
    'nav.web': 'Console',
    'nav.webDesc': 'CCX Web control panel',
    'tab.statusTitle': 'Gateway Status',
    'tab.agentTitle': 'Agent Config',
    'tab.channelsTitle': 'Channel Center',
    'tab.envTitle': 'Environment Settings',
    'tab.webTitle': 'Built-in Web UI',
    'sidebar.language': 'Language',
    'sidebar.languageEnglish': 'English',
    'sidebar.languageChinese': 'Chinese',
    'setup.loading': 'Initializing CCX Console',
    'setup.title': 'CCX Desktop Initial Setup',
    'setup.description': 'PROXY_ACCESS_KEY is the credential AI agents use to access upstream APIs through the CCX proxy. Every caller must have this key.',
    'setup.regenerate': 'Regenerate',
    'setup.regenerateTitle': 'Generate a new random key',
    'setup.hide': 'Hide',
    'setup.show': 'Show',
    'setup.copied': 'Copied',
    'setup.copyKey': 'Copy key',
    'setup.configPath': 'Config file path',
    'setup.copyPath': 'Copy path',
    'setup.hint': 'After saving, CCX will start automatically. You can adjust other settings later on the Environment page.',
    'setup.saving': 'Saving and starting...',
    'setup.submit': 'Finish setup and start',
    'webui.notRunning': 'CCX service is not running, so the Web UI cannot be displayed.',
    'webui.openInBrowser': 'Open in browser',
    'metrics.gatewayPort': 'Gateway port',
    'metrics.uptime': 'Uptime',
    'metrics.channels': 'Channels',
    'metrics.version': 'Gateway version',
    'actions.start': 'Start',
    'actions.stop': 'Stop',
    'actions.restart': 'Restart',
    'actions.openWebUI': 'Open Web UI',
    'actions.openBrowser': 'Open in browser',
    'actions.refreshStatus': 'Refresh status',
    'details.title': 'Service details',
    'details.binary': 'Binary',
    'details.binaryMissing': 'Not found',
    'details.dataDir': 'Data dir',
    'details.dataDirMissing': 'Not configured',
    'details.healthStatus': 'Health',
    'details.revealDir': 'Reveal directory',
    'details.openDir': 'Open directory',
    'channel.headerEyebrow': 'Channel Preset Center',
    'channel.title': 'Channel Center',
    'channel.description': 'Use DeepSeek, MiMo, Kimi, GLM, and MiniMax keys for both direct Agent routing and the unified CCX channel pool. Provider presets handle advanced switches automatically.',
    'channel.hasKey': 'Key saved',
    'channel.promo': 'Register via promotion link to claim a ¥5 trial credit',
    'channel.console': 'Open official console',
    'channel.target': 'Target',
    'channel.keySavedPlaceholder': 'Saved locally; leave empty to reuse this Provider Key',
    'channel.keyInputPlaceholder': 'Enter API Key. It is stored only in local Desktop config',
    'channel.missingKey': 'Enter an API Key, or save this Provider key in Agent Config first.',
    'channel.reuseKey': 'Will reuse the saved local {provider} key',
    'channel.name': 'Channel name',
    'channel.nameHint': 'A channel with the same name will be overwritten. Use another name to create a separate channel.',
    'channel.presetWrites': 'Preset will write automatically',
    'channel.capabilityHint': 'Capability toggles: reasoning / vision / model list / compatibility fields are configured by Provider.',
    'channel.addToCCX': 'Add to CCX',
    'channel.badgeDirectAgent': 'Agent direct',
    'channel.badgeNativeMessages': 'Native Messages',
    'env.title': 'Environment config',
    'env.pathDetecting': 'Detecting',
    'env.refresh': 'Refresh',
    'env.openingEditor': 'Opening…',
    'env.openWithEditor': 'Open in {editor}',
    'env.openInEditor': 'Open in editor',
    'env.selectEditor': 'Choose editor…',
    'env.save': 'Save',
    'env.saving': 'Saving',
    'env.saved': 'Saved',
    'env.failed': 'Failed',
    'env.hide': 'Hide',
    'env.show': 'Show',
    'env.copied': 'Copied',
    'env.copy': 'Copy',
    'env.fieldRequired': '{field} is required',
    'env.fieldDisallow': '{field} cannot use a sample placeholder value',
    'env.fieldNumber': '{field} must be a number',
    'env.fieldInteger': '{field} must be an integer',
    'env.fieldMin': '{field} must be at least {min}',
    'env.fieldMax': '{field} must be at most {max}',
    'env.loadFailed': 'Failed to load config: {error}',
    'env.saveSuccess': '.env saved; changes take effect after restarting the service',
    'agent.statusDetecting': 'Detecting',
    'agent.statusConfigured': 'Configured',
    'agent.statusPortMismatch': 'Port mismatch',
    'agent.statusUnconfigured': 'Not configured',
    'agent.localGateway': 'Current CCX gateway',
    'agent.custom': 'Custom',
    'agent.currentProvider': 'Current provider',
    'agent.currentUrl': 'Current URL',
    'agent.targetUrl': 'Target URL',
    'agent.notSet': 'Not configured',
    'agent.configPath': 'Config file',
    'agent.authPath': 'Auth file',
    'agent.openFileInEditor': 'Open in editor',
    'agent.applyConfig': 'Apply config',
    'agent.restoreConfig': 'Restore original config',
    'agent.openConsole': 'Open official console',
    'agent.codexPlaceholderSaved': 'Saved locally; leave empty to reuse the key',
    'agent.codexPlaceholderRequired': 'Required: enter API Key',
    'agent.codexPlaceholderWriteOnly': 'Written to Codex config only',
    'agent.diffPreviewApply': 'Apply config preview',
    'agent.diffPreviewRestore': 'Restore config preview',
    'agent.diffConfirmApply': 'Confirm apply',
    'agent.diffConfirmRestore': 'Confirm restore',
    'agent.diffActionCreate': 'Create',
    'agent.diffActionDelete': 'Delete',
    'agent.diffActionModify': 'Modify',
    'agent.diffComputing': 'Computing changes…',
    'agent.diffNoChanges': 'No changes',
    'agent.diffExpandContext': 'Expand {count} unchanged lines',
    'agent.diffCollapseContext': 'Collapse {count} lines',
    'agent.diffCancel': 'Cancel',
    'agent.provider.localGateway': 'CCX local gateway',
    'agent.provider.deepseekDirect': 'DeepSeek direct',
    'agent.provider.mimoDirect': 'MiMo direct',
    'agent.provider.compshareDirect': 'Compshare direct',
    'agent.provider.kimiDirect': 'Kimi direct',
    'agent.provider.glmDirect': 'GLM direct',
    'agent.provider.minimaxDirect': 'MiniMax direct',
    'agent.provider.dashscopeDirect': 'DashScope direct',
    'agent.provider.opencodeZenDirect': 'OpenCode Zen direct',
    'agent.provider.opencodeGoDirect': 'OpenCode Go direct',
    'agent.provider.openaiDirect': 'OpenAI direct',
    'agent.promo': 'Register via promotion link to claim a ¥5 trial credit',
    'agent.planPayAsYouGo': 'Pay-as-you-go',
    'agent.planChina': 'Subscription - China',
    'agent.planSingapore': 'Subscription - Singapore',
    'agent.planEurope': 'Subscription - Europe',
    'agent.planSubscription': 'Subscription',
    'agent.billingModeMiMo': 'MiMo billing mode',
    'agent.billingModeDashScope': 'DashScope billing mode',
    'agent.placeholderSaved': 'Saved locally; leave empty to reuse the key',
    'agent.placeholderMimo': 'Required: MiMo API Key (tp-xxx or account key)',
    'agent.placeholderDashScope': 'Required: DashScope API Key (sk-xxx or sk-sp-xxx)',
    'agent.placeholderRequired': 'Required: enter API Key',
    'env.groupAccess': 'Access control',
    'env.groupAccessDesc': 'Access keys for the proxy and admin endpoints.',
    'env.fieldProxyAccessKey': 'Proxy access key',
    'env.placeholderProxyAccessKey': 'Enter a strong random key',
    'env.fieldAdminAccessKey': 'Admin API key',
    'env.placeholderAdminAccessKey': 'Leave empty to fall back to PROXY_ACCESS_KEY',
    'env.descAdminAccessKey': 'Used for the Web UI and /api/* endpoints.',
    'env.groupServer': 'Server config',
    'env.groupServerDesc': 'Desktop injects some runtime values on start; this section still matches .env.example.',
    'env.fieldPort': 'Port',
    'env.descPort': 'Preferred startup port. If occupied, the next available port is used automatically.',
    'env.fieldEnv': 'Environment',
    'env.descEnv': 'production is recommended.',
    'env.groupWebUI': 'Web UI config',
    'env.groupWebUIDesc': 'Control whether the Web UI is enabled and which language it uses by default.',
    'env.fieldEnableWebUI': 'Enable Web UI',
    'env.descEnableWebUI': 'Desktop mode usually forces this on.',
    'env.fieldAppUILanguage': 'Default language',
    'env.groupLogs': 'Logging config',
    'env.groupLogsDesc': 'Control request/response logs, SSE debugging, and response model rewriting.',
    'env.fieldLogLevel': 'Log level',
    'env.fieldEnableRequestLogs': 'Enable request logs',
    'env.fieldEnableResponseLogs': 'Enable response logs',
    'env.descEnableResponseLogs': 'Response logs may expose sensitive content.',
    'env.fieldQuietPollingLogs': 'Quiet polling logs',
    'env.fieldRawLogOutput': 'Raw log output',
    'env.fieldSseDebugLevel': 'SSE debug level',
    'env.fieldRewriteResponseModel': 'Rewrite response model',
    'env.groupPerformance': 'Performance config',
    'env.groupPerformanceDesc': 'Request-chain timeouts and request body size limits.',
    'env.fieldRequestTimeout': 'Request timeout (ms)',
    'env.fieldServerReadTimeout': 'Server read timeout (ms)',
    'env.fieldMaxRequestBodySize': 'Max request body size (MB)',
    'env.fieldResponseHeaderTimeout': 'Response header timeout (s)',
    'env.groupCors': 'CORS config',
    'env.groupCorsDesc': 'Cross-origin access control.',
    'env.fieldEnableCors': 'Enable CORS',
    'env.fieldCorsOrigin': 'Allowed Origin',
    'env.groupCircuitBreaker': 'Circuit-breaker metrics config',
    'env.groupCircuitBreakerDesc': 'Control scheduler metric windows and failure-rate thresholds.',
    'env.fieldMetricsWindowSize': 'Window size',
    'env.fieldMetricsFailureThreshold': 'Failure-rate threshold',
    'env.groupMetricsPersistence': 'Metrics persistence config',
    'env.groupMetricsPersistenceDesc': 'Control SQLite metrics persistence and data retention.',
    'env.fieldMetricsPersistenceEnabled': 'Enable metrics persistence',
    'env.fieldMetricsRetentionDays': 'Metrics retention days',
    'logs.searchPlaceholder': 'Search logs...',
    'logs.autoScroll': 'Auto-scroll to bottom',
    'logs.copied': 'Copied!',
    'logs.copyAll': 'Copy all logs',
    'logs.clear': 'Clear log console',
    'logs.noSearchResults': 'No matching log lines found',
    'logs.empty': 'No logs yet. Start the service to view output.',
    'diagnostic.binaryTitle': 'Binary not found',
    'diagnostic.binarySuggestionBuild': 'Confirm the CCX binary has been built: cd backend-go && make build',
    'diagnostic.binarySuggestionCheckDataDir': 'Check whether ccx-go / ccx-go.exe exists in the Desktop data directory',
    'diagnostic.binarySuggestionDownload': 'Build the backend first, or download a prebuilt release from the Release page',
    'diagnostic.portTitle': 'Port conflict',
    'diagnostic.portSuggestionInstance': 'Check whether another CCX instance is already running',
    'diagnostic.portSuggestionEnv': 'Change the PORT field in .env to another port',
    'diagnostic.portSuggestionInspect': 'Use lsof -i :3688 (macOS/Linux) or netstat -ano | findstr :3688 (Windows) to inspect port usage',
    'diagnostic.healthTitle': 'Health check timeout',
    'diagnostic.healthSuggestionLogs': 'Check the log panel for startup errors',
    'diagnostic.healthSuggestionEnv': 'Check whether .env has syntax errors',
    'diagnostic.healthSuggestionChannels': 'Confirm upstream channel configuration is correct. First startup may take longer.',
    'diagnostic.healthSuggestionRestart': 'Try restarting the service manually',
    'diagnostic.permissionTitle': 'Insufficient permissions',
    'diagnostic.permissionSuggestionDataDir': 'Check whether the data directory is writable',
    'diagnostic.permissionSuggestionExecutable': 'macOS/Linux: confirm the binary is executable (chmod +x)',
    'diagnostic.permissionSuggestionWindows': 'Windows: try running as administrator',
    'diagnostic.genericTitle': 'Startup failed',
    'diagnostic.genericSuggestionLogs': 'Check the log panel below for detailed errors',
    'diagnostic.genericSuggestionRestart': 'Try restarting the service',
    'setup.errorEmptyKey': 'PROXY_ACCESS_KEY cannot be empty',
    'env.saveSuccessHint': '.env saved; changes take effect after restarting the service',
    'env.openedInEditor': '.env opened in editor',
    'sidebar.versionHintStore': 'Microsoft Store version updates automatically via Store',
    'sidebar.versionHintTray': 'Check for updates from the tray menu',
    'sidebar.updateAvailable': 'New {version}',
    'sidebar.updateAvailableHint': 'Click to view release notes',
    'sidebar.theme': 'Theme',
    'sidebar.themeAuto': 'System',
    'sidebar.themeLight': 'Light',
    'sidebar.themeDark': 'Dark',
  },
  'zh-CN': {
    'common.gatewayLabel': 'CCX CORE',
    'common.online': 'GATEWAY ONLINE',
    'common.connecting': 'CONNECTING...',
    'common.offline': 'GATEWAY OFFLINE',
    'common.refreshWebUI': '刷新 Web UI',
    'common.version': '当前版本',
    'common.gatewayPort': '网关端口',
    'common.daemonPid': '守护 PID',
    'common.autoStart': '开机自启',
    'common.autoStartOn': '已开启',
    'common.autoStartOff': '已关闭',
    'common.serviceHealthy': '运行正常',
    'common.serviceStarting': '网关启动中',
    'common.serviceDisconnected': '服务已断开',
    'common.settings': '设置',
    'common.save': '保存',
    'common.cancel': '取消',
    'common.retry': '重试',
    'nav.status': '网关监控',
    'nav.statusDesc': '实时状态及核心日志',
    'nav.agent': 'Agent 配置',
    'nav.agentDesc': '本地开发代理控制',
    'nav.channels': '渠道中心',
    'nav.channelsDesc': '一键添加上游渠道',
    'nav.env': '环境参数',
    'nav.envDesc': '网关配置文件编辑',
    'nav.web': '管理界面',
    'nav.webDesc': 'CCX Web 控制面板',
    'tab.statusTitle': '网关状态监控',
    'tab.agentTitle': 'Agent 代理配置',
    'tab.channelsTitle': '渠道中心',
    'tab.envTitle': '环境参数管理',
    'tab.webTitle': '内置控制台 Web UI',
    'sidebar.language': '语言',
    'sidebar.languageEnglish': 'English',
    'sidebar.languageChinese': '中文',
    'setup.loading': '正在初始化 CCX 控制台',
    'setup.title': 'CCX Desktop 初始配置',
    'setup.description': 'PROXY_ACCESS_KEY 是 AI Agent 通过 CCX 代理访问上游 API 的身份凭证，所有调用方必须持有此密钥。',
    'setup.regenerate': '重新生成',
    'setup.regenerateTitle': '重新生成随机密钥',
    'setup.hide': '隐藏',
    'setup.show': '显示',
    'setup.copied': '已复制',
    'setup.copyKey': '复制密钥',
    'setup.configPath': '配置文件路径',
    'setup.copyPath': '复制路径',
    'setup.hint': '保存后 CCX 将自动启动。后续可在主界面【环境参数】页继续调整其他配置。',
    'setup.saving': '正在保存并启动...',
    'setup.submit': '完成配置并启动',
    'webui.notRunning': 'CCX 服务尚未启动，无法显示 Web UI。',
    'webui.openInBrowser': '浏览器打开',
    'metrics.gatewayPort': '网关端口',
    'metrics.uptime': '运行时长',
    'metrics.channels': '调度信道',
    'metrics.version': '网关版本',
    'actions.start': '启动服务',
    'actions.stop': '停止服务',
    'actions.restart': '重启服务',
    'actions.openWebUI': '进入 Web UI',
    'actions.openBrowser': '浏览器直达',
    'actions.refreshStatus': '刷新当前状态',
    'details.title': '服务详情',
    'details.binary': '二进制',
    'details.binaryMissing': '未发现',
    'details.dataDir': '数据目录',
    'details.dataDirMissing': '未设置',
    'details.healthStatus': '健康状态',
    'details.revealDir': '打开所在目录',
    'details.openDir': '打开目录',
    'channel.headerEyebrow': 'Channel Preset Center',
    'channel.title': '渠道中心',
    'channel.description': '统一把 DeepSeek、MiMo、Kimi、GLM、MiniMax Key 可同时用于 Agent 直连和 CCX 统一渠道池，复杂开关由预设自动处理。',
    'channel.hasKey': '已有 Key',
    'channel.promo': '通过推广链接注册，领取 5 元平台试用金',
    'channel.console': '访问官方控制台',
    'channel.target': '添加目标',
    'channel.keySavedPlaceholder': '已保存，留空则复用该 Provider Key',
    'channel.keyInputPlaceholder': '输入 API Key，仅保存在本机 Desktop 配置中',
    'channel.missingKey': '请填写 API Key，或先在 Agent 配置中保存该 Provider 的 key。',
    'channel.reuseKey': '将复用本机已保存的 {provider} Key',
    'channel.name': '渠道名称',
    'channel.nameHint': '同名渠道会被直接覆盖更新；如需新建独立渠道，请改用不同名称。',
    'channel.presetWrites': '预设将自动写入',
    'channel.capabilityHint': '能力开关：reasoning / vision / model list / 兼容字段会按 Provider 自动配置。',
    'channel.addToCCX': '添加到 CCX',
    'channel.badgeDirectAgent': 'Agent 直连',
    'channel.badgeNativeMessages': 'Messages 原生',
    'env.title': '环境配置',
    'env.pathDetecting': '检测中',
    'env.refresh': '刷新',
    'env.openingEditor': '打开中…',
    'env.openWithEditor': '用 {editor} 打开',
    'env.openInEditor': '用编辑器打开',
    'env.selectEditor': '选择编辑器…',
    'env.save': '保存',
    'env.saving': '保存中',
    'env.saved': '已保存',
    'env.failed': '失败',
    'env.hide': '隐藏',
    'env.show': '显示',
    'env.copied': '已复制',
    'env.copy': '复制',
    'env.fieldRequired': '{field}不能为空',
    'env.fieldDisallow': '{field}不能使用示例占位值',
    'env.fieldNumber': '{field}必须是数字',
    'env.fieldInteger': '{field}必须是整数',
    'env.fieldMin': '{field}不能小于 {min}',
    'env.fieldMax': '{field}不能大于 {max}',
    'env.loadFailed': '加载配置失败：{error}',
    'env.saveSuccess': '.env 已保存，重启服务后生效',
    'agent.statusDetecting': '检测中',
    'agent.statusConfigured': '已配置',
    'agent.statusPortMismatch': '端口不匹配',
    'agent.statusUnconfigured': '未配置',
    'agent.localGateway': '当前 CCX 网关',
    'agent.custom': '自定义',
    'agent.currentProvider': '当前 Provider',
    'agent.currentUrl': '当前 URL',
    'agent.targetUrl': '目标 URL',
    'agent.notSet': '未设置',
    'agent.configPath': '配置文件',
    'agent.authPath': '认证文件',
    'agent.openFileInEditor': '用编辑器打开',
    'agent.applyConfig': '应用配置',
    'agent.restoreConfig': '恢复原始配置',
    'agent.openConsole': '访问官方控制台',
    'agent.codexPlaceholderSaved': '已保存，留空则使用已保存的 key',
    'agent.codexPlaceholderRequired': '必填：输入 API Key',
    'agent.codexPlaceholderWriteOnly': '仅写入 Codex 配置',
    'agent.diffPreviewApply': '应用配置预览',
    'agent.diffPreviewRestore': '恢复配置预览',
    'agent.diffConfirmApply': '确认应用',
    'agent.diffConfirmRestore': '确认恢复',
    'agent.diffActionCreate': '创建',
    'agent.diffActionDelete': '删除',
    'agent.diffActionModify': '修改',
    'agent.diffComputing': '计算变更中...',
    'agent.diffNoChanges': '无变更',
    'agent.diffExpandContext': '展开 {count} 行未变更内容',
    'agent.diffCollapseContext': '收起 {count} 行',
    'agent.diffCancel': '取消',
    'agent.provider.localGateway': 'CCX 本地网关',
    'agent.provider.deepseekDirect': 'DeepSeek 直连',
    'agent.provider.mimoDirect': 'MiMo 直连',
    'agent.provider.compshareDirect': 'Compshare 直连',
    'agent.provider.kimiDirect': 'Kimi 直连',
    'agent.provider.glmDirect': 'GLM 直连',
    'agent.provider.minimaxDirect': 'MiniMax 直连',
    'agent.provider.dashscopeDirect': 'DashScope 直连',
    'agent.provider.opencodeZenDirect': 'OpenCode Zen 直连',
    'agent.provider.opencodeGoDirect': 'OpenCode Go 直连',
    'agent.provider.openaiDirect': 'OpenAI 直连',
    'agent.promo': '通过推广链接注册，领取 5 元平台试用金',
    'agent.planPayAsYouGo': '按量',
    'agent.planChina': '订阅套餐 - 中国',
    'agent.planSingapore': '订阅套餐 - 新加坡',
    'agent.planEurope': '订阅套餐 - 欧洲',
    'agent.planSubscription': '订阅套餐',
    'agent.billingModeMiMo': 'MiMo 计费模式',
    'agent.billingModeDashScope': 'DashScope 计费模式',
    'agent.placeholderSaved': '已保存，留空则使用已保存的 key',
    'agent.placeholderMimo': '必填：MiMo API Key（tp-xxx 或账号 key）',
    'agent.placeholderDashScope': '必填：DashScope API Key（sk-xxx 或 sk-sp-xxx）',
    'agent.placeholderRequired': '必填：输入 API Key',
    'env.groupAccess': '访问控制',
    'env.groupAccessDesc': '代理入口与管理入口的访问密钥。',
    'env.fieldProxyAccessKey': '代理访问密钥',
    'env.placeholderProxyAccessKey': '请输入强随机密钥',
    'env.fieldAdminAccessKey': '管理 API 独立密钥',
    'env.placeholderAdminAccessKey': '留空则回退到 PROXY_ACCESS_KEY',
    'env.descAdminAccessKey': '用于管理界面和 /api/* 端点。',
    'env.groupServer': '服务器配置',
    'env.groupServerDesc': 'Desktop 会在启动时注入部分运行参数；这里仍完整覆盖 .env.example。',
    'env.fieldPort': '服务端口',
    'env.descPort': '启动时优先使用此端口，被占用时自动递增分配。',
    'env.fieldEnv': '运行环境',
    'env.descEnv': 'production 为推荐值。',
    'env.groupWebUI': 'Web UI 配置',
    'env.groupWebUIDesc': '控制管理界面是否启用以及默认语言。',
    'env.fieldEnableWebUI': '启用 Web UI',
    'env.descEnableWebUI': 'Desktop 模式通常会强制启用。',
    'env.fieldAppUILanguage': '默认语言',
    'env.groupLogs': '日志配置',
    'env.groupLogsDesc': '控制请求/响应日志、SSE 调试和模型字段改写。',
    'env.fieldLogLevel': '日志级别',
    'env.fieldEnableRequestLogs': '启用请求日志',
    'env.fieldEnableResponseLogs': '启用响应日志',
    'env.descEnableResponseLogs': '响应日志可能增加敏感内容暴露风险。',
    'env.fieldQuietPollingLogs': '静默轮询日志',
    'env.fieldRawLogOutput': '原始日志输出',
    'env.fieldSseDebugLevel': 'SSE 调试级别',
    'env.fieldRewriteResponseModel': '改写响应 model',
    'env.groupPerformance': '性能配置',
    'env.groupPerformanceDesc': '请求链路超时和请求体大小限制。',
    'env.fieldRequestTimeout': '请求超时（毫秒）',
    'env.fieldServerReadTimeout': '服务端读取超时（毫秒）',
    'env.fieldMaxRequestBodySize': '请求体最大大小（MB）',
    'env.fieldResponseHeaderTimeout': '响应头超时（秒）',
    'env.groupCors': 'CORS 配置',
    'env.groupCorsDesc': '跨域访问控制。',
    'env.fieldEnableCors': '启用 CORS',
    'env.fieldCorsOrigin': '允许的 Origin',
    'env.groupCircuitBreaker': '熔断指标配置',
    'env.groupCircuitBreakerDesc': '控制调度指标窗口与失败率阈值。',
    'env.fieldMetricsWindowSize': '滑动窗口大小',
    'env.fieldMetricsFailureThreshold': '失败率阈值',
    'env.groupMetricsPersistence': '指标持久化配置',
    'env.groupMetricsPersistenceDesc': '控制 SQLite 指标持久化与数据保留。',
    'env.fieldMetricsPersistenceEnabled': '启用指标持久化',
    'env.fieldMetricsRetentionDays': '指标保留天数',
    'logs.searchPlaceholder': '搜索日志...',
    'logs.autoScroll': '自动滚动到底部',
    'logs.copied': '已复制！',
    'logs.copyAll': '复制全部日志',
    'logs.clear': '清空日志控制台',
    'logs.noSearchResults': '未找到匹配的日志行',
    'logs.empty': '暂无日志输出，启动服务后即可查看',
    'diagnostic.binaryTitle': '二进制文件未找到',
    'diagnostic.binarySuggestionBuild': '确认 CCX 二进制已构建: cd backend-go && make build',
    'diagnostic.binarySuggestionCheckDataDir': '检查 Desktop 数据目录中是否存在 ccx-go / ccx-go.exe',
    'diagnostic.binarySuggestionDownload': '首次使用需先构建后端，或从 Release 页面下载预编译版本',
    'diagnostic.portTitle': '端口冲突',
    'diagnostic.portSuggestionInstance': '检查是否有其他 CCX 实例已在运行',
    'diagnostic.portSuggestionEnv': '修改 .env 中 PORT 字段使用其他端口',
    'diagnostic.portSuggestionInspect': '使用 lsof -i :3688 (macOS/Linux) 或 netstat -ano | findstr :3688 (Windows) 检查端口占用',
    'diagnostic.healthTitle': '健康检查超时',
    'diagnostic.healthSuggestionLogs': '查看日志面板中是否有启动错误信息',
    'diagnostic.healthSuggestionEnv': '检查 .env 配置是否有语法错误',
    'diagnostic.healthSuggestionChannels': '确认上游渠道配置正确，首次启动可能需要较长时间',
    'diagnostic.healthSuggestionRestart': '尝试手动重启服务',
    'diagnostic.permissionTitle': '权限不足',
    'diagnostic.permissionSuggestionDataDir': '检查数据目录是否有写入权限',
    'diagnostic.permissionSuggestionExecutable': 'macOS/Linux: 确认二进制文件有执行权限 (chmod +x)',
    'diagnostic.permissionSuggestionWindows': 'Windows: 尝试以管理员身份运行',
    'diagnostic.genericTitle': '启动失败',
    'diagnostic.genericSuggestionLogs': '查看下方日志面板获取详细错误信息',
    'diagnostic.genericSuggestionRestart': '尝试重启服务',
    'setup.errorEmptyKey': 'PROXY_ACCESS_KEY 不能为空',
    'env.saveSuccessHint': '.env 已保存，重启服务后生效',
    'env.openedInEditor': '已在编辑器中打开 .env 文件',
    'sidebar.versionHintStore': 'Microsoft Store 版本由 Store 自动更新',
    'sidebar.versionHintTray': '通过托盘菜单检查更新',
    'sidebar.updateAvailable': '新版 {version}',
    'sidebar.updateAvailableHint': '点击查看发布说明',
    'sidebar.theme': '主题',
    'sidebar.themeAuto': '跟随系统',
    'sidebar.themeLight': '亮色',
    'sidebar.themeDark': '暗色',
  },
}
