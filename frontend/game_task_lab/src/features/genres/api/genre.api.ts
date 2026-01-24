import { ApiHelper, getApiConfig } from "../../../app/config/api.config";
import type { GenreDto } from "../types/genre.types";

export class GenreApi {
  private config = getApiConfig();

  async getAllGenres(): Promise<GenreDto[]> {
    const { genres } = this.config.endpoints;
    const url = ApiHelper.buildUrl(this.config.baseURL, genres.list);

    const response = await fetch(url, {
      method: genres.list.method,
      headers: ApiHelper.getHeaders(genres.list),
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch genres: ${response.statusText}`);
    }

    return response.json();
  }
}

export const genreApi = new GenreApi();

