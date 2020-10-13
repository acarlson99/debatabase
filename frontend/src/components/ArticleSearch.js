import React, { useState } from "react";
import axios from "axios";
import { SERVER_PORT, SERVER_HOST } from "../Const";
import Article from "./Article.js";
import "../App.css";

const ArticleSearch = () => {
  const [articles, setArticles] = useState([]);
  const [tags, setTags] = useState("");

  let makeSearch = () => {
    axios
      .get(
        `http://${SERVER_HOST}:${SERVER_PORT}/api/search/article/?orderby=name&tags=${tags}`
      )
      .then((res) => {
        setArticles(res.data);
      })
      .catch((err) => console.log(err));
  };

  return (
    <div className="ArticleSearch">
      <div className="ArticleSearchBar">
        <input
          placeholder="search"
          value={tags}
          onChange={(e) => setTags(e.target.value)}
          onKeyUp={(e) => {
            if (e.key === "Enter") {
              makeSearch();
            }
          }}
          required
        />
        <button onClick={makeSearch} type="submit">
          search
        </button>
      </div>
      <div className="ArticleSearchRes scroll">
        {articles.map((element, i) => (
          <ul key={element.id}>
            <Article article={element} />
          </ul>
        ))}
      </div>
    </div>
  );
};

export default ArticleSearch;
