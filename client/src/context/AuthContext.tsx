import {
  getUserInfo,
  loginUser,
  loginWithoutAuth,
  logoutUser,
} from "@/services/authService";
import { UserProps } from "@/types";
import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";

interface AuthContextType {
  user: UserProps | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  loginWithoutData: () => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<UserProps | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const currentUser = await getUserInfo();
        setUser(currentUser);
      } catch (error) {
        console.error("Failed to fetch user", error);
      } finally {
        setLoading(false);
      }
    };
    fetchUser();
  }, []);

  const login = async (username: string, password: string) => {
    try {
      const data = await loginUser(username, password);
      setUser(data);
    } catch (error) {
      console.log("Failed to login", error);
    }
  };

  const loginWithoutData = async () => {
    try {
      const data = await loginWithoutAuth();
      setUser(data);
    } catch (error) {
      console.log("Failed to login", error);
    }
  };

  const logout = async () => {
    try {
      await logoutUser();
      setUser(null);
    } catch (error) {
      console.error("Failed to logout", error);
    }
  };

  return (
    <AuthContext.Provider
      value={{ user, loading, login, loginWithoutData, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
