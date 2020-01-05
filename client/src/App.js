import React from "react";
import "./App.css";
// import the Container Component from the semantic-ui-react
import { Container } from "semantic-ui-react";
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