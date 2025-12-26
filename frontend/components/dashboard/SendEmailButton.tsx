'use client';

import { useState } from 'react';
import { gmailAPI } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';

interface SendEmailButtonProps {
  messageId: string;
  messageContent: string;
  replyableMessageId: string | null; // The seller message ID to reply to
}

export default function SendEmailButton({ messageId, messageContent, replyableMessageId }: SendEmailButtonProps) {
  const [sending, setSending] = useState(false);

  const handleSendEmail = async () => {
    if (!replyableMessageId) {
      toast.error('No email message found to reply to in this thread');
      return;
    }

    setSending(true);
    try {
      await gmailAPI.replyViaGmail(replyableMessageId, messageContent);
      toast.success('Email sent via Gmail!');
    } catch (error: unknown) {
      const errorMessage = (error as { response?: { data?: { error?: string } } }).response?.data?.error || '';

      // Handle specific error cases from backend
      if (errorMessage === 'gmail not connected') {
        toast.error('Gmail not connected. Connect in your profile menu.');
      } else if (errorMessage === 'message not found') {
        toast.error('Message not found');
      } else if (errorMessage === 'message was not received via email') {
        toast.error('This message was not received via email and cannot be replied to');
      } else {
        toast.error('Failed to send email via Gmail');
      }
    } finally {
      setSending(false);
    }
  };

  if (!replyableMessageId) {
    return null;
  }

  return (
    <Button
      onClick={handleSendEmail}
      disabled={sending}
      variant="outline"
      size="sm"
      className="mt-1"
    >
      {sending ? 'Sending...' : 'Send Email'}
    </Button>
  );
}

