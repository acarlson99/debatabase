import React, { useState } from "react";
import axios from "axios";
import { SERVER_URL } from "../Const";

const QueryServer = (props) => {
  const [articles, setArticles] = useState([]);
  const [tags, setTags] = useState("");
  const [limit, setLimit] = useState("");
  const [offset, setOffset] = useState("");
  const [lookslike, setLookslike] = useState("");

  return (
    <div className="List">
      <input
        placeholder="tags"
        value={tags}
        onChange={(e) => setTags(e.target.value)}
      />
      <br />
      <input
        placeholder="limit"
        value={limit}
        onChange={(e) => setLimit(e.target.value)}
      />
      <br />
      <input
        placeholder="offset"
        value={offset}
        onChange={(e) => setOffset(e.target.value)}
      />
      <br />
      <input
        placeholder="lookslike"
        value={lookslike}
        onChange={(e) => setLookslike(e.target.value)}
      />
      <br />
      <button
        onClick={() => {
          axios
            .get(
              `${SERVER_URL}/api/search/${props.searchType}?orderby=name&tags=${tags}&limit=${limit}&offset=${offset}&lookslike=${lookslike}`
            )
            .then((res) => setArticles(res.data))
            .catch((err) => console.log(err));
        }}
      >
        Update {props.searchType}
      </button>
      <div>
        {articles.map((element) => (
          <ul key={"query-res" + JSON.stringify(element)}>
            {JSON.stringify(element)}
          </ul>
        ))}
      </div>
    </div>
  );
};

export default QueryServer;
