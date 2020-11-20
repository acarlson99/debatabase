import React, { useState, useEffect } from "react";
import axios from "axios";
import { SERVER_URL } from "../Const";

const ArticlePage = () => {
  const pn = window.location.pathname.split("/");
  // TODO: check if num
  const articleNum = pn[pn.length - 1];
  const [info, setInfo] = useState();

  useEffect(() => {
    axios
      .get(`${SERVER_URL}/api/search/article/${articleNum}`)
      .then((res) => {
        setInfo(res.data);
      })
      .catch((e) => {
        console.log("error retrieving article: ", e);
        // TODO: handle error better than this
      });
  }, [articleNum]);

  // TODO: make `edit` and `delete` options
  // TODO: add similar page for tags
  return <div>{JSON.stringify(info)}</div>;
};

export default ArticlePage;
