'use client';

import { useState, useEffect } from 'react';
import { InboxMessage, Thread, messageAPI, Message } from '@/lib/api';

interface ChatPaneProps {
  selectedThreadId: string | null;
  selectedInboxMessage?: InboxMessage | null;
  threads?: Thread[];
  onInboxMessageAssigned?: () => void;
}

export default function ChatPane({ selectedThreadId, selectedInboxMessage, threads = [], onInboxMessageAssigned }: ChatPaneProps) {
  const [showAssignModal, setShowAssignModal] = useState(false);
  const [assigningToThread, setAssigningToThread] = useState(false);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loadingMessages, setLoadingMessages] = useState(false);

  useEffect(() => {
    if (selectedThreadId) {
      loadThreadMessages(selectedThreadId);
    }
  }, [selectedThreadId]);

  const loadThreadMessages = async (threadId: string) => {
    setLoadingMessages(true);
    try {
      const response = await messageAPI.getThreadMessages(threadId);
      setMessages(response.messages);
    } catch (error) {
      console.error('Failed to load thread messages:', error);
    } finally {
      setLoadingMessages(false);
    }
  };

  const handleAssignToThread = async (threadId: string) => {
    if (!selectedInboxMessage) return;

    setAssigningToThread(true);
    try {
      await messageAPI.assignInboxMessageToThread(selectedInboxMessage.id, threadId);
      setShowAssignModal(false);

      // Trigger inbox refresh
      window.dispatchEvent(new Event('refreshInboxMessages'));

      if (onInboxMessageAssigned) {
        onInboxMessageAssigned();
      }
    } catch (error) {
      console.error('Failed to assign message to thread:', error);
      alert('Failed to assign message to thread');
    } finally {
      setAssigningToThread(false);
    }
  };

  // Show inbox message if selected
  if (selectedInboxMessage) {
    return (
      <div className="flex-1 flex flex-col bg-background">
        {/* Email Header */}
        <div className="border-b border-border px-6 py-4 bg-card">
          <h2 className="text-xl font-semibold text-card-foreground mb-2">
            {selectedInboxMessage.subject || 'No Subject'}
          </h2>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <span className="font-medium">From:</span>
            <span>{selectedInboxMessage.senderEmail}</span>
          </div>
          <div className="text-xs text-muted-foreground mt-1">
            {new Date(selectedInboxMessage.timestamp).toLocaleString()}
          </div>
        </div>

        {/* Email Content */}
        <div className="flex-1 overflow-y-auto px-6 py-6 bg-muted/50">
          <div className="bg-card rounded-lg border border-border p-6">
            <div className="whitespace-pre-wrap text-card-foreground">
              {selectedInboxMessage.content}
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="border-t border-border px-6 py-4 bg-card">
          <div className="flex gap-3">
            <button
              onClick={() => setShowAssignModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors cursor-pointer"
            >
              Assign to Thread
            </button>
            <button className="bg-secondary hover:bg-secondary/80 text-secondary-foreground px-6 py-2 rounded-lg font-medium transition-colors cursor-pointer">
              Archive
            </button>
          </div>
        </div>

        {/* Assign to Thread Modal */}
        {showAssignModal && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={() => !assigningToThread && setShowAssignModal(false)}>
            <div className="bg-popover rounded-lg p-6 max-w-md w-full mx-4 border border-border" onClick={(e) => e.stopPropagation()}>
              <h2 className="text-xl font-bold text-popover-foreground mb-4">Assign to Thread</h2>
              <p className="text-sm text-muted-foreground mb-4">
                Choose which seller thread to assign this email to:
              </p>

              {threads.length === 0 ? (
                <div className="text-muted-foreground text-sm italic py-4">
                  No threads available. Create a new seller thread first.
                </div>
              ) : (
                <div className="space-y-2 mb-4 max-h-96 overflow-y-auto">
                  {threads.map((thread) => (
                    <button
                      key={thread.id}
                      onClick={() => handleAssignToThread(thread.id)}
                      disabled={assigningToThread}
                      className="w-full text-left px-4 py-3 rounded-md border border-border hover:bg-accent hover:text-accent-foreground transition-colors disabled:opacity-50 cursor-pointer"
                    >
                      <div className="font-medium">{thread.sellerName}</div>
                      <div className="text-xs text-muted-foreground mt-0.5 capitalize">{thread.sellerType}</div>
                    </button>
                  ))}
                </div>
              )}

              <button
                onClick={() => setShowAssignModal(false)}
                disabled={assigningToThread}
                className="w-full bg-secondary hover:bg-secondary/80 text-secondary-foreground font-medium py-2 px-4 rounded-md transition-colors disabled:opacity-50 cursor-pointer"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>
    );
  }

  if (!selectedThreadId) {
    return (
      <div className="flex-1 flex items-center justify-center bg-background">
        <div className="text-center max-w-md">
          <div className="mb-6 flex justify-center">
            <svg
              className="w-24 h-24 text-muted"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              />
              <circle cx="19" cy="7" r="3" fill="currentColor" stroke="none" />
              <text x="19" y="8.5" fontSize="2.5" textAnchor="middle" fill="white" fontWeight="bold">+</text>
            </svg>
          </div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">
            No Active Negotiations
          </h2>
          <p className="text-muted-foreground">
            Add a seller in the sidebar to begin chatting.
          </p>
        </div>
      </div>
    );
  }

  const selectedThread = threads.find(t => t.id === selectedThreadId);

  return (
    <div className="flex-1 flex flex-col bg-background">
      {/* Thread Header */}
      <div className="border-b border-border px-6 py-4 bg-card">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold text-card-foreground">
              {selectedThread?.sellerName || 'Thread'}
            </h2>
            {selectedThread && (
              <div className="text-xs text-muted-foreground mt-1 capitalize">
                {selectedThread.sellerType}
              </div>
            )}
          </div>
          <button className="text-muted-foreground hover:text-foreground">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
        </div>
      </div>

      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto px-6 py-4 bg-muted/50">
        {loadingMessages ? (
          <div className="text-center text-muted-foreground text-sm py-8">
            Loading messages...
          </div>
        ) : messages.length === 0 ? (
          <div className="text-center text-muted-foreground text-sm py-8">
            No messages yet. Start the conversation!
          </div>
        ) : (
          <div className="space-y-4">
            {messages.map((message) => {
              const isUser = message.sender === 'user';
              const isAgent = message.sender === 'agent';
              const isSeller = message.sender === 'seller';

              return (
                <div key={message.id} className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}>
                  <div className={`max-w-[70%] ${isUser ? 'order-2' : 'order-1'}`}>
                    {/* Sender label */}
                    <div className={`text-xs text-muted-foreground mb-1 ${isUser ? 'text-right' : 'text-left'}`}>
                      {isUser && 'You'}
                      {isAgent && 'AI Agent'}
                      {isSeller && selectedThread?.sellerName}
                    </div>

                    {/* Message bubble */}
                    <div
                      className={`rounded-lg px-4 py-3 ${
                        isUser
                          ? 'bg-blue-600 text-white'
                          : isAgent
                          ? 'bg-purple-100 dark:bg-purple-950 text-foreground border border-purple-200 dark:border-purple-800'
                          : 'bg-card text-card-foreground border border-border'
                      }`}
                    >
                      <div className="whitespace-pre-wrap break-words">
                        {message.content}
                      </div>
                      <div
                        className={`text-xs mt-2 ${
                          isUser ? 'text-blue-100' : 'text-muted-foreground'
                        }`}
                      >
                        {new Date(message.timestamp).toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Message Input */}
      <div className="border-t border-border px-6 py-4 bg-card">
        <div className="flex items-end gap-3">
          <div className="flex-1">
            <textarea
              placeholder="Type message... AI will assist"
              rows={1}
              className="w-full px-4 py-3 bg-background border border-input rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-ring text-foreground"
            />
          </div>
          <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors flex items-center gap-2 cursor-pointer">
            <span>SEND</span>
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
