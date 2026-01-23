// features/games/stores/game.store.ts
import { createStore } from 'solid-js/store';
import type {CreateGameDto, GameDto, UpdateGameDto} from '../types/game.types.ts';
import { gameApi } from '../api/game.api';

interface GameState {
    games: GameDto[];
    data: GameDto[];
    isLoading: boolean;
    error: string | null;
    selectedGame: GameDto | null;
}

const initialState: GameState = {
    games: [],
    data:[],
    isLoading: false,
    error: null,
    selectedGame: null,
};

export type GameStore = {
    state: GameState;
    actions: {
        loadGames: () => Promise<void>;
        createGame: (dto: CreateGameDto) => Promise<GameDto>;
        updateGame: (id: string, dto: UpdateGameDto) => Promise<GameDto>;
        deleteGame: (id: string) => Promise<void>;
        filterGames: (search: string) => Promise<void>;
        setErrorNull: () => Promise<void>;
        setSelectedGame: (game: GameDto | null) => void;
    };
};

export const createGameStore  = () : GameStore => {
    const [state, setState] = createStore(initialState);

    const actions = {
        async loadGames() {
            setState('isLoading', true);
            setState('error', null);

            try {
                const games = await gameApi.getAllGames();
                setState('data', games)
                setState('games', games);
            } catch (error) {
                console.error('[loadGames] error', error);
                setState('error', error instanceof Error ? error.message : 'Unknown error');
            } finally {
                setState('isLoading', false);
            }
        },

        async createGame(dto: CreateGameDto) {
            setState('isLoading', true);
            // Не устанавливаем error здесь, так как ошибки обрабатываются в модальном окне

            try {
                const newGame = await gameApi.createGame(dto);
                setState('data', [...state.data, newGame]);
                setState('games', [...state.games, newGame]);
                return newGame;
            } catch (error) {
                // Пробрасываем ошибку для обработки в модальном окне
                throw error;
            } finally {
                setState('isLoading', false);
            }
        },

        async updateGame(id: string, dto: UpdateGameDto) {
            setState('isLoading', true);
            // Не устанавливаем error здесь, так как ошибки обрабатываются в модальном окне

            try {
                const updatedGame = await gameApi.updateGame(id, dto);
                setState('data', state.data.map(game => game.id === id ? updatedGame : game));
                setState('games', state.games.map(game => game.id === id ? updatedGame : game));
                if (state.selectedGame?.id === id) {
                    setState('selectedGame', updatedGame);
                }
                return updatedGame;
            } catch (error) {
                // Пробрасываем ошибку для обработки в модальном окне
                throw error;
            } finally {
                setState('isLoading', false);
            }
        },

        async deleteGame(id: string) {
            setState('isLoading', true);

            try {
                await gameApi.deleteGame(id);
                setState('data', state.data.filter(game => game.id !== id));
                setState('games', state.games.filter(game => game.id !== id));
                if (state.selectedGame?.id === id) {
                    setState('selectedGame', null);
                }
            } catch (error) {
                setState('error', error instanceof Error ? error.message : 'Failed to delete game');
                throw error;
            } finally {
                setState('isLoading', false);
            }
        },

        async filterGames(search: string): Promise<void> {
            const games = await gameApi.getAllGames();

            const gamesFiltered = games.filter(game => game.title.includes(search));

            setState('data', games)
            setState('games', gamesFiltered);
        },
        async setErrorNull(): Promise<void> {
            setState('isLoading', false);
            setState('error', null);
        },

        setSelectedGame(game: GameDto | null) {
            setState('selectedGame', game);
        }
    };

    return { state, actions };
};



// Глобальный store
export const gameStore = createGameStore();