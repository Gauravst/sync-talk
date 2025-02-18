import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import ChatPage from "./pages/ChatPage";
import RoomListPage from "./pages/ChatListPage";

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/rooms" element={<RoomListPage />} />
        <Route path="/chat/:name" element={<ChatPage />} />
      </Routes>
    </Router>
  );
};

export default App;
