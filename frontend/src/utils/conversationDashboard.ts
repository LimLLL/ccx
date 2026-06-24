import type { ConversationInfo } from '@/services/api'

export type BoardColumnKey = 'streaming' | 'subagents' | 'active' | 'idle'

export function buildConversationColumnBuckets(items: ConversationInfo[]): Record<BoardColumnKey, ConversationInfo[]> {
  const buckets: Record<BoardColumnKey, ConversationInfo[]> = {
    streaming: [],
    subagents: [],
    active: [],
    idle: [],
  }

  for (const item of items) {
    buckets[getConversationBoardColumnKey(item)].push(item)
  }

  return buckets
}

export function getConversationBoardColumnKey(conversation: ConversationInfo): BoardColumnKey {
  if (conversation.hasSubagents) return 'subagents'
  if (conversation.status === 'streaming') return 'streaming'
  if (conversation.status === 'idle') return 'idle'
  return 'active'
}
