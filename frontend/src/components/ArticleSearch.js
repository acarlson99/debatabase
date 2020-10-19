import React, { useState, useEffect } from "react";
import axios from "axios";
import { SERVER_PORT, SERVER_HOST } from "../Const";
import Article from "./Article.js";
import "../App.css";

const ArticleSearch = () => {
  const [checkedTags, setCheckedTags] = useState({});
  const [articles, setArticles] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [tags, setTags] = useState([]);

  useEffect(() => {
    axios
      .get(`http://${SERVER_HOST}:${SERVER_PORT}/api/search/tag`)
      .then((res) => {
        console.log("retrieving tags");
        setTags(res.data);
      })
      .catch((e) => console.log(e));
  }, []);

  // TODO: add search options (checkbox list of tags, lookslike, etc.)
  let makeSearch = () => {
    let tags = Object.keys(checkedTags).filter((k) => checkedTags[k]);
    axios
      .get(
        `http://${SERVER_HOST}:${SERVER_PORT}/api/search/article/?orderby=name&tags=${tags}&lookslike=${searchTerm}`
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
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          onKeyUp={(e) => {
            if (e.key === "Enter") {
              makeSearch();
            }
          }}
          required
        />
        <button onClick={makeSearch} type="submit">
          <span role="img" aria-label="search">
            🔎
          </span>
        </button>
      </div>
      <div>
        <div className="TagList">
          {tags.map((tag, i) => (
            <div>
              <input
                type="checkbox"
                onClick={() => {
                  let nc = { ...checkedTags };
                  nc[tag.name] ^= true;
                  setCheckedTags(nc);
                }}
              />
              {tag.name}
            </div>
          ))}
        </div>
        <div className="ArticleSearchRes scroll">
          {articles.map((element, i) => (
            <ul key={element.id}>
              <Article article={element} />
            </ul>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ArticleSearch;
