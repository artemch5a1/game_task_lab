import { ApiHelper, getApiConfig } from "../../../app/config/api.config";

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export class AuthApi {
  private config = getApiConfig();

  async login(dto: LoginRequest): Promise<LoginResponse> {
    const { auth } = this.config.endpoints;
    const url = ApiHelper.buildUrl(this.config.baseURL, auth.login);

    const response = await fetch(url, {
      method: auth.login.method,
      headers: ApiHelper.getHeaders(auth.login),
      body: JSON.stringify(dto),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMessage =
        errorData.error || errorData.message || `Login failed: ${response.statusText}`;
      throw new Error(errorMessage);
    }

    return response.json();
  }

  async register(dto: RegisterRequest): Promise<LoginResponse> {
    const { auth } = this.config.endpoints;
    const url = ApiHelper.buildUrl(this.config.baseURL, auth.register);

    const response = await fetch(url, {
      method: auth.register.method,
      headers: ApiHelper.getHeaders(auth.register),
      body: JSON.stringify(dto),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMessage =
        errorData.error || errorData.message || `Register failed: ${response.statusText}`;
      throw new Error(errorMessage);
    }

    return response.json();
  }
}

export const authApi = new AuthApi();

