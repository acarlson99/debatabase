import React, { useState } from "react";
import Cookies from "universal-cookie";
// import PostArticle from "./components/PostArticle";
// import PostTag from "./components/PostTag";
// import QueryServer from "./components/QueryServer";
import ArticleSearch from "./components/ArticleSearch";
import ArticleSidebar from "./components/ArticleSidebar";
import { ArticleListProvider } from "./contexts/ArticleList";
import TopBar from "./components/TopBar";
import Present from "./components/Present";
import "./App.css";

const cookies = new Cookies();

const loadArticles = () => {
  let a = cookies.get("articles");
  if (a === undefined) {
    return [];
  }

  try {
    return JSON.parse(atob(a));
  } catch (e) {
    console.log("Unable to parse data:", a);
    console.log(e);
    return [];
  }
};

/*
need article sidebar div
search div
*/
function App() {
  const [articles, setArticles] = useState(loadArticles());
  const [page, setPage] = useState(0);
  const arr = [
    <div className="articleBody rowC">
      {/* TODO: fix height */}
      <ArticleListProvider
        value={{ articles: articles, setArticles: setArticles }}
      >
        <ArticleSearch />
        <ArticleSidebar />
      </ArticleListProvider>
    </div>,

    <div>
      <ArticleListProvider
        value={{ articles: articles, setArticles: setArticles }}
      >
        <Present />
      </ArticleListProvider>
    </div>,
  ];

  // update cookie
  try {
    cookies.set("articles", btoa(JSON.stringify(articles)));
  } catch (e) {
    console.log("COULD NOT UPDATE COOKIE", e);
  }

  return (
    <div className="App">
      {/* <header className="App-header"></header> */}
      <body className="App-Body">
        <TopBar callback={() => setPage((page + 1) % arr.length)} />
        {arr[page % arr.length]}
      </body>
    </div>
  );
}

export default App;
