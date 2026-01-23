// features/games/components/GameList.tsx
import { createEffect, Show } from 'solid-js';
import { gameStore } from '../store/game.store';
import './GameList.css';

const formatDate = (dateString: string): string => {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    } catch {
        return dateString;
    }
};

export const GameList = () => {
    const { state, actions } = gameStore;

    createEffect(() => {
        actions.loadGames();
    });

    return (
        <div class="game-list-container">
            <Show when={state.isLoading}>
                <div class="loading-container">
                    <div class="spinner"></div>
                    <p>Загрузка игр...</p>
                </div>
            </Show>

            <Show when={state.error}>
                <div class="error-container">
                    <p class="error-message">Ошибка: {state.error}</p>
                    <button 
                        class="retry-button"
                        onClick={() => actions.loadGames()}
                    >
                        Попробовать снова
                    </button>
                </div>
            </Show>

            <Show when={!state.isLoading && !state.error && state.games.length === 0}>
                <div class="empty-state">
                    <p>Игры не найдены</p>
                </div>
            </Show>

            <Show when={!state.isLoading && !state.error && state.games.length > 0}>
                <div class="games-grid">
                    {state.games.map(game => (
                        <div class="game-card" id={game.id}>
                            <div class="game-card-header">
                                <h3 class="game-title">{game.title}</h3>
                            </div>
                            <div class="game-card-body">
                                <Show when={game.description}>
                                    <p class="game-description">{game.description}</p>
                                </Show>
                                <div class="game-meta">
                                    <span class="game-date">
                                        📅 {formatDate(game.releaseDate)}
                                    </span>
                                </div>
                            </div>
                            <div class="game-card-footer">
                                <button
                                    class="delete-button"
                                    onClick={() => {
                                        if (confirm(`Вы уверены, что хотите удалить игру "${game.title}"?`)) {
                                            actions.deleteGame(game.id);
                                        }
                                    }}
                                    disabled={state.isLoading}
                                >
                                    Удалить
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            </Show>
        </div>
    );
};