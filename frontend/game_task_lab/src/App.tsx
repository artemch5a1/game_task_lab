import "./App.css";
import { Navigate, Route, Router, useNavigate } from "@solidjs/router";
import { AppLayout } from "./assets/AppLayout.tsx";
import { GamesPage } from "./features/games/pages/GamesPage.tsx";
import { LoginPage } from "./features/auth/pages/LoginPage";
import { authStore } from "./features/auth/store/auth.store";
import { createEffect, Show } from "solid-js";

function App() {
    const RequireAuth = (props: { children: any }) => {
        const navigate = useNavigate();
        createEffect(() => {
            if (!authStore.actions.isAuthenticated()) {
                navigate("/login", { replace: true });
            }
        });

        return (
            <Show when={authStore.actions.isAuthenticated()}>
                {props.children}
            </Show>
        );
    };

    return (
        <Router>
            <Route path="/login" component={() => (
                authStore.actions.isAuthenticated()
                    ? <Navigate href="/" />
                    : <LoginPage />
            )} />
            <Route path="/" component={() => (
                <RequireAuth>
                    <AppLayout showSidebar={true}>
                        <GamesPage />
                    </AppLayout>
                </RequireAuth>
            )} />
        </Router>
    );
}

export default App;