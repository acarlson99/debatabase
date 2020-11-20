import React from "react";
import { Navbar, Nav } from "react-bootstrap";
import "../App.css";

const NLink = ({ href, body }) => (
  <Nav.Link className="BarLink" href={href}>
    {body}
  </Nav.Link>
);

const TopBar = () => {
  const elems = [
    ["/", "Home"],
    ["/present", "present"],
    ["/upload", "upload"],
    ["/upload/tag", "upload tag"],
    ["/upload/article", "upload article"],
  ];

  return (
    <Navbar className="TopBar">
      <Nav>
        {elems
          .map((e, i) => <NLink key={e + i} href={e[0]} body={e[1]} />)
          .reduce((a, v) => [...a, v, "|"], [])
          .slice(0, -1)}
      </Nav>
    </Navbar>
  );
};

export default TopBar;
