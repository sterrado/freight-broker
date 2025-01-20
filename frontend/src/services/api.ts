// src/services/api.ts

import axios from 'axios';
import { Load, LoadsResponse } from '../types/load.types';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL,
});

// Add auth token to all requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authService = {
  login: async () => {
    const response = await api.post('/auth/login', {
      username: 'admin',
      password: 'password'
    });
    localStorage.setItem('token', response.data.token);
    return response.data;
  }
};

export const loadService = {
  getLoads: async (page: number = 1, size: number = 10): Promise<LoadsResponse> => {
    const response = await api.get<LoadsResponse>(`/loads?page=${page}&size=${size}`);
    return response.data;
  },

  getLoadById: async (id: string): Promise<Load> => {
    const response = await api.get<Load>(`/loads/${id}`);
    return response.data;
  },

  createLoad: async (load: Omit<Load, 'id' | 'createdAt' | 'updatedAt'>): Promise<Load> => {
    const response = await api.post<Load>('/loads', load);
    return response.data;
  },
};