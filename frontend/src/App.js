import React, { useState } from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import Cookies from "universal-cookie";
// import PostArticle from "./components/PostArticle";
// import PostTag from "./components/PostTag";
// import QueryServer from "./components/QueryServer";
import ArticleSearch from "./components/ArticleSearch";
import ArticleSidebar from "./components/ArticleSidebar";
import { ArticleListProvider } from "./contexts/ArticleList";
import TopBar from "./components/TopBar";
import Present from "./components/Present";
import Upload from "./components/Upload";
import ArticlePage from "./components/ArticlePage";
import PostArticle from "./components/PostArticle";
import PostTag from "./components/PostTag";
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

const Search = () => {
  const [articles, setArticles] = useState(loadArticles());
  return (
    <div className="articleBody rowC">
      <ArticleListProvider
        value={{
          articles: articles,
          setArticles: setArticles,
        }}
      >
        <ArticleSearch />
        <ArticleSidebar />
      </ArticleListProvider>
    </div>
  );
};

const LoadPresent = () => {
  const [articles, setArticles] = useState(loadArticles());
  return (
    <ArticleListProvider
      value={{
        articles: articles,
        setArticles: setArticles,
      }}
    >
      <Present />
    </ArticleListProvider>
  );
};

/*
need article sidebar div
search div
*/
function App() {
  const [articles] = useState(loadArticles());
  // const [page, setPage] = useState(0);

  // update cookie
  try {
    cookies.set("articles", btoa(JSON.stringify(articles)));
  } catch (e) {
    console.log("COULD NOT UPDATE COOKIE", e);
  }

  return (
    <div className="App">
      <Router>
        {/* <header className="App-header"></header> */}
        <body className="App-Body">
          <TopBar />
          <Route exact path="/" component={Search} />
          <Route exact path="/present" component={LoadPresent} />
          <Route exact path="/upload/article" component={PostArticle} />
          <Route exact path="/upload/tag" component={PostTag} />
          <Route exact path="/upload" component={Upload} />
          <Route path="/article" component={ArticlePage} />
        </body>
      </Router>
    </div>
  );
}

export default App;
