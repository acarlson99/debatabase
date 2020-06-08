import React, { useState } from "react";
import axios from "axios";
import { SERVER_PORT, SERVER_HOST } from "../Const";

const PostTag = (props) => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  return (
    <div className="Post">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          const data = {
            name: name,
            description: description,
          };
          console.log("POSTING TAG:", data);
          axios
            .post(
              `http://${SERVER_HOST}:${SERVER_PORT}/api/upload/tag`,
              JSON.stringify(data)
            )
            .then((res) => console.log("RES:", res))
            .catch((err) => console.log(err));
          setName("");
          setDescription("");
        }}
      >
        <input
          placeholder="tag name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          maxLength={16}
          required
        />
        <br />
        <textarea
          placeholder="tag description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          maxLength={256}
          required
        />
        <br />
        <button type="submit">Create Tag</button>
      </form>
    </div>
  );
};

export default PostTag;
