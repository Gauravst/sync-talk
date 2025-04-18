import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { ThemeProvider } from "@/components/theme-provider";
import LoginPage from "./pages/LoginPage";
import ChatPage from "./pages/ChatPage";
import ProtectedRoute from "./components/ProtectedRoute";
import HomePage from "./pages/HomePage";
import RoomsPage from "./pages/RoomsPage";
import { NotFoundPage } from "./pages/NotFoundPage";

const App = () => {
  return (
    <ThemeProvider defaultTheme="dark" forcedTheme="dark">
      <AuthProvider>
        <Router>
          <Routes>
            {/* Public Routes */}
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />

            {/* Protected Routes */}
            <Route element={<ProtectedRoute />}>
              <Route path="/rooms" element={<RoomsPage />} />
              <Route path="/chat/:name" element={<ChatPage />} />
              <Route path="/chat" element={<ChatPage />} />
            </Route>

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App;
