import axios from 'axios';

const API_BASE_URL = '/api/v1';

// Create axios instance
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add token to request headers
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Handle response errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Device APIs
export const deviceAPI = {
  register: (data) => apiClient.post('/devices/register', data),
  list: (params) => apiClient.get('/devices', { params }),
  get: (deviceId) => apiClient.get(`/devices/${deviceId}`),
  update: (deviceId, data) => apiClient.put(`/devices/${deviceId}`, data),
  delete: (deviceId) => apiClient.delete(`/devices/${deviceId}`),
  getHeartbeats: (deviceId, params) => apiClient.get(`/devices/${deviceId}/heartbeats`, { params })
};

// Session APIs
export const sessionAPI = {
  create: (data) => apiClient.post('/sessions', data),
  list: (params) => apiClient.get('/sessions', { params }),
  get: (sessionId) => apiClient.get(`/sessions/${sessionId}`),
  terminate: (sessionId) => apiClient.post(`/sessions/${sessionId}/terminate`)
};

// Admin APIs
export const adminAPI = {
  getStats: () => apiClient.get('/admin/stats'),
  getDevicesStatus: () => apiClient.get('/admin/devices/status')
};

// Health check
export const healthCheck = () => apiClient.get('/health');

export default apiClient;
