import React, { useContext, useState } from "react";
import ArticleListContext from "../contexts/ArticleList";
import ArrowSVG from "../resources/arrow.svg";

const PresentArticle = ({ value }) => {
  return (
    <div>
      {"id: " + value.id}
      <br />
      {"name: " + value.name}
      <br />
      {"description: " + value.description}
    </div>
  );
};

// TODO: use `<svg>` tag
const ArrowIMG = ({ scroll, checkOp, rot, alt }) => (
  <img
    className="DivButton"
    src={ArrowSVG}
    style={{
      opacity: checkOp(),
      float: "right",
      transform: `rotate(${rot}deg)`,
    }}
    alt={alt}
    onClick={scroll}
  />
);

const min = (a, b) => {
  if (a < b) {
    return a;
  } else {
    return b;
  }
};

const max = (a, b) => {
  if (a > b) {
    return a;
  } else {
    return b;
  }
};

const Present = () => {
  const c = useContext(ArticleListContext);
  const articles = c.articles;
  console.log(ArrowSVG);
  // const setArticles = c.setArticles;

  const [idx, setIdx] = useState(0);

  const checkOpacity = (n) => {
    if (idx === n) {
      return 0.25;
    } else {
      return 1;
    }
  };

  return (
    <div className="PresentArticleList">
      <ArrowIMG
        scroll={() => setIdx(max(idx - 1, 0))}
        checkOp={() => checkOpacity(0)}
        rot={180}
        alt="LeftArrow"
      />
      <PresentArticle
        value={articles[idx]}
        style={{ marginRight: "auto", marginLeft: "auto" }}
      />
      <ArrowIMG
        scroll={() => setIdx(min(idx + 1, articles.length - 1))}
        checkOp={() => checkOpacity(articles.length - 1)}
        rot={0}
        alt="RightArrow"
      />
    </div>
  );
};

export default Present;
