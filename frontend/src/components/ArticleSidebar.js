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

const SortableItem = SortableElement(({ value }) => (
  <div className="SortableElement">
    <DragHandle />
    {JSON.stringify(value)}
  </div>
));

const SortableList = SortableContainer(({ items }) => {
  return (
    <div>
      {items.map((value, index) => (
        <SortableItem key={`item-${index}`} index={index} value={value} />
      ))}
    </div>
  );
});

const ArticleSidebar = () => {
  const c = useContext(ArticleListContext);
  const articles = c.articles;
  const setArticles = c.setArticles;

  return (
    <div className="ArticleSidebar">
      <SortableList
        items={articles}
        onSortEnd={({ oldIndex, newIndex }) => {
          setArticles(arrayMove(articles, oldIndex, newIndex));
        }}
        useDragHandle
      />
    </div>
  );
};

export default ArticleSidebar;
