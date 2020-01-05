import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import { Header, Segment, Form, Message, Input, TextArea, Button, Table, Icon } from "semantic-ui-react";
//import CreateMessage from "./CreateMessage";

let endpoint = "http://127.0.0.1:8080";

// could make it a react funcitons
class ComposeForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      to: "",
      subject: "",
      body: "",
      scheduledAt: null,
      replySent: false
    }
  }

  handleToChange = (e) => {
    this.setState({
      to: e.target.value
    });
  }

  handleSubjectChange = (e) => {
    this.setState({
      subject: e.target.value
    });
  }

  handleBodyChange = (e) => {
    this.setState({
      body: e.target.value
    });
  }

  sentStatus = () => {
    if(this.state.replySent) {
      return   <Message success header='Sent' content='Your Message has been sent.'/>
    } else {
      return ""
    }
  }

  sendMessage = () => {
    let a = this.props.GetAccountIDFn();
    let newSender = this.props.AccountEmail;
    let newRecipients = [];
    let newSubject = this.state.subject;

    let tostr = "";
    let apiStr = "/message";
    apiStr += "?"+a.name+"="+a.value;

    tostr = this.state.to; 
    newRecipients = tostr.toString().split(", ");

    let data = JSON.stringify({ 
      ParentMid: 0,
      ScheduledAt: null,
      SenderEmail: newSender,
      Recipients: newRecipients,
      Subject: newSubject,
      Body: btoa(this.state.body)
    })

    axios
      .post(endpoint + apiStr, data, 
        {
          headers: {
            "Content-Type": "application/json"
        }
      })
      .then(res => {
        console.log(res);
        this.setState({ 
          replySent: true, 
          body: ""
        });
      },(error)=>{
        console.log(error);
        this.setState({ 
          replySent: false
        });
      });    
  }

  createDisplay = () => {
    if(this.props.IsLoggedIn) {
      return (
      <Form success>
        <Table compact>
          <Table.Body>
            <Table.Row>
              <Table.Cell>From: </Table.Cell>            
              <Table.Cell>{this.props.AccountEmail}</Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>To: </Table.Cell>
              <Table.Cell>
                <Input fluid id="to" placeholder="Recipient, Recipient, ..." rows={1} value={this.state.toRecipients} onChange={this.handleToChange}/>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Subject: </Table.Cell>
              <Table.Cell>
                <Input fluid id="subject" placeholder="Subject" rows={1} onChange={this.handleSubjectChange}/>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Body: </Table.Cell>
              <Table.Cell>
                <TextArea id="body" placeholder="Message" rows={3} onChange={this.handleBodyChange}/>
              </Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        {this.sentStatus()}
        <Button icon onClick={this.sendMessage}>
          <Icon name='send'/>Send</Button>
      </Form>
      );
    } else { 
      return (
        <></>
      );
    }
  }

  render() {
    return (
      <div>
        <Segment>
          <Header className="header" as="h3">{this.props.ComponentName}</Header>
          {this.createDisplay()}          
        </Segment>
      </div>
    );
  }
}



ComposeForm.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  AccountEmail: PropTypes.string.isRequired,

  ComponentName: PropTypes.string,
  FormatTimeFn: PropTypes.func
}

ComposeForm.defaultProps = {
  ComponentName: 'Compose Form',
  FormatTimeFn: identityFn
}

function identityFn(x) { return x; }

export default ComposeForm;
