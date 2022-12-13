import React from "react";
import { Navigate, Route } from "react-router-dom";
import AccountService from "../service/AccountService";

export type RequireAuthProps = {
    authenticationPath: string;
    outlet: JSX.Element;
};

export default function ProtectedRoute({ authenticationPath, outlet }: RequireAuthProps) {
    if (AccountService.getCurrentUser() === null) {
        return <Navigate to={{ pathname: authenticationPath }} />;
    } else {
        return outlet;
    }
};