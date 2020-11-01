import React, { useContext } from "react";
import ArticleListContext from "../contexts/ArticleList";

const Article = (props) => {
  const c = useContext(ArticleListContext);
  const articles = c.articles;
  const setArticles = c.setArticles;

  return (
    <div className="Article">
      {"Name: " + props.article.name}
      <br></br>
      {props.article.url}
      <br></br>
      {props.article.description}
      <button style={{ float: 'right'}}
        onClick={() => {
          setArticles([...articles, props.article]);
        }}
      >
        +
      </button>
    </div>
  );
};

export default Article;
