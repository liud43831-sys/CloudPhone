import { create } from 'zustand';

export const useAuthStore = create((set) => ({
  token: localStorage.getItem('token') || null,
  user: null,

  setToken: (token) => {
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
    set({ token });
  },

  setUser: (user) => set({ user }),

  logout: () => {
    localStorage.removeItem('token');
    set({ token: null, user: null });
  },

  isAuthenticated: () => !!localStorage.getItem('token')
}));

export const useDeviceStore = create((set) => ({
  devices: [],
  loading: false,
  error: null,

  setDevices: (devices) => set({ devices }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),

  addDevice: (device) =>
    set((state) => ({
      devices: [...state.devices, device]
    })),

  updateDevice: (id, updatedDevice) =>
    set((state) => ({
      devices: state.devices.map((d) => (d.id === id ? updatedDevice : d))
    })),

  removeDevice: (id) =>
    set((state) => ({
      devices: state.devices.filter((d) => d.id !== id)
    }))
}));

export const useSessionStore = create((set) => ({
  sessions: [],
  loading: false,
  error: null,

  setSessions: (sessions) => set({ sessions }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),

  addSession: (session) =>
    set((state) => ({
      sessions: [...state.sessions, session]
    })),

  updateSession: (id, updatedSession) =>
    set((state) => ({
      sessions: state.sessions.map((s) => (s.id === id ? updatedSession : s))
    })),

  removeSession: (id) =>
    set((state) => ({
      sessions: state.sessions.filter((s) => s.id !== id)
    }))
}));
