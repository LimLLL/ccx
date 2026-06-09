import { Browser } from '@wailsio/runtime'

export const providerConsoleLinks: Record<string, string> = {
  deepseek: 'https://platform.deepseek.com/usage',
  mimo: 'https://platform.xiaomimimo.com/console/balance',
  compshare: 'https://console.compshare.cn/light-gpu/model-manage',
  runapi: 'https://runapi.co/console',
  'tencent-lkeap': 'https://console.cloud.tencent.com/lkeap/token-plan',
  'kimi-code': 'https://www.kimi.com/code/console',
  'volc-ark': 'https://console.volcengine.com/ark',
  qianfan: 'https://console.bce.baidu.com/qianfan/resource/subscribe',
  originrouter: 'https://easytransnote.com/ai/console/#key',
  kimi: 'https://platform.moonshot.cn/console/account',
  glm: 'https://open.bigmodel.cn/coding-plan/personal/overview',
  minimax: 'https://platform.minimaxi.com/user-center/payment/balance',
  dashscope: 'https://bailian.console.aliyun.com/cn-beijing?tab=model#/api-key',
  'opencode-zen': 'https://opencode.ai/',
  'opencode-go': 'https://opencode.ai/',
  openrouter: 'https://openrouter.ai/keys',
  modelscope: 'https://modelscope.cn/my/myaccesstoken',
  openai: 'https://platform.openai.com',
}

export const providerPromotionLinks: Record<string, string> = {
  compshare: 'https://www.compshare.cn/?ytag=GPU_YY_git_ccx',
  runapi: 'https://runapi.co/register?aff=CqQO',
}

export function openExternalLink(url: string) {
  return Browser.OpenURL(url)
}

export function openProviderConsole(provider: string) {
  const url = providerConsoleLinks[provider]
  if (!url) return Promise.resolve()
  return openExternalLink(url)
}

export function openProviderPromotion(provider: string) {
  const url = providerPromotionLinks[provider]
  if (!url) return Promise.resolve()
  return openExternalLink(url)
}
