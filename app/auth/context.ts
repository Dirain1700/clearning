"use strict";

import React from "react";

export interface IUserInfo {
    sub: string;
    name: string;
    avatarUrl: string;
}

export const UserContext = React.createContext<IUserInfo | null>(null);
export const setUserContext = React.createContext<React.Dispatch<React.SetStateAction<IUserInfo | null>>>(() => null);

export function useSession() {
    return React.useContext(UserContext);
}

export function useSetSession() {
    return React.useContext(setUserContext);
}
