import api from "./api";

export interface User {
  id: string;
  username: string;
  email: string;
  avatar?: string;
}

// 🔹 User Login
export const login = async (
  username: string,
  password: string,
): Promise<{ token: string }> => {
  try {
    const response = await api.post("/auth/login", { username, password });

    // Store token in localStorage
    // localStorage.setItem("token", response.data.token);
    return response.data;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};

// 🔹 Fetch User Info (Protected Route)
export const getUserInfo = async (): Promise<User> => {
  try {
    const response = await api.get("/user");
    return response.data;
  } catch (error) {
    console.error("Error fetching user info:", error);
    throw error;
  }
};

// export const logout = () => {
//   localStorage.removeItem("token");
//   window.location.href = "/login";
// };
