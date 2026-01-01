'use client';

import { useState } from 'react';
import { twilioAPI } from '@/lib/api';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { toast } from 'sonner';

interface SendSMSButtonProps {
  messageId: string;
  messageContent: string;
  replyableMessageId: string | null; // The seller message ID to reply to
  phoneNumber: string; // The dealer's phone number to send to
  onSuccess?: () => void;
}

export default function SendSMSButton({ 
  messageId, 
  messageContent, 
  replyableMessageId,
  phoneNumber,
  onSuccess 
}: SendSMSButtonProps) {
  const [sending, setSending] = useState(false);
  const [showDialog, setShowDialog] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const formatPhoneNumber = (phone: string) => {
    // Format E.164 to (XXX) XXX-XXXX
    const cleaned = phone.replace(/\D/g, '');
    if (cleaned.length === 11 && cleaned.startsWith('1')) {
      return `(${cleaned.slice(1, 4)}) ${cleaned.slice(4, 7)}-${cleaned.slice(7)}`;
    } else if (cleaned.length === 10) {
      return `(${cleaned.slice(0, 3)}) ${cleaned.slice(3, 6)}-${cleaned.slice(6)}`;
    }
    return phone;
  };

  const handleSendSMS = async () => {
    if (!replyableMessageId) {
      toast.error('No SMS message found to reply to in this thread');
      return;
    }

    setSending(true);
    setError(null);
    
    try {
      await twilioAPI.sendSMS(replyableMessageId, messageContent);
      toast.success('SMS sent successfully!');
      setShowDialog(false);
      if (onSuccess) {
        onSuccess();
      }
    } catch (error: unknown) {
      const errorMessage = (error as { response?: { data?: { error?: string } } }).response?.data?.error || 'Failed to send SMS';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setSending(false);
    }
  };

  if (!replyableMessageId || !phoneNumber) {
    return null;
  }

  return (
    <>
      <Button
        onClick={() => setShowDialog(true)}
        variant="outline"
        size="sm"
        className="min-w-[100px]"
      >
        Send SMS
      </Button>

      <Dialog open={showDialog} onOpenChange={setShowDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Send SMS</DialogTitle>
            <DialogDescription>
              Are you sure you want to send this message to {formatPhoneNumber(phoneNumber)}?
            </DialogDescription>
          </DialogHeader>
          
          {error && (
            <div className="text-sm text-red-500 bg-red-50 dark:bg-red-900/20 p-3 rounded">
              {error}
            </div>
          )}

          <div className="py-4">
            <p className="text-sm text-muted-foreground mb-2">Message:</p>
            <div className="bg-muted p-3 rounded text-sm whitespace-pre-wrap">
              {messageContent}
            </div>
          </div>

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setShowDialog(false);
                setError(null);
              }}
              disabled={sending}
            >
              Discard
            </Button>
            <Button
              onClick={handleSendSMS}
              disabled={sending}
            >
              {sending ? 'Sending...' : 'Send'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

