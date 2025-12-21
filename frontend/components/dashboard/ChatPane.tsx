'use client';

import { useState, useEffect, useRef } from 'react';
import { InboxMessage, Thread, messageAPI, Message, threadAPI } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { Separator } from '@/components/ui/separator';
import ArchiveConfirmDialog from './ArchiveConfirmDialog';
import TrackOfferButton from './TrackOfferButton';

interface ChatPaneProps {
  selectedThreadId: string | null;
  selectedInboxMessage?: InboxMessage | null;
  threads?: Thread[];
  onInboxMessageAssigned?: () => void;
}

export default function ChatPane({ selectedThreadId, selectedInboxMessage, threads = [], onInboxMessageAssigned }: ChatPaneProps) {
  const [showAssignModal, setShowAssignModal] = useState(false);
  const [showArchiveDialog, setShowArchiveDialog] = useState(false);
  const [assigningToThread, setAssigningToThread] = useState(false);
  const [archiving, setArchiving] = useState(false);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const [messageInput, setMessageInput] = useState('');
  const [sendingMessage, setSendingMessage] = useState(false);
  const [showGearMenu, setShowGearMenu] = useState(false);
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

      // Close dialogs
      setShowArchiveThreadDialog(false);
      setShowGearMenu(false);

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
    setSendingMessage(true);

    try {
      const response = await messageAPI.createMessage(selectedThreadId, {
        content,
        sender: 'user'
      });

      // Add both user and agent messages to local state
      if ('userMessage' in response) {
        setMessages(prev => [
          ...prev,
          response.userMessage,
          ...(response.agentMessage ? [response.agentMessage] : [])
        ]);
      }

      // Clear input
      setMessageInput('');
    } catch (error) {
      console.error('Failed to send message:', error);
      alert('Failed to send message. Please try again.');
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
  }, [messages]);

  useEffect(() => {
    const handleClickOutside = () => setShowGearMenu(false);
    if (showGearMenu) {
      window.addEventListener('click', handleClickOutside);
      return () => window.removeEventListener('click', handleClickOutside);
    }
  }, [showGearMenu]);

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
                <button
                  onClick={() => setShowAssignModal(true)}
                  disabled={archiving}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Assign to Thread
                </button>
                <button
                  onClick={() => setShowArchiveDialog(true)}
                  disabled={archiving}
                  className="bg-secondary hover:bg-secondary/80 text-secondary-foreground px-6 py-2 rounded-lg font-medium transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
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
            <div className="relative">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  setShowGearMenu(!showGearMenu);
                }}
                className="text-muted-foreground hover:text-foreground"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </button>

              {showGearMenu && (
                <div className="absolute right-0 mt-2 w-48 bg-popover rounded-md shadow-lg border border-border z-10">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      setShowGearMenu(false);
                      setShowArchiveThreadDialog(true);
                    }}
                    className="w-full text-left px-4 py-2 text-sm hover:bg-accent hover:text-accent-foreground rounded-md"
                  >
                    Archive Thread
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </header>
      <div className="flex flex-1 flex-col overflow-hidden">
        <div className="flex-1 flex flex-col bg-background overflow-hidden">
          {/* Messages Area */}
          <div className="flex-1 overflow-y-auto px-6 py-4 bg-muted/50" style={{ minHeight: 0 }}>
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
                          {isSeller && selectedThreadId && (
                            <div className="mt-2">
                              <TrackOfferButton
                                threadId={selectedThreadId}
                                messageId={message.id}
                                messageContent={message.content}
                              />
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
                <div ref={messagesEndRef} />
              </div>
            )}
          </div>

          {/* Message Input */}
          <div className="border-t border-border px-6 py-4 bg-card shrink-0">
            <div className="flex w-full items-center gap-2">
              <Input
                value={messageInput}
                onChange={(e) => setMessageInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Type a message... (Try asking about the seller's offer)"
                disabled={sendingMessage}
              />
              <Button
                onClick={handleSendMessage}
                disabled={sendingMessage || !messageInput.trim()}
                className="bg-green-600 hover:bg-green-700 text-white"
              >
                {sendingMessage ? 'Sending...' : 'Send'}
              </Button>
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
