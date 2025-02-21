import axios from "axios";

// const API_URL = process.env.REACT_APP_API_URL || "https://your-api-url.com/api";

// const API_URL = import.meta.env.REACT_APP_API_URL;
const API_URL = "http://localhost:8080/api";
const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
  withCredentials: true,
});

export default api;
