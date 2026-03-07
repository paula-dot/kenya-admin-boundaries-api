import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import DocLayout from "./components/layout/DocLayout";
import ApiDocs from "./views/ApiDocs";
import Playground from "./views/Playground";
import "./index.css";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/counties" replace />} />
        <Route path="/counties" element={<DocLayout><ApiDocs /></DocLayout>} />
        <Route path="/constituencies" element={<DocLayout><ApiDocs /></DocLayout>} />
        <Route path="/sub-counties" element={<DocLayout><ApiDocs /></DocLayout>} />
        <Route path="/map" element={<DocLayout noPadding><Playground /></DocLayout>} />
      </Routes>
    </Router>
  );
}

export default App;
