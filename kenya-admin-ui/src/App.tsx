import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import DocLayout from "./components/layout/DocLayout";
import ApiDocs from "./views/ApiDocs";
import Playground from "./views/Playground";
import "./index.css";

function App() {
  return (
    <Router>
      <DocLayout>
        <Routes>
          <Route path="/" element={<Navigate to="/counties" replace />} />
          <Route path="/counties" element={<ApiDocs />} />
          <Route path="/constituencies" element={<ApiDocs />} />
          <Route path="/spatial" element={<ApiDocs />} />
          <Route path="/map" element={<Playground />} />
        </Routes>
      </DocLayout>
    </Router>
  );
}

export default App;
