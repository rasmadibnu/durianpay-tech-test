import { NavLink } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Blocks, LogOut } from "lucide-react";
import type { ReactNode } from "react";

export default function Navigation({ children }: { children: ReactNode }) {
  const { email, role, logout } = useAuth();
  const isOperation = role === "operation";

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="flex justify-between items-center px-6 py-3 border-b border-border sticky top-0 z-10 backdrop-blur-xl bg-background/80">
        <div className="flex items-center gap-3">
          <Blocks className="text-primary size-8" />
          <h1 className="text-lg font-bold text-foreground tracking-tight mr-4">
            Payment Dashboard
          </h1>
          <nav className="flex items-center gap-1 bg-muted/50 border border-border rounded-lg p-1">
            <NavLink
              to="/"
              end
              className={({ isActive }) =>
                cn(
                  "px-3 py-1.5 rounded-md text-sm font-medium transition-colors",
                  isActive
                    ? "bg-background text-foreground shadow-sm"
                    : "text-muted-foreground hover:text-foreground",
                )
              }
            >
              Payments
            </NavLink>
            {isOperation && (
              <>
                <NavLink
                  to="/merchants"
                  className={({ isActive }) =>
                    cn(
                      "px-3 py-1.5 rounded-md text-sm font-medium transition-colors",
                      isActive
                        ? "bg-background text-foreground shadow-sm"
                        : "text-muted-foreground hover:text-foreground",
                    )
                  }
                >
                  Merchants
                </NavLink>
                <NavLink
                  to="/users"
                  className={({ isActive }) =>
                    cn(
                      "px-3 py-1.5 rounded-md text-sm font-medium transition-colors",
                      isActive
                        ? "bg-background text-foreground shadow-sm"
                        : "text-muted-foreground hover:text-foreground",
                    )
                  }
                >
                  Users
                </NavLink>
              </>
            )}
          </nav>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 text-right">
            <span className="text-sm text-muted-foreground">{email}</span>
            <span className="text-xs text-primary bg-primary/10 border border-primary/20 px-2 py-0.5 rounded-full uppercase font-semibold tracking-wide">
              {role}
            </span>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={logout}
            className="text-muted-foreground hover:text-destructive"
          >
            <LogOut className="size-4" />
            Sign out
          </Button>
        </div>
      </header>
      <main>{children}</main>
    </div>
  );
}
