import axios from "axios";

const API_URL = import.meta.env.VITE_REACT_APP_API_URL;

const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
  withCredentials: true,
});

export default api;
