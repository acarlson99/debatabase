import React, { useState } from "react";
// import PostArticle from "./components/PostArticle";
// import PostTag from "./components/PostTag";
// import QueryServer from "./components/QueryServer";
import ArticleSearch from "./components/ArticleSearch";
import ArticleSidebar from "./components/ArticleSidebar";
import { ArticleListProvider } from "./contexts/ArticleList";
import TopBar from "./components/TopBar";
import "./App.css";

/*
need article sidebar div
search div
*/
function App() {
  const [articles, setArticles] = useState([]);

  return (
    <div className="App">
      <header className="App-header">
        <TopBar />
        <div className="articleBody rowC">
          {/* TODO: fix height */}
          <ArticleListProvider
            value={{ articles: articles, setArticles: setArticles }}
          >
            <ArticleSearch />
            <ArticleSidebar />
          </ArticleListProvider>
        </div>
      </header>
    </div>
  );
}

export default App;
