import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  token: string | null;
  userInfo: any | null; 
  routers: any[] | null;
  setToken: (token: string) => void;
  setUserInfo: (info: any) => void;
  setRouters: (routers: any[]) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      userInfo: null,
      routers: null,
      setToken: (token: string) => set({ token }),
      setUserInfo: (info: any) => set({ userInfo: info }),
      setRouters: (routers: any[]) => set({ routers }),
      logout: () => {
        set({ token: null, userInfo: null, routers: null });
      },
    }),
    {
      name: 'auth-storage', 
    }
  )
);
