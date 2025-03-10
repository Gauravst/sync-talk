import api from "./api";

export interface User {
  id: string;
  username: string;
  email: string;
  avatar?: string;
}

//  User Login
export const login = async (
  username: string,
  password: string,
): Promise<boolean> => {
  try {
    const response = await api.post("/auth/login", { username, password });
    return response.status === 200 || response.status === 201;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};

//  User loginWithoutAuth
export const loginWithoutAuth = async (): Promise<boolean> => {
  try {
    const response = await api.post("/auth/loginWithoutAuth");
    return response.status === 201;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};

//  Fetch User Info (Protected Route)
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
