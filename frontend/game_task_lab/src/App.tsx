import "./App.css";
import { AppLayout } from "./assets/AppLayout.tsx";
import { GamesPage } from "./features/games/pages/GamesPage";

function App() {
  return (
    <AppLayout>
      <GamesPage />
    </AppLayout>
  );
}

export default App;
