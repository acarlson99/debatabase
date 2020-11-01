import React, { useContext, useState } from "react";
import KeyHandler, { KEYPRESS } from "react-key-handler";
import ArticleListContext from "../contexts/ArticleList";
import ArrowSVG from "../resources/arrow.svg";

const fixURL = (s) => {
  if (!s.startsWith("http://") && !s.startsWith("https://")) {
    s = "http://" + s;
  }
  return s;
};

const PresentArticle = ({ value }) => {
  const url = fixURL(value.url);
  var u = { hostname: url, pathname: "" };
  try {
    u = new URL(url);
  } catch {
    console.log("BAD:", u);
  }
  return (
    <div className="PresentArticle">
      {"id: " + value.id}
      <br />
      url:{" "}
      <a target="_blank" rel="noopener noreferrer" href={u.href}>
        {u.hostname + u.pathname}
      </a>
      <br />
      {"name: " + value.name}
      <br />
      {"description: " + value.description}
    </div>
  );
};

// TODO: use `<svg>` tag
const ArrowIMG = ({ scroll, checkOp, rot, alt }) => (
  <div className="DivButton noHL">
    <img
      src={ArrowSVG}
      style={{
        opacity: checkOp(),
        float: "right",
        transform: `rotate(${rot}deg)`,
        width: "75%",
      }}
      alt={alt}
      onClick={scroll}
    />
  </div>
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
    <div>
      <KeyHandler
        keyEventName={KEYPRESS}
        keyValue="l"
        onKeyHandle={() => setIdx(min(idx + 1, articles.length - 1))}
      />
      <KeyHandler
        keyEventName={KEYPRESS}
        keyValue="j"
        onKeyHandle={() => setIdx(max(idx - 1, 0))}
      />

      <div className="PresentArticleList rowC">
        <ArrowIMG
          scroll={() => setIdx(max(idx - 1, 0))}
          checkOp={() => checkOpacity(0)}
          rot={180}
          alt="LeftArrow"
        />
        <PresentArticle
          value={articles[idx]}
          // style={{ marginRight: "auto", marginLeft: "auto" }}
        />
        <ArrowIMG
          scroll={() => setIdx(min(idx + 1, articles.length - 1))}
          checkOp={() => checkOpacity(articles.length - 1)}
          rot={0}
          alt="RightArrow"
        />
      </div>
    </div>
  );
};

export default Present;
