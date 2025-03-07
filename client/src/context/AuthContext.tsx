import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import axios from "axios";

// Define the shape of the user object
interface User {
  userId: number;
  username: string;
}

// Define context value shape
interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (credentials: { username: string; password: string }) => Promise<void>;
  logout: () => Promise<void>;
}

// Create Auth Context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Auth Provider Component
export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const baseApi = import.meta.env.VITE_REACT_APP_API_URL;

  // Check login status on page load
  useEffect(() => {
    axios
      .get<User>(`${baseApi}/user`, {
        withCredentials: true,
      })
      .then((response) => setUser(response.data))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  // NOTE : this login and logout is not used and these fuction is temp
  // Login function
  const login = async (credentials: { username: string; password: string }) => {
    try {
      const response = await axios.post<User>(
        "http://localhost:8080/api/auth/login",
        credentials,
        { withCredentials: true },
      );
      setUser(response.data);
    } catch (error) {
      throw error;
    }
  };

  // Logout function
  const logout = async () => {
    try {
      await axios.post(
        "http://localhost:8080/api/auth/logout",
        {},
        { withCredentials: true },
      );
      setUser(null);
    } catch (error) {
      console.error("Logout failed", error);
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Custom hook to use AuthContext
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
