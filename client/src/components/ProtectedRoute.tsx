import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const ProtectedRoute = () => {
  const { user } = useAuth();
  const location = useLocation();

  if (user && location.pathname === "/login") {
    return <Navigate to="/chat" replace />;
  }

  return user ? <Outlet /> : <Navigate to="/login" replace />;
};

export default ProtectedRoute;
