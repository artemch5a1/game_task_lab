import {ApiHelper, getApiConfig} from "../../../app/config/api.config.ts";
import {CreateGameDto, GameDto, UpdateGameDto} from "../types/game.types.ts";

export class GameApi {
    private config = getApiConfig();

    async getAllGames(): Promise<GameDto[]> {
        const { games } = this.config.endpoints;
        const url = ApiHelper.buildUrl(this.config.baseURL, games.list);

        const response = await fetch(url, {
            method: 'GET',
            headers: ApiHelper.getHeaders(games.list),
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch games: ${response.statusText}`);
        }

        return response.json();
    }

    // Получить игру по ID
    async getGameById(id: string): Promise<GameDto> {
        const { games } = this.config.endpoints;
        const url = ApiHelper.buildUrl(this.config.baseURL, games.detail, { id });

        const response = await fetch(url, {
            method: games.detail.method,
            headers: ApiHelper.getHeaders(games.detail),
        });

        if (!response.ok) {
            if (response.status === 404) {
                throw new Error('Game not found');
            }
            throw new Error(`Failed to fetch game: ${response.statusText}`);
        }

        return response.json();
    }

    // Создать игру
    async createGame(dto: CreateGameDto): Promise<GameDto> {
        const { games } = this.config.endpoints;
        const url = ApiHelper.buildUrl(this.config.baseURL, games.create);

        // Убираем undefined поля и форматируем дату
        const payload: any = {
            title: dto.title,
            releaseDate: dto.releaseDate,
            genreId: dto.genreId,
        };
        
        if (dto.description) {
            payload.description = dto.description;
        }

        const response = await fetch(url, {
            method: games.create.method,
            headers: ApiHelper.getHeaders(games.create),
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const errorMessage = errorData.error || errorData.message || `Failed to create game: ${response.statusText}`;
            throw new Error(errorMessage);
        }

        return response.json();
    }

    // Обновить игру
    async updateGame(id: string, dto: UpdateGameDto): Promise<GameDto> {
        const { games } = this.config.endpoints;
        const url = ApiHelper.buildUrl(this.config.baseURL, games.update, { id });

        // Убираем undefined поля и форматируем дату
        const payload: any = {
            title: dto.title,
            releaseDate: dto.releaseDate,
            genreId: dto.genreId,
        };
        
        if (dto.description) {
            payload.description = dto.description;
        }

        const response = await fetch(url, {
            method: games.update.method,
            headers: ApiHelper.getHeaders(games.update),
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const errorMessage = errorData.error || errorData.message || `Failed to update game: ${response.statusText}`;
            throw new Error(errorMessage);
        }

        return response.json();
    }

    // Удалить игру
    async deleteGame(id: string): Promise<void> {
        const { games } = this.config.endpoints;
        const url = ApiHelper.buildUrl(this.config.baseURL, games.delete, { id });

        const response = await fetch(url, {
            method: games.delete.method,
            headers: ApiHelper.getHeaders(games.delete),
        });

        if (!response.ok) {
            throw new Error('Failed to delete game');
        }
    }
}

export const gameApi = new GameApi();