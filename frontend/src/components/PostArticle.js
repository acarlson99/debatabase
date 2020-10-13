import React, { useState } from "react";
import axios from "axios";
import { SERVER_PORT, SERVER_HOST } from "../Const";

const PostArticle = (props) => {
  const [name, setName] = useState("");
  const [url, setURL] = useState("");
  const [tags, setTags] = useState("");
  const [description, setDescription] = useState("");

  return (
    <div className="Post">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          const data = {
            name: name,
            url: url,
            tags: tags.split(","),
            description: description,
          };
          console.log("POSTING TAG:", data);
          axios
            .post(
              `http://${SERVER_HOST}:${SERVER_PORT}/api/upload/article`,
              JSON.stringify(data)
            )
            .then((res) => console.log("RES:", res))
            .catch((err) => console.log(err));
          setName("");
          setURL("");
          setTags("");
          setDescription("");
        }}
      >
        <input
          placeholder="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <input
          placeholder="URL"
          value={url}
          onChange={(e) => setURL(e.target.value)}
          required
        />
        <br />
        <textarea
          placeholder="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
        />
        <br />
        <textarea
          placeholder="CSV of tags"
          value={tags}
          onChange={(e) => setTags(e.target.value)}
        />
        <br />
        <button type="submit">Create Article</button>
      </form>
    </div>
  );
};

export default PostArticle;
