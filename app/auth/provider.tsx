"use strict";

import React from "react";

interface IUserInfo {
    sub: string;
    name: string;
    avatarUrl: string;
}

const UserContext = React.createContext<IUserInfo | null>(null);
const setUserContext = React.createContext<React.Dispatch<React.SetStateAction<IUserInfo | null>>>(() => null);

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

export function useSession() {
    return React.useContext(UserContext);
}

export function useSetSession() {
    return React.useContext(setUserContext);
}
