"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import {
  ChevronsUpDown,
  LogOut,
  Mail,
  Settings,
  User,
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
}: {
  user: {
    name: string
    email: string
  }
}) {
  const { isMobile } = useSidebar()
  const { logout } = useAuth()
  const router = useRouter()
  const [gmailConnected, setGmailConnected] = useState(false)
  const [gmailEmail, setGmailEmail] = useState<string>()
  const [loading, setLoading] = useState(true)

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

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="h-8 w-8 rounded-lg bg-muted flex items-center justify-center">
                <User className="h-4 w-4" />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-lg font-bold">Lolo AI</span>
                <span className="truncate text-xs">{user.email}</span>
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
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
