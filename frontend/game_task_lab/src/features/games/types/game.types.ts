import {DateTime, UUID} from "../../../types/global";

export interface GameDto {
    id: UUID;
    title: string;
    description: string;
    releaseDate: DateTime;
    genreId: UUID;
}

export interface GameDtoWithStats extends GameDto {
    averageRating: number;
    ratingCount: number;
}

export interface CreateGameDto {
    title: string;
    description?: string;
    releaseDate: DateTime;
    genreId: UUID;
}

export interface UpdateGameDto extends CreateGameDto {
    id: UUID;
}