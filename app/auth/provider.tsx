"use strict";

import React from "react";

import { UserContext, setUserContext } from "./context";

import type { IUserInfo } from "./context";

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = React.useState<IUserInfo | null>(null);

    React.useEffect(() => {
        fetch("/api/info")
            .then((res) => {
                if (!res.ok) {
                    throw new Error(res.statusText);
                }
                return res.json();
            })
            .then((data: IUserInfo) => {
                console.log(data);
                setUser(data);
            })
            .catch((e) => {
                console.error("Failed to fetch user info:", e);
                setUser(null);
            });
    });

    return (
        <UserContext.Provider value={user}>
            <setUserContext.Provider value={setUser}>{children}</setUserContext.Provider>
        </UserContext.Provider>
    );
}
