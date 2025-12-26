'use client';

export default function TypingIndicator() {
  return (
    <div className="flex justify-start">
      <div className="max-w-[70%]">
        <div className="text-xs text-muted-foreground mb-1 text-left">
          AI Agent
        </div>
        <div className="rounded-lg px-4 py-3 bg-purple-100 dark:bg-purple-950 border border-purple-200 dark:border-purple-800">
          <div className="flex items-center gap-1">
            <span className="text-muted-foreground text-sm mr-2">AI is thinking</span>
            <div className="flex gap-1">
              <div className="w-2 h-2 bg-purple-400 rounded-full animate-bounce [animation-delay:-0.3s]"></div>
              <div className="w-2 h-2 bg-purple-400 rounded-full animate-bounce [animation-delay:-0.15s]"></div>
              <div className="w-2 h-2 bg-purple-400 rounded-full animate-bounce"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
