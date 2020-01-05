import React from "react";
import "./App.css";
// import the Container Component from the semantic-ui-react
import { Container, Header, Image } from "semantic-ui-react";
// import the ToDoList component
import ToDoList from "./To-Do-List";
import Top from "./Top";

function App() {
  return (
    <div>
      <Container>
        <Top />        
      </Container>
    </div>
  );
}
export default App;
//<ToDoList />