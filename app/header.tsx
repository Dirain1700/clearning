"use strict";

import { Link, useLocation } from "react-router";

import { useSession } from "./auth/provider";

export function Header() {
    const user = useSession();
    const location = useLocation();
    const encoded = encodeURIComponent(location.pathname + location.search);
    console.log("Current path:", location);

    return (
        <header className="sticky top-0 z-10 h-16 w-full border-b-[1px] border-[foreground] bg-[background] opacity-100 [&_*]:box-border">
            <div className="absolute box-border flex h-full w-full items-center px-4 py-8 sm:px-8">
                <Link to="/" className="text-xl font-semibold">
                    Crestal Learning
                </Link>
                <div className="ml-auto">
                    {user ? (
                        <div className="flex items-center gap-4">
                            <div>Logged in as {user.name}</div>
                            <a href="/api/logout">Logout</a>
                        </div>
                    ) : (
                        <a href={`/api/auth?from=${encoded}`}>Log in now</a>
                    )}
                </div>
            </div>
        </header>
    );
}
