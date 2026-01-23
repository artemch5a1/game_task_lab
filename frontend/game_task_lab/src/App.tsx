import "./App.css";
import { Route, Router } from "@solidjs/router";
import { AppLayout } from "./assets/AppLayout.tsx";
import { GamesPage } from "./features/games/pages/GamesPage.tsx";

function App() {
    return (
        <Router>
            <Route path="/" component={() => (
                <AppLayout showSidebar={true}>
                    <GamesPage />
                </AppLayout>
            )} />
        </Router>
    );
}

export default App;