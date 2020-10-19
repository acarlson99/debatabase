import React from "react";
import "../App.css";

const TopBar = ({callback}) => {
  return <div className="TopBar" onClick={callback} />;
};
export default TopBar;
