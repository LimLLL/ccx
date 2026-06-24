import { describe, expect, it } from 'vitest'
import type { ConversationInfo } from '@/services/api'
import { buildConversationColumnBuckets, getConversationBoardColumnKey } from './conversationDashboard'

function conversation(overrides: Partial<ConversationInfo>): ConversationInfo {
  return {
    id: overrides.id ?? 'conv-1',
    kind: overrides.kind ?? 'messages',
    userId: overrides.userId ?? 'user',
    createdAt: '2026-01-01T00:00:00Z',
    lastActiveAt: '2026-01-01T00:00:00Z',
    requestCount: 1,
    models: ['model'],
    currentChannel: 0,
    channelName: 'primary',
    status: overrides.status ?? 'active',
    lastModel: 'model',
    lastRequestId: '',
    ...overrides,
  }
}

describe('conversation dashboard columns', () => {
  it('prioritizes subagent conversations over streaming status', () => {
    const item = conversation({ status: 'streaming', hasSubagents: true })

    expect(getConversationBoardColumnKey(item)).toBe('subagents')
  })

  it('buckets conversations into mutually exclusive columns', () => {
    const buckets = buildConversationColumnBuckets([
      conversation({ id: 'streaming', status: 'streaming' }),
      conversation({ id: 'subagents', hasSubagents: true }),
      conversation({ id: 'idle', status: 'idle' }),
      conversation({ id: 'active', status: 'active' }),
    ])

    expect(buckets.streaming.map(c => c.id)).toEqual(['streaming'])
    expect(buckets.subagents.map(c => c.id)).toEqual(['subagents'])
    expect(buckets.idle.map(c => c.id)).toEqual(['idle'])
    expect(buckets.active.map(c => c.id)).toEqual(['active'])
  })
})
