import "../App.css";

import axios from "axios";
import React, { useEffect, useState } from "react";
import ReactTooltip from "react-tooltip";

import { SERVER_URL } from "../Const";

import Article from "./Article.js";

const genTagList = (tags, checkedTags, setCheckedTags) => {
  console.log(tags);
  if (tags !== undefined) {
    return tags.map((tag, i) => {
      let desc = "N/A";

      if (tag.description.length > 0) {
        desc = tag.description;
      }

      const bkgs = ["#E8EEF2", "#D6C9C9"];

      return (
        <div
          key={"tag-" + JSON.stringify(tag)}
          className="rowC"
          style={{ backgroundColor: bkgs[i % bkgs.length] }}
        >
          <input
            type="checkbox"
            onClick={() => {
              let nc = { ...checkedTags };
              nc[tag.name] ^= true;
              setCheckedTags(nc);
            }}
          />
          <span
            data-tip
            data-for={"tag-" + tag.id}
            style={{ minWidth: "100%" }}
          >
            {tag.name}
          </span>
          <ReactTooltip id={"tag-" + tag.id} place="right" effect="solid">
            {desc}
          </ReactTooltip>
        </div>
      );
    });
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
      .get(`${SERVER_URL}/api/search/tag?orderby=name`)
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
  let makeSearch = (e) => {
    e.preventDefault();

    let ct = Object.keys(checkedTags).filter((k) => checkedTags[k]);
    axios
      .get(
        `${SERVER_URL}/api/search/article/?orderby=name&tags=${ct}&lookslike=${searchTerm}`
      )
      .then((res) => {
        setArticles(res.data);
      })
      .catch((err) => console.log(err));
  };

  return (
    <div className="ArticleSearch">
      <form onSubmit={makeSearch}>
        <input
          placeholder="search"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          required
        />
        <button onClick={makeSearch} type="submit">
          <span role="img" aria-label="search">
            🔎
          </span>
        </button>
      </form>
      <br />
      <div className="rowC">
        <div className="TagList scroll">
          {genTagList(tags, checkedTags, setCheckedTags)}
        </div>
        <div className="ArticleSearchRes scroll">
          {articles.map((element) => (
            <Article
              key={"article-" + JSON.stringify(element)}
              article={element}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

export default ArticleSearch;
