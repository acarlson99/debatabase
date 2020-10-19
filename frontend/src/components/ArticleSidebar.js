import React, { useContext } from "react";
import {
  SortableContainer,
  SortableElement,
  sortableHandle,
} from "react-sortable-hoc";
import arrayMove from "array-move";
import ArticleListContext from "../contexts/ArticleList";

// TODO: fix `findDOMNode` error (prolly use refs)

const DragHandle = sortableHandle(() => <span className="DragHandle">::</span>);

const SortableItem = SortableElement(({ value, deleteArticle }) => (
  <div className="SortableElement Article">
    <DragHandle />
    {JSON.stringify(value)}
    <button onClick={deleteArticle}>-</button>
  </div>
));

const SortableList = SortableContainer(({ items, deleteArticle }) => (
  <div>
    {items.map((value, index) => (
      <SortableItem
        key={`item-${index}`}
        index={index}
        value={value}
        deleteArticle={() => {
          console.log("deleting " + index);
          deleteArticle(index);
        }}
      />
    ))}
  </div>
));

const ArticleSidebar = () => {
  const c = useContext(ArticleListContext);
  const articles = c.articles;
  const setArticles = c.setArticles;

  return (
    <div className="ArticleSidebar">
      {/* TODO: confirm reset */}
      <button onClick={() => setArticles([])}>clear</button>{" "}
      <SortableList
        deleteArticle={(i) => {
          let newA = [...articles];
          newA.splice(i, 1);
          setArticles(newA);
        }}
        items={articles}
        onSortEnd={({ oldIndex, newIndex }) =>
          setArticles(arrayMove(articles, oldIndex, newIndex))
        }
        useDragHandle
      />
    </div>
  );
};

export default ArticleSidebar;
