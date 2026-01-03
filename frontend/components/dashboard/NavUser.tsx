"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import Image from "next/image"
import { useTheme } from "next-themes"
import {
  ChevronsUpDown,
  Home,
  LogOut,
  Mail,
  Settings,
  User,
  Copy,
} from "lucide-react"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { useAuth } from "@/contexts/AuthContext"
import { gmailAPI } from "@/lib/api"
import { toast } from "sonner"

export function NavUser({
  user,
  onGoToDashboard,
}: {
  user: {
    name: string
    email: string
    inboxEmail?: string
  }
  onGoToDashboard?: () => void
}) {
  const { isMobile } = useSidebar()
  const { logout } = useAuth()
  const router = useRouter()
  const { resolvedTheme } = useTheme()
  const [mounted, setMounted] = useState(false)
  const [gmailConnected, setGmailConnected] = useState(false)
  const [gmailEmail, setGmailEmail] = useState<string>()
  const [loading, setLoading] = useState(true)
  const [showEmailDialog, setShowEmailDialog] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    const fetchGmailStatus = async () => {
      try {
        const status = await gmailAPI.getStatus()
        setGmailConnected(status.connected)
        setGmailEmail(status.gmailEmail)
      } catch (error) {
        console.error('Failed to fetch Gmail status:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchGmailStatus()
  }, [])

  const handleGmailConnect = async () => {
    try {
      const { authUrl } = await gmailAPI.getAuthUrl()
      window.location.href = authUrl
    } catch (error) {
      console.error('Failed to get Gmail auth URL:', error)
      toast.error('Failed to start Gmail connection')
    }
  }

  const handleGmailDisconnect = async () => {
    try {
      await gmailAPI.disconnect()
      setGmailConnected(false)
      setGmailEmail(undefined)
      toast.success('Gmail disconnected successfully')
    } catch (error) {
      console.error('Failed to disconnect Gmail:', error)
      toast.error('Failed to disconnect Gmail')
    }
  }

  // Use dark logo in light mode, light logo in dark mode
  const logoSrc = mounted && resolvedTheme === 'light'
    ? '/logo-dark-v2.png'
    : '/logo-light-v2.png'

  return (
    <>
    <SidebarMenu>
      <SidebarMenuItem>
        <div className="flex items-center gap-2 w-full">
          <SidebarMenuButton
            onClick={() => {
              if (onGoToDashboard) {
                onGoToDashboard();
              } else {
                router.push('/dashboard');
              }
            }}
            className="shrink-0 w-10 h-10 p-0 flex items-center justify-center hover:bg-sidebar-accent"
            title="Go to Dashboard"
            asChild={false}
          >
            <Home className="size-5" />
          </SidebarMenuButton>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground flex-1"
              >
                <div className="h-8 w-auto flex items-center flex-1">
                  {mounted && (
                    <Image
                      src={logoSrc}
                      alt="Otto"
                      width={80}
                      height={24}
                      className="h-6 w-auto"
                    />
                  )}
                </div>
                <ChevronsUpDown className="ml-auto size-4" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="start"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <div className="h-8 w-8 rounded-lg bg-muted flex items-center justify-center">
                  <User className="h-4 w-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user.name}</span>
                  <span className="truncate text-xs">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              {!loading && (
                <>
                  {gmailConnected ? (
                    <DropdownMenuItem onClick={handleGmailDisconnect}>
                      <Mail />
                      <div className="flex flex-col">
                        <span>Gmail Connected ✓</span>
                        {gmailEmail && (
                          <span className="text-xs text-muted-foreground">{gmailEmail}</span>
                        )}
                      </div>
                    </DropdownMenuItem>
                  ) : (
                    <DropdownMenuItem onClick={handleGmailConnect}>
                      <Mail />
                      Connect Gmail
                    </DropdownMenuItem>
                  )}
                </>
              )}
              {user.inboxEmail && (
                <DropdownMenuItem onClick={() => setShowEmailDialog(true)}>
                  <Mail />
                  Email Forwarding
                </DropdownMenuItem>
              )}
              <DropdownMenuItem onClick={() => router.push('/settings')}>
                <Settings />
                Settings
              </DropdownMenuItem>
              <DropdownMenuItem onClick={logout}>
                <LogOut />
                Log out
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
        </div>
      </SidebarMenuItem>
    </SidebarMenu>

    {/* Email Forwarding Dialog */}
    {user.inboxEmail && (
      <Dialog open={showEmailDialog} onOpenChange={setShowEmailDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Forward Emails</DialogTitle>
            <DialogDescription>
              Forward or BCC emails from sellers to this address and they'll appear in your inbox:
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center gap-2 p-3 bg-muted rounded-md">
            <code className="flex-1 text-sm break-all">{user.inboxEmail}</code>
            <Button
              onClick={async () => {
                try {
                  await navigator.clipboard.writeText(user.inboxEmail!);
                  toast.success('Email copied to clipboard');
                } catch (err) {
                  toast.error('Failed to copy email');
                }
              }}
              variant="outline"
              size="sm"
            >
              <Copy className="h-4 w-4 mr-2" />
              Copy
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    )}
  </>
  )
}
