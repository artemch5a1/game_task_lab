
export interface ApiEndpoint {
    path: string;
    method: 'GET' | 'POST' | 'PUT' | 'DELETE';
    requiresAuth: boolean;
}

export interface ApiConfig {
    baseURL: string;
    timeout: number;
    endpoints: {
        games: {
            list: ApiEndpoint;
            create: ApiEndpoint;
            detail: ApiEndpoint;
            update: ApiEndpoint;
            delete: ApiEndpoint;
        };
        genres: {
            list: ApiEndpoint;
        };
        auth: {
            login: ApiEndpoint;
            register: ApiEndpoint;
        };
    };
}

export const apiEndpoints = {
    games: {
        list: {
            path: '/games',
            method: 'GET',
            requiresAuth: false,
        },
        create: {
            path: '/games',
            method: 'POST',
            requiresAuth: true,
        },
        detail: {
            path: '/games/:id',
            method: 'GET',
            requiresAuth: false,
        },
        update: {
            path: '/games/:id',
            method: 'PUT',
            requiresAuth: true,
        },
        delete: {
            path: '/games/:id',
            method: 'DELETE',
            requiresAuth: true,
        },
    },
    genres: {
        list: {
            path: '/genres',
            method: 'GET',
            requiresAuth: false,
        },
    },
    auth: {
        login: {
            path: '/auth/login',
            method: 'POST',
            requiresAuth: false,
        },
        register: {
            path: '/auth/register',
            method: 'POST',
            requiresAuth: false,
        },
    },
} as const;

export const createApiConfig = (): ApiConfig => {
    const baseURL = import.meta.env.VITE_API_URL?.replace(/\/$/, '') || 'http://localhost:8080';

    return {
        baseURL,
        timeout: 10000,
        endpoints: apiEndpoints,
    };
};

export class ApiHelper {
    static buildUrl(baseURL: string, endpoint: ApiEndpoint, params?: Record<string, string>): string {
        let path = endpoint.path;

        if (params) {
            for (const [key, value] of Object.entries(params)) {
                path = path.replace(`:${key}`, encodeURIComponent(value));
            }
        }

        return `${baseURL}${path}`;
    }

    static getGameUrl(config: ApiConfig, gameId: string): string {
        return this.buildUrl(config.baseURL, config.endpoints.games.detail, { id: gameId });
    }

    static getHeaders(endpoint: ApiEndpoint, extraHeaders?: Record<string, string>): HeadersInit {
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            ...extraHeaders,
        };

        if (endpoint.requiresAuth) {
            const token = localStorage.getItem('auth_token');
            if (token) {
                headers['Authorization'] = `Bearer ${token}`;
            }
        }

        return headers;
    }
}

let apiConfigInstance: ApiConfig | null = null;

export const getApiConfig = (): ApiConfig => {
    if (!apiConfigInstance) {
        apiConfigInstance = createApiConfig();
        console.log('API Config initialized:', {
            baseURL: apiConfigInstance.baseURL,
            endpoints: Object.keys(apiConfigInstance.endpoints),
        });
    }
    return apiConfigInstance;
};

export const useApiConfig = () => {
    const config = getApiConfig();

    return {
        config,
        endpoints: config.endpoints,
        helper: ApiHelper,
    };
};

// Экспорт по умолчанию
export default getApiConfig;