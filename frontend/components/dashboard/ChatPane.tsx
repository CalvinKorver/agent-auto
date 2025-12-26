'use client';

import { useState, useEffect, useRef } from 'react';
import Image from 'next/image';
import { InboxMessage, Thread, messageAPI, Message, threadAPI } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { Separator } from '@/components/ui/separator';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { IconTrash } from '@tabler/icons-react';
import ArchiveConfirmDialog from './ArchiveConfirmDialog';
import TrackOfferButton from './TrackOfferButton';
import SendEmailButton from './SendEmailButton';
import TypingIndicator from './TypingIndicator';
import { cn } from '@/lib/utils';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from 'sonner';

type MessageStatus = 'sending' | 'sent' | 'error';

interface DisplayMessage extends Message {
  status?: MessageStatus;
  tempId?: string;
  errorMessage?: string;
}

interface ChatPaneProps {
  selectedThreadId: string | null;
  selectedInboxMessage?: InboxMessage | null;
  threads?: Thread[];
  onInboxMessageAssigned?: () => void;
}

export default function ChatPane({ selectedThreadId, selectedInboxMessage, threads = [], onInboxMessageAssigned }: ChatPaneProps) {
  const { user } = useAuth();
  const [showAssignModal, setShowAssignModal] = useState(false);
  const [showArchiveDialog, setShowArchiveDialog] = useState(false);
  const [assigningToThread, setAssigningToThread] = useState(false);
  const [archiving, setArchiving] = useState(false);
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const [messageInput, setMessageInput] = useState('');
  const [sendingMessage, setSendingMessage] = useState(false);
  const [showTypingIndicator, setShowTypingIndicator] = useState(false);
  const [showArchiveThreadDialog, setShowArchiveThreadDialog] = useState(false);
  const [archivingThread, setArchivingThread] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

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

  const handleArchiveConfirm = async () => {
    if (!selectedInboxMessage) return;

    setArchiving(true);
    try {
      await messageAPI.archiveInboxMessage(selectedInboxMessage.id);

      // Close dialog
      setShowArchiveDialog(false);

      // Trigger inbox refresh
      window.dispatchEvent(new Event('refreshInboxMessages'));

      if (onInboxMessageAssigned) {
        onInboxMessageAssigned();
      }
    } catch (error) {
      console.error('Failed to archive message:', error);
      alert('Failed to archive message');
    } finally {
      setArchiving(false);
    }
  };

  const handleArchiveThread = async () => {
    if (!selectedThreadId) return;

    setArchivingThread(true);
    try {
      await threadAPI.archive(selectedThreadId);

      // Close dialog
      setShowArchiveThreadDialog(false);

      // Trigger thread refresh
      window.dispatchEvent(new Event('refreshThreads'));
    } catch (error) {
      console.error('Failed to archive thread:', error);
      alert('Failed to archive thread');
    } finally {
      setArchivingThread(false);
    }
  };

  const handleSendMessage = async () => {
    if (!selectedThreadId || !messageInput.trim() || sendingMessage) {
      return;
    }

    const content = messageInput.trim();
    const tempId = `temp-${Date.now()}-${Math.random()}`;
    const threadIdAtSendTime = selectedThreadId;

    // Create optimistic message
    const optimisticMessage: DisplayMessage = {
      id: tempId,
      tempId: tempId,
      threadId: selectedThreadId,
      sender: 'user',
      content: content,
      timestamp: new Date().toISOString(),
      status: 'sent'
    };

    // Add to state immediately
    setMessages(prev => [...prev, optimisticMessage]);
    setMessageInput('');
    setShowTypingIndicator(true);
    setSendingMessage(true);

    try {
      const response = await messageAPI.createMessage(selectedThreadId, {
        content,
        sender: 'user'
      });

      // Guard: only update if still on same thread
      if (selectedThreadId !== threadIdAtSendTime) {
        return;
      }

      // Reconcile: replace optimistic message with server messages
      if ('userMessage' in response) {
        setMessages(prev => {
          const withoutOptimistic = prev.filter(m => m.tempId !== tempId);
          return [
            ...withoutOptimistic,
            { ...response.userMessage, status: 'sent' as const },
            ...(response.agentMessage ? [{ ...response.agentMessage, status: 'sent' as const }] : [])
          ];
        });
      }

      setShowTypingIndicator(false);

    } catch (error) {
      console.error('Failed to send message:', error);

      // Mark message as failed (keep visible with error indicator)
      setMessages(prev =>
        prev.map(m =>
          m.tempId === tempId
            ? { ...m, status: 'error' as const, errorMessage: 'Failed to send' }
            : m
        )
      );

      setShowTypingIndicator(false);
    } finally {
      setSendingMessage(false);
    }
  };

  // Helper function to find the most recent replyable message (seller message with externalMessageId)
  // that came before a given message index
  const findReplyableMessageId = (messageIndex: number): string | null => {
    // Look backwards from the current message to find the most recent seller message with externalMessageId
    for (let i = messageIndex - 1; i >= 0; i--) {
      const msg = messages[i];
      if (msg.sender === 'seller' && msg.externalMessageId) {
        return msg.id;
      }
    }
    return null;
  };

  const handleRetry = async (failedMessage: DisplayMessage) => {
    if (!failedMessage.tempId || !selectedThreadId) return;

    // Update to sent state (remove error)
    setMessages(prev =>
      prev.map(m =>
        m.tempId === failedMessage.tempId
          ? { ...m, status: 'sent' as const, errorMessage: undefined }
          : m
      )
    );

    setShowTypingIndicator(true);
    setSendingMessage(true);

    try {
      const response = await messageAPI.createMessage(selectedThreadId, {
        content: failedMessage.content,
        sender: 'user'
      });

      if ('userMessage' in response) {
        setMessages(prev => {
          const withoutFailed = prev.filter(m => m.tempId !== failedMessage.tempId);
          return [
            ...withoutFailed,
            { ...response.userMessage, status: 'sent' as const },
            ...(response.agentMessage ? [{ ...response.agentMessage, status: 'sent' as const }] : [])
          ];
        });
      }

      setShowTypingIndicator(false);
    } catch (error) {
      console.error('Retry failed:', error);

      setMessages(prev =>
        prev.map(m =>
          m.tempId === failedMessage.tempId
            ? { ...m, status: 'error' as const, errorMessage: 'Failed to send' }
            : m
        )
      );

      setShowTypingIndicator(false);
    } finally {
      setSendingMessage(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, showTypingIndicator]);

  // Show inbox message if selected
  if (selectedInboxMessage) {
    return (
      <>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-16">
          <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mx-2 data-[orientation=vertical]:h-4"
            />
            <h1 className="text-base font-medium">
              {selectedInboxMessage.subject || 'No Subject'}
            </h1>
          </div>
        </header>
        <div className="flex flex-1 flex-col overflow-hidden">
          <div className="flex-1 flex flex-col bg-background overflow-hidden">
            {/* Email Header */}
            <div className="border-b border-border px-6 py-4 bg-card shrink-0">
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
            <div className="flex-1 overflow-y-auto px-6 py-6 bg-muted/50" style={{ overflowY: 'auto' }}>
              <div className="bg-card rounded-lg border border-border p-6">
                <div className="whitespace-pre-wrap text-card-foreground break-words">
                  {selectedInboxMessage.content}
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="border-t border-border px-6 py-4 bg-card">
              <div className="flex gap-3">
                <Button onClick={() => setShowAssignModal(true)} disabled={archiving}>
                  Assign to Thread
                </Button>
                <Button
                  onClick={() => setShowArchiveDialog(true)}
                  disabled={archiving}
                  variant="secondary"
                >
                  Archive
                </Button>
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

                  <Button
                    onClick={() => setShowAssignModal(false)}
                    disabled={assigningToThread}
                    variant="secondary"
                    className="w-full"
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            )}

            {/* Archive Confirmation Dialog */}
            <ArchiveConfirmDialog
              open={showArchiveDialog}
              onOpenChange={setShowArchiveDialog}
              onConfirm={handleArchiveConfirm}
              archiving={archiving}
            />
          </div>
        </div>
      </>
    );
  }

  if (!selectedThreadId) {
    return (
      <>
        <header className="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-16">
          <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mx-2 data-[orientation=vertical]:h-4"
            />
            <h1 className="text-base font-medium">Chat</h1>
          </div>
        </header>
        <div className="flex flex-1 flex-col overflow-hidden">
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
        </div>
      </>
    );
  }

  const selectedThread = threads.find(t => t.id === selectedThreadId);

  return (
    <>
      <header className="flex h-16 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-16">
        <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mx-2 data-[orientation=vertical]:h-4"
          />
          <div className="flex flex-1 items-center justify-between">
            <div>
              <h1 className="text-base font-medium">
                {selectedThread?.sellerName || 'Thread'}
              </h1>
              {selectedThread && (
                <div className="text-xs text-muted-foreground capitalize">
                  {selectedThread.sellerType}
                </div>
              )}
            </div>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button className="text-muted-foreground hover:text-foreground">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => setShowArchiveThreadDialog(true)}>
                  <IconTrash size={16} />
                  Archive Thread
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>
      <div className="flex flex-1 flex-col overflow-hidden">
        <div className="flex-1 flex flex-col bg-background overflow-hidden">
          <div className="mx-auto flex h-full w-full max-w-[760px] flex-col">
            {/* Messages Area */}
            <div className="flex-1 overflow-y-auto px-6 py-4" style={{ minHeight: 0 }}>
            {loadingMessages ? (
              <div className="text-center text-muted-foreground text-sm py-8">
                Loading messages...
              </div>
            ) : messages.length === 0 ? (
              <div className="text-center text-muted-foreground text-sm py-8">
                No messages yet. Start the conversation!
              </div>
            ) : (
              <div className="space-y-12">
                {messages.map((message, messageIndex) => {
                  const isUser = message.sender === 'user';
                  const isAgent = message.sender === 'agent';
                  const isSeller = message.sender === 'seller';
                  const hasError = message.status === 'error';

                  // Agent replies render as plain text (no bubble)
                  if (isAgent) {
                    return (
                      <div
                        key={message.tempId || message.id}
                        className="flex items-start gap-3 justify-start mt-4"
                      >
                        <Image
                          src="/random-headshot_cropped.png"
                          alt="Agent avatar"
                          width={40}
                          height={40}
                          className="rounded-full border border-border bg-card shrink-0 -translate-y-2 -translate-x-3"
                        />
                        <div className="space-y-2 text-base leading-relaxed text-foreground">
                          <div className="whitespace-pre-wrap break-words">{message.content}</div>
                          <div className="text-xs text-muted-foreground">
                            {new Date(message.timestamp).toLocaleString()}
                          </div>
                          {selectedThreadId && !hasError && user?.gmailConnected && (
                            <div className="pt-1">
                              <SendEmailButton
                                messageId={message.id}
                                messageContent={message.content}
                                replyableMessageId={findReplyableMessageId(messageIndex)}
                              />
                            </div>
                          )}
                        </div>
                      </div>
                    );
                  }

                  return (
                    <div
                      key={message.tempId || message.id}
                      className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}
                    >
                      <div
                        className={
                          isUser
                            ? 'max-w-[70%] order-2'
                            : isSeller
                            ? 'w-full max-w-full order-1'
                            : 'max-w-[70%] order-1'
                        }
                      >
                        {/* Sender label */}
                        <div className={`text-xs text-muted-foreground mb-1 ${isUser ? 'text-right' : 'text-left'}`}>
                          {isUser && 'You'}
                          {isSeller && selectedThread?.sellerName}
                        </div>

                        {/* Message bubble */}
                        <div
                          className={cn(
                            'rounded-lg px-4 py-3',
                            isSeller && 'w-full',
                            isUser
                              ? 'bg-card text-card-foreground border border-border'
                              : 'bg-card text-card-foreground border border-border',
                            hasError && 'border-2 border-red-500'
                          )}
                        >
                          <div className="whitespace-pre-wrap break-words">
                            {message.content}
                          </div>

                          {/* Timestamp or status */}
                          <div
                            className={cn(
                              'text-xs mt-2',
                              'text-muted-foreground'
                            )}
                          >
                            {hasError ? (
                              <span className="text-red-500 flex items-center gap-1">
                                <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                                </svg>
                                Failed to send
                              </span>
                            ) : (
                              new Date(message.timestamp).toLocaleString()
                            )}
                          </div>

                          {/* Seller track offer button */}
                          {isSeller && selectedThreadId && !hasError && (
                            <div className="mt-2">
                              <TrackOfferButton
                                threadId={selectedThreadId}
                                messageId={message.id}
                                messageContent={message.content}
                              />
                            </div>
                          )}

                          {/* Error retry button */}
                          {hasError && (
                            <button
                              type="button"
                              onClick={() => handleRetry(message)}
                              className="mt-2 text-xs underline hover:no-underline"
                            >
                              Tap to retry
                            </button>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}

                {/* Typing Indicator */}
                {showTypingIndicator && <TypingIndicator />}

                <div ref={messagesEndRef} />
              </div>
            )}
            </div>

            {/* Message Input */}
            <div className="px-6 py-4 shrink-0">
              <div className="flex w-full items-center gap-2">
                <Input
                  value={messageInput}
                  onChange={(e) => setMessageInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Type a message... (Ctrl/Cmd+Enter to send)"
                  disabled={sendingMessage}
                />
                <div className="flex gap-2">
                  <Button onClick={handleSendMessage} disabled={sendingMessage || !messageInput.trim()}>
                    {sendingMessage ? 'Sending...' : 'Send'}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Archive Thread Confirmation Dialog */}
      <ArchiveConfirmDialog
        open={showArchiveThreadDialog}
        onOpenChange={setShowArchiveThreadDialog}
        onConfirm={handleArchiveThread}
        archiving={archivingThread}
        title="Archive this thread?"
        description="This will remove the thread from your list. All messages and tracked offers will be preserved but the thread won't be visible."
      />
    </>
  );
}
