import { UserProps } from "@/types";
import api from "./api";

//  User Login
export const loginUser = async (
  username: string,
  password: string,
): Promise<UserProps> => {
  try {
    const response = await api.post("/auth/login", { username, password });
    return response.data;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};

//  User loginWithoutAuth
export const loginWithoutAuth = async (): Promise<UserProps> => {
  try {
    const response = await api.post("/auth/loginWithoutAuth");
    return response.data;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};

//  Fetch User Info (Protected Route)
export const getUserInfo = async (): Promise<UserProps> => {
  try {
    const response = await api.get("/user");
    return response.data;
  } catch (error) {
    console.error("Error fetching user info:", error);
    throw error;
  }
};

export const logoutUser = async (): Promise<void> => {
  try {
    const response = await api.post("/user/logout");
    console.log(response);
    return;
  } catch (error) {
    console.error("Error fetching user info:", error);
    throw error;
  }
};
