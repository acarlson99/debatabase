import React, { useContext } from "react";
import { SortableContainer,
  SortableElement,
  sortableHandle, } from "react-sortable-hoc";
import arrayMove from "array-move";
import ArticleListContext from "../contexts/ArticleList";

// TODO: fix `findDOMNode` error (prolly use refs)

const DragHandle = sortableHandle(() => (
  <div className="DragHandle noHL">::</div>
));

const SortableItem = SortableElement(({ value, deleteArticle }) => (
  <div className="SortableElement Article">
    <div className="rowC">
      <DragHandle />
      {JSON.stringify(value)}
    </div>
    <button onClick={deleteArticle}>-</button>
  </div>
));

const SortableList = SortableContainer(({ items, deleteArticle }) => (
  <div>
    {items.map((value, index) => (
      <SortableItem
        key={"sortable-item" + JSON.stringify(items)}
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
    <div className="ArticleSidebar scroll">
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
