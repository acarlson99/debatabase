import React, { useState, useEffect } from "react";
import axios from "axios";
import { SERVER_PORT, SERVER_HOST } from "../Const";
import Article from "./Article.js";
import "../App.css";

const genTagList = (tags, checkedTags, setCheckedTags) => {
  if (tags !== undefined) {
    return tags.map((tag, i) => (
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
    ));
  } else {
    return <div className="QueryError">ERROR QUERYING TAGS</div>;
  }
};

const ArticleSearch = () => {
  const [checkedTags, setCheckedTags] = useState({});
  const [articles, setArticles] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [tags, setTags] = useState([]);

  useEffect(() => {
    axios
      .get(`http://${SERVER_HOST}:${SERVER_PORT}/api/search/tag`)
      .then((res) => {
        console.log("retrieved tags");
        setTags(res.data);
      })
      .catch((e) => {
        console.log("error retrieving tags: ", e);
        setTags(undefined);
      });
  }, []);

  // TODO: add search options (checkbox list of tags, lookslike, etc.)
  let makeSearch = () => {
    let ct = Object.keys(checkedTags).filter((k) => checkedTags[k]);
    axios
      .get(
        `http://${SERVER_HOST}:${SERVER_PORT}/api/search/article/?orderby=name&tags=${ct}&lookslike=${searchTerm}`
      )
      .then((res) => {
        setArticles(res.data);
      })
      .catch((err) => console.log(err));
  };

  return (
    <div className="ArticleSearch">
      <div>
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
      <br />
      <div className="rowC">
        <div className="TagList scroll">
          {genTagList(tags, checkedTags, setCheckedTags)}
        </div>
        <div className="ArticleSearchRes scroll">
          {articles.map((element, i) => (
            <Article article={element} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default ArticleSearch;
