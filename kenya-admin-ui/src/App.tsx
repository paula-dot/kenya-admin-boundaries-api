import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Playground from "./views/Playground";
import "./index.css";
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Playground />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
