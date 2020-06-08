import React from "react";
import PostArticle from "./components/PostArticle";
import PostTag from "./components/PostTag";
import QueryServer from "./components/QueryServer";
import "./App.css";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <div className="Post-div">
          <PostTag />
          <PostArticle />
        </div>
        <div className="List-div">
          <QueryServer searchType="tag" />
          <QueryServer searchType="article" />
        </div>
      </header>
    </div>
  );
}

export default App;
